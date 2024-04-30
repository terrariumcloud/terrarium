package v1

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/terrariumcloud/terrarium/internal/provider/services"
	pb "github.com/terrariumcloud/terrarium/pkg/terrarium/provider"
)

func GetProviderNameFromRequest(r *http.Request) string {
	params := mux.Vars(r)
	orgName := params["organization_name"]
	providerName := params["name"]
	return fmt.Sprintf("%s/%s", orgName, providerName)
}

func GetProviderInputsFromRequest(r *http.Request) (string, string, string) {
	params := mux.Vars(r)
	return params["version"], params["os"], params["arch"]
}

func GetProviderLocationFromRequest(r *http.Request) *services.ProviderRequest {
	params := mux.Vars(r)
	orgName := params["organization_name"]
	providerName := params["name"]
	version := params["version"]
	os := params["os"]
	arch := params["arch"]
	return &services.ProviderRequest{
		Name:    fmt.Sprintf("%s/%s", orgName, providerName),
		Version: version,
		Os:      os,
		Arch:    arch,
	}
}

func GetVersionedProviderFromRequest(r *http.Request) *pb.Provider {
	params := mux.Vars(r)
	orgName := params["organization_name"]
	providerName := params["name"]
	version := params["version"]
	return &pb.Provider{
		Name:    fmt.Sprintf("%s/%s", orgName, providerName),
		Version: version,
	}
}
