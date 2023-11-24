package browse

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/terrariumcloud/terrarium/internal/module/services/registrar"
	"github.com/terrariumcloud/terrarium/internal/module/services/version_manager"
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
	responseHandler restapi.ResponseHandler
	errorHandler    restapi.ErrorHandler
}

func (h *browseHttpService) GetHttpHandler(mountPath string) http.Handler {
	router := h.createRouter(mountPath)

	return handlers.CombinedLoggingHandler(os.Stdout, router)
}

func New() *browseHttpService {
	return &browseHttpService{}
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
	apiRouter.Handle("/releases/{age}", h.getReleasesHandler()).Methods(http.MethodGet)
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
		conn, err := services.CreateGRPCConnection(registrar.RegistrarServiceEndpoint)
		if err != nil {
			log.Printf("Failed to connect to '%s': %v", registrar.RegistrarServiceEndpoint, err)
			h.errorHandler.Write(rw, errors.New("failed connecting to the registrar backend service"), http.StatusInternalServerError)
			return
		}
		defer closeClient(conn)

		client := services.NewRegistrarClient(conn)

		registrarResponse, err2 := client.ListModules(r.Context(), &services.ListModulesRequest{})

		if err2 != nil {
			log.Printf("Failed GRPC call with error: %v", err2)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of modules from backend service"), http.StatusInternalServerError)
			return
		}

		data, _ := json.Marshal(createModulesResponse(registrarResponse.Modules))

		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(data)
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

		conn, err := services.CreateGRPCConnection(registrar.RegistrarServiceEndpoint)
		if err != nil {
			log.Printf("Failed to connect to '%s': %v", registrar.RegistrarServiceEndpoint, err)
			h.errorHandler.Write(rw, errors.New("failed connecting to the registrar backend service"), http.StatusInternalServerError)
			return
		}
		defer closeClient(conn)

		clientRegistrar := services.NewRegistrarClient(conn)
		registrarResponse, err := clientRegistrar.GetModule(ctx, &services.GetModuleRequest{Name: moduleName})
		if err != nil {
			log.Printf("Failed GRPC call with error: %v", err)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of modules from backend service"), http.StatusInternalServerError)
			return
		}

		connVersion, err := services.CreateGRPCConnection(version_manager.VersionManagerEndpoint)
		if err != nil {
			log.Printf("Failed to connect to '%s': %v", version_manager.VersionManagerEndpoint, err)
			h.errorHandler.Write(rw, errors.New("failed connecting to the version manager backend service"), http.StatusInternalServerError)
			return
		}
		defer closeClient(connVersion)

		clientVersion := services.NewVersionManagerClient(connVersion)
		versionResponse, err := clientVersion.ListModuleVersions(ctx, &services.ListModuleVersionsRequest{Module: moduleName})

		var filteredVersions []string
		for _, moduleVersion := range versionResponse.Versions {
			parsedVersion := versions.MustParseVersion(moduleVersion)

			if parsedVersion.GreaterThan(versions.MustParseVersion("0.0.0")) {
				filteredVersions = append(filteredVersions, moduleVersion)
			}

		}
		versionResponse.Versions = filteredVersions

		if err != nil {
			log.Printf("Failed GRPC call with error: %v", err)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of versions from backend service"), http.StatusInternalServerError)
			return
		}

		data := createModuleMetadataResponse(registrarResponse.GetModule(), versionResponse.Versions)
		h.responseHandler.Write(rw, data, http.StatusOK)
	})
}

// GetReleasesHandler will return a list of all releases published.
func (h *browseHttpService) getReleasesHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var maxAge uint64

		ageStr := mux.Vars(r)["age"]
		if ageStr == "" {
			fmt.Println("Age not found setting default age to 7 days")
			maxAge = 604800
		}

		maxAge, err := strconv.ParseUint(ageStr, 10, 64)
		if err != nil {
			fmt.Println("Error converting age to uint64 setting default value to 7 days", err)
			maxAge = 604800
		}

		conn, err := services.CreateGRPCConnection(release.ReleaseServiceEndpoint)
		if err != nil {
			log.Printf("Failed to connect to '%s': %v", release.ReleaseServiceEndpoint, err)
			h.errorHandler.Write(rw, errors.New("failed connecting to the release backend service"), http.StatusInternalServerError)
			return
		}
		defer closeClient(conn)

		client := releaseServices.NewBrowseClient(conn)

		releaseResponse, err2 := client.ListReleases(r.Context(), &releaseServices.ListReleasesRequest{
			MaxAgeSeconds: &maxAge,
		})

		if err2 != nil {
			log.Printf("Failed GRPC call with error: %v", err2)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of releases from backend service"), http.StatusInternalServerError)
			return
		}

		data, _ := json.Marshal(createReleaseResponse(releaseResponse.Releases))

		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(data)
	})
}

// getReleaseTypesHandler will return a list of all types available.
func (h *browseHttpService) getReleaseTypesHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		conn, err := services.CreateGRPCConnection(release.ReleaseServiceEndpoint)
		if err != nil {
			log.Printf("Failed to connect to '%s': %v", release.ReleaseServiceEndpoint, err)
			h.errorHandler.Write(rw, errors.New("failed connecting to the release backend service"), http.StatusInternalServerError)
			return
		}
		defer closeClient(conn)

		client := releaseServices.NewBrowseClient(conn)

		releaseTypesResponse, err2 := client.ListReleaseTypes(r.Context(), &releaseServices.ListReleaseTypesRequest{})

		if err2 != nil {
			log.Printf("Failed GRPC call with error: %v", err2)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of release types from backend service"), http.StatusInternalServerError)
			return
		}

		data, _ := json.Marshal(releaseTypesResponse.Types)

		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(data)
	})
}

// getOrganizationsHandler will return a list of all organizations available.
func (h *browseHttpService) getOrganizationsHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		conn, err := services.CreateGRPCConnection(release.ReleaseServiceEndpoint)
		if err != nil {
			log.Printf("Failed to connect to '%s': %v", release.ReleaseServiceEndpoint, err)
			h.errorHandler.Write(rw, errors.New("failed connecting to the release backend service"), http.StatusInternalServerError)
			return
		}
		defer closeClient(conn)

		client := releaseServices.NewBrowseClient(conn)

		releaseOrganizationsResponse, err2 := client.ListOrganization(r.Context(), &releaseServices.ListOrganizationRequest{})

		if err2 != nil {
			log.Printf("Failed GRPC call with error: %v", err2)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of release types from backend service"), http.StatusInternalServerError)
			return
		}

		data, _ := json.Marshal(releaseOrganizationsResponse.Organizations)

		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(data)
	})
}
