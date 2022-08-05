protoc --go_out=. --go_opt=module=github.com/terrariumcloud/terrarium-grpc-gateway \
    --go-grpc_out=. --go-grpc_opt=module=github.com/terrariumcloud/terrarium-grpc-gateway \
    pb/terrarium/module/module.proto \
    pb/terrarium/module/services/registrar.proto \
    pb/terrarium/module/services/version_manager.proto \
    pb/terrarium/module/services/dependency_resolver.proto \
    pb/terrarium/module/services/storage.proto
