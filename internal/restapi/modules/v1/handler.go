package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/module/services/storage"
	"github.com/terrariumcloud/terrarium/internal/module/services/version_manager"
	"github.com/terrariumcloud/terrarium/internal/restapi"
	pb "github.com/terrariumcloud/terrarium/pkg/terrarium/module"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/errgo.v2/errors"
	"io"
	"log"
	"net/http"
	"os"
)

type modulesV1HttpService struct {
	responseHandler restapi.ResponseHandler
	errorHandler    restapi.ErrorHandler
}

type ModuleVersionItem struct {
	Version string `json:"version"`
}

type ModuleVersions struct {
	Versions []*ModuleVersionItem `json:"versions"`
}

type ModuleVersionResponse struct {
	Modules []*ModuleVersions `json:"modules"`
}

func (h *modulesV1HttpService) GetHttpHandler(mountPath string) http.Handler {
	router := h.createRouter(mountPath)
	return handlers.CombinedLoggingHandler(os.Stdout, router)
}

func New() *modulesV1HttpService {
	return &modulesV1HttpService{}
}

func (h *modulesV1HttpService) createRouter(mountPath string) *mux.Router {
	prefix := fmt.Sprintf("%s/v1", mountPath)
	log.Printf("Prefix for registry implementation: %s", prefix)
	r := mux.NewRouter()
	r.Handle("/healthz", h.healthHandler()).Methods(http.MethodGet)
	sr := r.PathPrefix(prefix).Subrouter()
	sr.StrictSlash(true)
	sr.Handle("/{organization_name}/{name}/{provider}/versions", h.getModuleVersionHandler()).Methods(http.MethodGet)
	sr.Handle("/{organization_name}/{name}/{provider}/{version}/download", h.downloadModuleHandler()).Methods(http.MethodGet)
	sr.Handle("/{organization_name}/{name}/{provider}/{version}/archive", h.archiveHandler()).Methods(http.MethodGet)
	return r
}

func (h *modulesV1HttpService) healthHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		h.responseHandler.Write(rw, "OK", http.StatusOK)
	})
}

// GetModuleVersionHandler will return a list of available versions for a given module.
// This signifies to the requesting CLI if that module is available to consume from the registry.
// Will return a 404 if a non-existent organization and/or module is requested.
// This handler complies with the following implementation from the module protocol
// https://www.terraform.io/internals/module-registry-protocol#download-source-code-for-a-specific-module-version
func (h *modulesV1HttpService) getModuleVersionHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("getModuleVersionHandler")
		moduleName := GetModuleNameFromRequest(r)
		conn, err := grpc.Dial(version_manager.VersionManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to connect to '%s': %v", version_manager.VersionManagerEndpoint, err)
			h.errorHandler.Write(rw, errors.New("failed connecting to the version manager backend service"), http.StatusInternalServerError)
			return
		}
		defer closeClient(conn)

		client := services.NewVersionManagerClient(conn)

		versionResponse, err2 := client.ListModuleVersions(context.TODO(), &services.ListModuleVersionsRequest{Module: moduleName})
		if err2 != nil {
			log.Printf("Failed GRPC call with error: %v", err2)
			h.errorHandler.Write(rw, errors.New("failed to retrieve the list of versions from backend service"), http.StatusInternalServerError)
			return
		}
		data, _ := json.Marshal(createModuleVersionsResponse(versionResponse.Versions))
		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(data)
	})
}

func (h *modulesV1HttpService) downloadModuleHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("downloadModuleHandler")
		// At this stage there is no validation about parameters given to that function, validation is done as part of download.
		rw.Header().Add("X-Terraform-Get", "./archive?archive=zip")
		h.responseHandler.Write(rw, nil, http.StatusNoContent)
	})
}

// archiveHandler performs a fetch of the restapi.d module source code from the chosen backing store and presents it to the client
// As part of the module flow clients are redirected here from the DownloadModuleHandler x-terraform-get header. his handler
// makes the stored registry code available to the client
func (h *modulesV1HttpService) archiveHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		conn, err := grpc.Dial(storage.StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to connect: %v", err)
			h.errorHandler.Write(rw, errors.New("failed connecting to the storage backend service"), http.StatusInternalServerError)
			return
		}
		defer closeClient(conn)

		client := services.NewStorageClient(conn)

		downloadStream, err2 := client.DownloadSourceZip(context.TODO(), &pb.DownloadSourceZipRequest{
			Module: getVersionedModuleFromRequest(r),
		})
		if err2 != nil {
			log.Printf("Failed to connect: %v", err2)
			h.errorHandler.Write(rw, errors.New("failed to initiate the downlowd of the archive from storage backend service"), http.StatusInternalServerError)
			return
		}

		r.Header.Set("Content-Type", "application/zip")
		for {
			chunk, err := downloadStream.Recv()
			if err == io.EOF {
				return
			}
			_, _ = rw.Write(chunk.ZipDataChunk)
		}
	})
}
