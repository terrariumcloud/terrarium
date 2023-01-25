package browse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/terrariumcloud/terrarium/internal/module/services/registrar"
	"github.com/terrariumcloud/terrarium/internal/module/services/version_manager"
	v1 "github.com/terrariumcloud/terrarium/internal/restapi/modules/v1"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/restapi"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

type browseHttpService struct {
	responseHandler restapi.ResponseHandler
	errorHandler    restapi.ErrorHandler
}

func (h *browseHttpService) GetHttpHandler(mountPath string) http.Handler {
	router := h.createRouter(mountPath)
	router.Use(otelmux.Middleware("browse"))
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
	apiRouter.StrictSlash(true)
	apiRouter.Handle("/modules/{organization_name}/{name}/{provider}", h.getModuleMetadataHandler()).Methods(http.MethodGet)
	apiRouter.Handle("/modules", h.getModuleListHandler()).Methods(http.MethodGet)
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

		registrarResponse, err2 := client.ListModules(context.TODO(), &services.ListModulesRequest{})

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
		conn, err := services.CreateGRPCConnection(registrar.RegistrarServiceEndpoint)
		if err != nil {
			log.Printf("Failed to connect to '%s': %v", registrar.RegistrarServiceEndpoint, err)
			h.errorHandler.Write(rw, errors.New("failed connecting to the registrar backend service"), http.StatusInternalServerError)
			return
		}
		defer closeClient(conn)

		connVersion, err := services.CreateGRPCConnection(version_manager.VersionManagerEndpoint)
		if err != nil {
			log.Printf("Failed to connect to '%s': %v", version_manager.VersionManagerEndpoint, err)
			h.errorHandler.Write(rw, errors.New("failed connecting to the version manager backend service"), http.StatusInternalServerError)
			return
		}
		defer closeClient(connVersion)

		clientRegistrar := services.NewRegistrarClient(conn)

		registrarResponse, err := clientRegistrar.GetModule(context.TODO(), &services.GetModuleRequest{Name: moduleName})
		if err != nil {
			log.Printf("Failed GRPC call with error: %v", err)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of modules from backend service"), http.StatusInternalServerError)
			return
		}

		clientVersion := services.NewVersionManagerClient(connVersion)
		versionResponse, err := clientVersion.ListModuleVersions(context.TODO(), &services.ListModuleVersionsRequest{Module: moduleName})
		if err != nil {
			log.Printf("Failed GRPC call with error: %v", err)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of versions from backend service"), http.StatusInternalServerError)
			return
		}

		data := createModuleMetadataResponse(registrarResponse.GetModule(), versionResponse.Versions)
		h.responseHandler.Write(rw, data, http.StatusOK)
	})
}
