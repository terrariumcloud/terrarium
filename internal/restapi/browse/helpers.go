package browse

import (
	"github.com/terrariumcloud/terrarium/internal/module/services"

	"google.golang.org/grpc"
)

type moduleItem struct {
	Organization string   `json:"organization"`
	Name         string   `json:"name"`
	Provider     string   `json:"provider"`
	Description  string   `json:"description"`
	SourceUrl    string   `json:"source_url"`
	Maturity     string   `json:"maturity,omitempty"`
	Versions     []string `json:"versions,omitempty"`
}

type modulesResponse struct {
	Modules []*services.ModuleMetadata `json:"modules"`
}

func createModulesResponse(modules []*services.ModuleMetadata) *modulesResponse {
	return &modulesResponse{
		Modules: modules,
	}
}

func createModuleMetadataResponse(moduleMetadata *services.ModuleMetadata, moduleVersions []string) *moduleItem {
	return &moduleItem{
		Organization: moduleMetadata.Organization,
		Name:         moduleMetadata.Name,
		Provider:     moduleMetadata.Provider,
		Description:  moduleMetadata.Description,
		SourceUrl:    moduleMetadata.SourceUrl,
		Versions:     moduleVersions,
	}
}

func closeClient(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {

	}
}
