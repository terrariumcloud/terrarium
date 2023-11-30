package browse

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/terrariumcloud/terrarium/internal/release/services/release"
	v1 "github.com/terrariumcloud/terrarium/internal/restapi/modules/v1"

	"github.com/apparentlymart/go-versions/versions"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	releaseServices "github.com/terrariumcloud/terrarium/internal/release/services"
	"github.com/terrariumcloud/terrarium/internal/restapi"
)

type browseHttpService struct {
	registrarClient      services.RegistrarClient
	versionManagerClient services.VersionManagerClient
	releasesClient       releaseServices.BrowseClient
	responseHandler      restapi.ResponseHandler
	errorHandler         restapi.ErrorHandler
}

func (h *browseHttpService) GetHttpHandler(mountPath string) http.Handler {
	router := h.createRouter(mountPath)

	return handlers.CombinedLoggingHandler(os.Stdout, router)
}

func New(registrarClient services.RegistrarClient, versionManagerClient services.VersionManagerClient, releasesClient releaseServices.BrowseClient) *browseHttpService {
	return &browseHttpService{registrarClient: registrarClient, versionManagerClient: versionManagerClient, releasesClient: releasesClient}
}

func (h *browseHttpService) createRouter(mountPath string) *mux.Router {
	apiPrefix := fmt.Sprintf("%s/api", mountPath)
	log.Printf("Prefix for browse endpoint implementation: %s", apiPrefix)
	rootRouter := mux.NewRouter()
	rootRouter.Handle("/healthz", h.healthHandler()).Methods(http.MethodGet)
	apiRouter := rootRouter.PathPrefix(apiPrefix).Subrouter()
	apiRouter.Use(otelmux.Middleware("browse"))
	apiRouter.StrictSlash(true)
	apiRouter.Handle("/modules/{organization_name}/{name}/{provider}", h.getModuleMetadataHandler()).Methods(http.MethodGet)
	apiRouter.Handle("/modules", h.getModuleListHandler()).Methods(http.MethodGet)
	apiRouter.Handle("/releases", h.getReleasesHandler()).Methods(http.MethodGet)
	apiRouter.Handle("/organizations", h.getOrganizationsHandler()).Methods(http.MethodGet)
	apiRouter.Handle("/types", h.getReleaseTypesHandler()).Methods(http.MethodGet)
	rootRouter.PathPrefix("/").Handler(getFrontendSpaHandler())
	return rootRouter
}

func (h *browseHttpService) healthHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		h.responseHandler.Write(rw, "OK", http.StatusOK)
	})
}

// GetModuleListHandler will return a list of all published module.
func (h *browseHttpService) getModuleListHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		if registrarResponse, err := h.registrarClient.ListModules(r.Context(), &services.ListModulesRequest{}); err != nil {
			span.RecordError(err)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of modules from backend service"), http.StatusInternalServerError)
			return
		} else {
			data, _ := json.Marshal(createModulesResponse(registrarResponse.Modules))

			rw.Header().Add("Content-Type", "application/json")
			_, _ = rw.Write(data)
		}
	})
}

func (h *browseHttpService) getModuleMetadataHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		moduleName := v1.GetModuleNameFromRequest(r)
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(
			attribute.String("module.name", moduleName),
		)

		registrarResponse, err := h.registrarClient.GetModule(ctx, &services.GetModuleRequest{Name: moduleName})
		if err != nil {
			span.RecordError(err)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the module details from backend service"), http.StatusInternalServerError)
			return
		}

		versionResponse, err := h.versionManagerClient.ListModuleVersions(ctx, &services.ListModuleVersionsRequest{Module: moduleName})
		if err != nil {
			span.RecordError(err)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of versions from backend service"), http.StatusInternalServerError)
			return
		}

		var filteredVersions []string
		for _, moduleVersion := range versionResponse.Versions {
			parsedVersion := versions.MustParseVersion(moduleVersion)

			if parsedVersion.GreaterThan(versions.MustParseVersion("0.0.0")) {
				filteredVersions = append(filteredVersions, moduleVersion)
			}

		}
		versionResponse.Versions = filteredVersions

		data := createModuleMetadataResponse(registrarResponse.GetModule(), versionResponse.Versions)
		h.responseHandler.Write(rw, data, http.StatusOK)
	})
}

// GetReleasesHandler will return a list of all releases published.
func (h *browseHttpService) getReleasesHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var maxAge uint64

		values := r.URL.Query()
		age := values.Get("age")

		parsedAge, err := time.ParseDuration(age)
		if err != nil || parsedAge.Seconds() < 3600 {
			maxAge = 3600
		} else {
			maxAge = uint64(parsedAge.Seconds())
		}

		MaxAgeSeconds := release.ConvertUint64ToInt64(maxAge)

		ctx := r.Context()
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(
			attribute.Int64("release.maxAge", MaxAgeSeconds),
		)

		response, err := h.releasesClient.ListReleases(r.Context(), &releaseServices.ListReleasesRequest{
			MaxAgeSeconds: &maxAge,
		})
		if err != nil {
			span.RecordError(err)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of releases from backend service"), http.StatusInternalServerError)
			return
		}

		data, _ := json.Marshal(createReleaseResponse(response.Releases))

		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(data)
	})
}

// getReleaseTypesHandler will return a list of all types available.
func (h *browseHttpService) getReleaseTypesHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		releaseTypesResponse, err := h.releasesClient.ListReleaseTypes(r.Context(), &releaseServices.ListReleaseTypesRequest{})
		if err != nil {
			span.RecordError(err)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of release types from backend service"), http.StatusInternalServerError)
			return
		}

		types := releaseTypesResponse.Types
		if types == nil {
			types = make([]string, 0)
		}
		data, _ := json.Marshal(types)

		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(data)
	})
}

// getOrganizationsHandler will return a list of all organizations available.
func (h *browseHttpService) getOrganizationsHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		releaseOrganizationsResponse, err := h.releasesClient.ListOrganization(r.Context(), &releaseServices.ListOrganizationRequest{})
		if err != nil {
			span.RecordError(err)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of organizations from backend service"), http.StatusInternalServerError)
			return
		}
		organizations := releaseOrganizationsResponse.Organizations
		if organizations == nil {
			organizations = make([]string, 0)
		}
		data, _ := json.Marshal(organizations)

		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(data)
	})
}
