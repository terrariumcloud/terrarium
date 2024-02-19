package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/errgo.v2/errors"

	"github.com/terrariumcloud/terrarium/internal/restapi"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

)


type providersV1HttpService struct {
	responseHandler restapi.ResponseHandler
	errorHandler    restapi.ErrorHandler
}

// Structs to load data into (from a JSON file for now, will be from DB later)

type Protocols []string

type GPGPublicKey struct {
	KeyID          	string	`json:"key_id"`
	ASCIIArmor     	string 	`json:"ascii_armor"`
	TrustSignature	string 	`json:"trust_signature"`
	Source         	string 	`json:"source"`
	SourceURL      	string 	`json:"source_url"`
}

type SigningKeys struct {
	GPGPublicKeys	[]GPGPublicKey	`json:"gpg_public_keys"`
}

type ProviderMetadata struct {
	OS                  string       `json:"os"`
	Arch                string       `json:"arch"`
	Filename            string       `json:"filename,omitempty"`
	DownloadURL         string       `json:"download_url,omitempty"`
	ShasumsURL          string       `json:"shasums_url,omitempty"`
	ShasumsSignatureURL	string       `json:"shasums_signature_url,omitempty"`
	Shasum              string       `json:"shasum,omitempty"`
	SigningKeys         SigningKeys  `json:"signing_keys,omitempty"`
}

type VersionData struct {
	Protocols	Protocols    		`json:"protocols"`
	Platforms	[]ProviderMetadata	`json:"platforms"`
}

type ProviderData map[string]*VersionData

// Structs to load response into (for listing versions for a specific provider)

type Platform struct {
	OS  	string	`json:"os"`
	Arch	string	`json:"arch"`
}

type VersionItem struct {
	Version   	string		`json:"version"`
	Protocols	Protocols   `json:"protocols"`
	Platforms 	[]Platform	`json:"platforms"`
}

type ProviderVersionsResponse struct {
	Versions	[]VersionItem	`json:"versions"`
}

// Structs to load response into (for a provider's metadata)

type PlatformMetadataResponse struct {
	Protocols     	Protocols 	`json:"protocols"`
	OS            	string    	`json:"os"`
	Arch          	string    	`json:"arch"`
	Filename      	string    	`json:"filename"`
	DownloadURL   	string    	`json:"download_url"`
	ShasumsURL    	string    	`json:"shasums_url"`
	ShasumsSigURL	string    	`json:"shasums_signature_url"`
	Shasum        	string    	`json:"shasum"`
	SigningKeys   	SigningKeys	`json:"signing_keys"`
}


func New() *providersV1HttpService {
	return &providersV1HttpService{}
}

func (h *providersV1HttpService) GetHttpHandler(mountPath string) http.Handler {
	router := h.createRouter(mountPath)
	return handlers.CombinedLoggingHandler(os.Stdout, router)
}

func (h *providersV1HttpService) createRouter(mountPath string) *mux.Router {
	prefix := fmt.Sprintf("%s/v1", mountPath)
	log.Printf("prefix for registry implementation: %s", prefix)
	r := mux.NewRouter()
	r.Handle("/healthz", h.healthHandler()).Methods(http.MethodGet)
	sr := r.PathPrefix(prefix).Subrouter()
	sr.Use(otelmux.Middleware("providers-v1"))
	sr.StrictSlash(true)
	sr.Handle("/{organization_name}/{name}/versions", h.getProviderVersionHandler()).Methods(http.MethodGet)
	sr.Handle("/{organization_name}/{name}/{version}/download/{os}/{arch}", h.downloadProviderHandler()).Methods(http.MethodGet)
	return r
}

func (h *providersV1HttpService) healthHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		h.responseHandler.Write(rw, "OK", http.StatusOK)
	})
}

func loadJSONData() map[string]ProviderData {
	filePath := "./random.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var obj map[string]ProviderData

	if err := json.Unmarshal(data, &obj); err != nil {
		panic(err)
	}

	return obj
}

// GetProviderVersionHandler will return a list of available versions for a given provider.
// This signifies to the requesting CLI if that provider is available to consume from the registry.
// Will return a 404 if a non-existent organization and/or provider is requested.
// This handler complies with the following implementation from the provider protocol
// https://developer.hashicorp.com/terraform/internals/provider-registry-protocol#list-available-versions
func (h *providersV1HttpService) getProviderVersionHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("getProviderVersionHandler")
		
		providerName := GetProviderNameFromRequest(r)

		ctx := r.Context()
		span := trace.SpanFromContext(ctx)
		span.AddEvent("Returning a list of available versions for a provider", trace.WithAttributes(attribute.String("Provider Name", providerName)))
		span.SetAttributes(
			attribute.String("provider.name", providerName),
		)

		obj := loadJSONData()

		var providerVersions ProviderVersionsResponse
		var platform Platform

		// Check if the provider ID exists
		if providerData, exists := obj[providerName]; exists {
			// Add the matched provider's version details to the ProviderVersionsResponse
			for version, versionMetadata := range providerData {
				var versionItem VersionItem
				versionItem.Version = version
				versionItem.Protocols = versionMetadata.Protocols
				for _, versionPlatforms := range versionMetadata.Platforms {
					platform.OS = versionPlatforms.OS
					platform.Arch = versionPlatforms.Arch
					versionItem.Platforms = append(versionItem.Platforms, platform)
				}
				providerVersions.Versions = append(providerVersions.Versions, versionItem) 
			}
			
		} else {
			errMsg := fmt.Sprintf("failed to retrieve the list of versions for %s", providerName)
			span.RecordError(errors.New(errMsg))
			h.errorHandler.Write(rw, errors.New(errMsg), http.StatusNoContent)
			return
		}

		data, _ := json.Marshal(providerVersions)
		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(data)
	})
}

// DownloadProviderHandler returns the download URL of and associated metadata about the distribution package
// for a particular version of a provider for a particular operating system and architecture.
// Terraform CLI uses this operation after it has selected the newest available version matching the configured
// version constraints, in order to find the zip archive containing the plugin itself.
// Will return a 404 if a non-existent version or os or arch is requested.
// This handler complies with the following implementation from the provider protocol
// https://developer.hashicorp.com/terraform/internals/provider-registry-protocol#find-a-provider-package
func (h *providersV1HttpService) downloadProviderHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("downloadProviderHandler")

		providerName := GetProviderNameFromRequest(r)

		ctx := r.Context()
		span := trace.SpanFromContext(ctx)
		span.AddEvent("Fetching provider's metadata", trace.WithAttributes(attribute.String("Provider Name", providerName)))
		span.SetAttributes(
			attribute.String("module.name", providerName),
		)

		obj := loadJSONData()

		version, os, arch := GetProviderInputsFromRequest(r)

		var providerMetadata PlatformMetadataResponse
		var outputExists bool

		// Check if the provider ID exists
		if providerData, exists := obj[providerName]; exists {
			// Check if the version exists for the provider
			if versionData, exists := providerData[version]; exists {
				for _, platform := range versionData.Platforms {
					if platform.OS == os && platform.Arch == arch {
					outputExists = true
					// Add the matched platform details to the providerMetadata
					providerMetadata.Protocols 	   	= versionData.Protocols
					providerMetadata.OS 		   	= platform.OS
					providerMetadata.Arch 		   	= platform.Arch
					providerMetadata.Filename 	   	= platform.Filename
					providerMetadata.DownloadURL   	= platform.DownloadURL
					providerMetadata.ShasumsURL    	= platform.ShasumsURL
					providerMetadata.ShasumsSigURL 	= platform.ShasumsSignatureURL
					providerMetadata.Shasum 	   	= platform.Shasum
					providerMetadata.SigningKeys  	= platform.SigningKeys
					break
					} else {
						outputExists = false
					}
				}
			} else {
				errMsg := fmt.Sprintf("failed to retrieve version: %s for: %s", version, providerName)
				span.RecordError(errors.New(errMsg))
				h.errorHandler.Write(rw, errors.New(errMsg), http.StatusNoContent)
				return
			}
		} else {
			errMsg := fmt.Sprintf("failed to retrieve: %s", providerName)
			span.RecordError(errors.New(errMsg))
			h.errorHandler.Write(rw, errors.New(errMsg), http.StatusNoContent)
			return
		}

		if outputExists {
			data, _ := json.Marshal(providerMetadata)
			rw.Header().Add("Content-Type", "application/json")
			_, _ = rw.Write(data)
		} else {
			errMsg := fmt.Sprintf("failed to retrieve the requested provider: %s for version: %s, os: %s and arch: %s", providerName, version, os, arch)
			span.RecordError(errors.New(errMsg))
			h.errorHandler.Write(rw, errors.New(errMsg), http.StatusNoContent)
			return
		}
	})
}