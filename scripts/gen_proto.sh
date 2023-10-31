# TODO: need to validate the appropriate version of protoc is available 
protoc --go_out=. --go_opt=module=github.com/terrariumcloud/terrarium \
    --go-grpc_out=. --go-grpc_opt=module=github.com/terrariumcloud/terrarium \
    pb/terrarium/module/module.proto \
    pb/terrarium/module/services/registrar.proto \
    pb/terrarium/module/services/version_manager.proto \
    pb/terrarium/module/services/dependency_manager.proto \
    pb/terrarium/module/services/storage.proto \
    pb/terrarium/release/release.proto \
    pb/terrarium/release/services/release.proto \
    pb/terrarium/common/paging.proto

