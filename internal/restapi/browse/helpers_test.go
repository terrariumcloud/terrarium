package browse

import (
	"reflect"
	"testing"

	"github.com/terrariumcloud/terrarium/internal/module/services"
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
			name: "Module Metadata Reponse",
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
