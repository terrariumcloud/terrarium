package browse

import (
	"google.golang.org/grpc"
	//pb "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
	// "fmt"
	// "github.com/gorilla/mux"
	//"google.golang.org/grpc"
	"net/http"
)

func getModuleFromRequest(r *http.Request) string {
	return "to remove err"
}

func createModulesResponse(modules []string) *ModuleResponse {
	var structuredVersion []*ModuleItem

	for _, module := range modules {
		structuredVersion = append(structuredVersion, &ModuleItem{
			Name:        name,
			Provider:    provider,
			Description: description,
			SourceUrl:   source_url,
			Maturity:    maturity,
		})
	}
	return &ModuleResponse{
		Modules: []*Modules{
			{
				Modules: structuredVersion,
			},
		},
	}
}

func closeClient(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {

	}
}
