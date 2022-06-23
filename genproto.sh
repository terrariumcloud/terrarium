protoc --go_out=. --go_opt=module=github.com/terrariumcloud/terrarium-grpc-gateway \
    --go-grpc_out=. --go-grpc_opt=module=github.com/terrariumcloud/terrarium-grpc-gateway \
    pkg/terrarium/module.proto internal/services/module_create.proto internal/services/module_session_manager.proto


