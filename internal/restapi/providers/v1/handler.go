package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/terrariumcloud/terrarium/internal/provider/services"
	"github.com/terrariumcloud/terrarium/internal/restapi"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type providersV1HttpService struct {
	jsonHandler     services.ProviderVersionManager
	responseHandler restapi.ResponseHandler
	errorHandler    restapi.ErrorHandler
}

func New(jsonHandler services.ProviderVersionManager) *providersV1HttpService {
	return &providersV1HttpService{jsonHandler: jsonHandler}
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
		span.SetAttributes(
			attribute.String("provider.name", providerName),
		)

		providerVersions, err := h.jsonHandler.ListProviderVersions(providerName)
		if err != nil {
			h.errorHandler.Write(rw, err, http.StatusNotFound)
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
		span.SetAttributes(
			attribute.String("module.name", providerName),
		)

		providerVersion, providerOS, providerArch := GetProviderInputsFromRequest(r)

		providerMetadata, err := h.jsonHandler.GetVersionData(providerName, providerVersion, providerOS, providerArch)
		if err != nil {
			h.errorHandler.Write(rw, err, http.StatusNotFound)
			return
		}

		data, _ := json.Marshal(providerMetadata)
		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(data)
	})
}
