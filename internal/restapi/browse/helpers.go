package browse

import (
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"

	"google.golang.org/grpc"
)

func createModulesResponse(modules []*services.ModuleMetadata) *ModuleResponse {
	return &ModuleResponse{
		Modules: modules,
	}
}

func closeClient(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {

	}
}
