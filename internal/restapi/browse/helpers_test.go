package browse

import (
	"reflect"
	"testing"

	"github.com/terrariumcloud/terrarium/internal/module/services"
	providerServices "github.com/terrariumcloud/terrarium/internal/provider/services"
)

func Test_createModuleMetadataResponse(t *testing.T) {
	type args struct {
		moduleMetadata *services.ModuleMetadata
		moduleVersions []string
	}
	tests := []struct {
		name string
		args args
		want *moduleItem
	}{
		{
			name: "Module Metadata Response",
			args: args{
				moduleMetadata: &services.ModuleMetadata{
					Organization: "cie",
					Name:         "test-module",
					Provider:     "aws",
					Description:  "This is the description for the module it is supposedly a long text",
					SourceUrl:    "https://github.com/...",
					Maturity:     6,
				},
				moduleVersions: []string{
					"1.0.1",
					"1.0.2",
				},
			},
			want: &moduleItem{
				Organization: "cie",
				Name:         "test-module",
				Provider:     "aws",
				Description:  "This is the description for the module it is supposedly a long text",
				SourceUrl:    "https://github.com/...",
				Maturity:     "DEPRECATED",
				Versions: []string{
					"1.0.1",
					"1.0.2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createModuleMetadataResponse(tt.args.moduleMetadata, tt.args.moduleVersions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createModuleMetadataResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createProviderMetadataResponse(t *testing.T) {
	type args struct {
		providerMetadata *providerServices.ListProviderItem
		providerVersions []string
	}
	tests := []struct {
		name string
		args args
		want *providerItem
	}{
		{
			name: "Provider Metadata Response",
			args: args{
				providerMetadata: &providerServices.ListProviderItem{
					Organization:  "cie",
					Name:          "test-provider",
					Description:   "This is the description for the provider it is supposedly a long text",
					SourceRepoUrl: "https://github.com/...",
					Maturity:      6,
				},
				providerVersions: []string{
					"1.0.3",
					"1.0.4",
				},
			},
			want: &providerItem{
				Organization:  "cie",
				Name:          "test-provider",
				Description:   "This is the description for the provider it is supposedly a long text",
				SourceRepoUrl: "https://github.com/...",
				Maturity:      "DEPRECATED",
				Versions: []string{
					"1.0.3",
					"1.0.4",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createProviderMetadataResponse(tt.args.providerMetadata, tt.args.providerVersions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createProviderMetadataResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
