package v1

import (
	"fmt"
	"github.com/gorilla/mux"
	pb "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
	"google.golang.org/grpc"
	"net/http"
)

func getModuleNameFromRequest(r *http.Request) string {
	params := mux.Vars(r)
	orgName := params["organization_name"]
	moduleName := params["name"]
	providerName := params["provider"]
	return fmt.Sprintf("%s/%s/%s", orgName, moduleName, providerName)
}

func getVersionedModuleFromRequest(r *http.Request) *pb.Module {
	params := mux.Vars(r)
	orgName := params["organization_name"]
	moduleName := params["name"]
	providerName := params["provider"]
	version := params["version"]
	return &pb.Module{
		Name:    fmt.Sprintf("%s/%s/%s", orgName, moduleName, providerName),
		Version: version,
	}
}

func createModuleVersionsResponse(versions []string) *ModuleVersionResponse {
	var structuredVersion []*ModuleVersionItem

	for _, version := range versions {
		structuredVersion = append(structuredVersion, &ModuleVersionItem{
			Version: version,
		})
	}
	return &ModuleVersionResponse{
		Modules: []*ModuleVersions{
			{
				Versions: structuredVersion,
			},
		},
	}
}

func closeClient(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {

	}
}
