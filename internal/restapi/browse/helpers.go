package browse

import (
	"github.com/terrariumcloud/terrarium/internal/module/services"
	release "github.com/terrariumcloud/terrarium/internal/release/services"
	providerServices "github.com/terrariumcloud/terrarium/internal/provider/services"

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

type providerItem struct {
	Organization string   `json:"organization"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	SourceUrl    string   `json:"source_url"`
	Maturity     string   `json:"maturity,omitempty"`
	Versions     []string `json:"versions,omitempty"`
}

type modulesResponse struct {
	Modules []*services.ModuleMetadata `json:"modules"`
}

type providersResponse struct {
	Providers []*providerServices.ListProviderItem `json:"providers"`
}

func createModulesResponse(modules []*services.ModuleMetadata) *modulesResponse {
	return &modulesResponse{
		Modules: modules,
	}
}

func createProvidersResponse(providers []*providerServices.ListProviderItem) *providersResponse {
	return &providersResponse{
		Providers: providers,
	}
}

func createModuleMetadataResponse(moduleMetadata *services.ModuleMetadata, moduleVersions []string) *moduleItem {
	return &moduleItem{
		Organization: moduleMetadata.Organization,
		Name:         moduleMetadata.Name,
		Provider:     moduleMetadata.Provider,
		Description:  moduleMetadata.Description,
		SourceUrl:    moduleMetadata.SourceUrl,
		Maturity:     moduleMetadata.Maturity.String(),
		Versions:     moduleVersions,
	}
}

func createProviderMetadataResponse(providerMetadata *providerServices.ListProviderItem, providerVersions []string) *providerItem {
	return &providerItem{
		Organization: providerMetadata.Organization,
		Name: 		  providerMetadata.Name,
		Description:  providerMetadata.Description,
		SourceUrl: 	  providerMetadata.SourceUrl,
		Maturity: 	  providerMetadata.Maturity,
		Versions: 	  providerVersions,
	}
}

type releaseResponse struct {
	Releases []*release.Release `json:"releases"`
}

func createReleaseResponse(releases []*release.Release) *releaseResponse {
	return &releaseResponse{
		Releases: releases,
	}
}

func closeClient(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {

	}
}
