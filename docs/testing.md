# Unit testing

This code base contains unit tests that cover services. To run all tests execute:

```bash
go test -v ./...
```

Before pushing commits, please make sure you ran all tests and that they pass. Also, if you've contributed to this code base, be sure to add/update unit tests.

# GRPC service testing

For testing live GRPC services you can use [grpcurl](https://github.com/fullstorydev/grpcurl) tool to send requests.

It can be installed with go cli:

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

Here are some examples for testing terrarium gateway:

```bash
grpcurl -d '{"api_key": "123", "name": "module1", "description": "some description", "source_url": "http://my.dot.com", "maturity": "BETA"}' -plaintext -proto .pb/terrarium/module/module.proto 10.43.191.121:8080 terrarium.module.Publisher/Register

grpcurl -d '{"api_key": "123", "module": { "name": "module1", "version": "v1.0.0" }}' -plaintext -proto .pb/terrarium/module/module.proto 10.43.191.121:8080 terrarium.module.Publisher/BeginVersion

grpcurl -d '{"session_key": "7b15ef52-4055-4c94-8693-af8819c78ee1", "action": 0 }' -plaintext -proto .pb/terrarium/module/module.proto 10.43.191.121:8080 terrarium.module.Publisher/EndVersion

grpcurl -d '{"session_key": "a6d7792f-bb4e-493d-97f1-e5fa8c8ca43f", "action": 1 }' -plaintext -proto .pb/terrarium/module/module.proto 10.43.191.121:8080 terrarium.module.Publisher/EndVersion


grpcurl -d '{"session_key": "01e1891f-1f7e-46b8-b2f5-8b5f4081934a", "modules": [{ "name": "module1", "version": "v1.0.0" }] }' -plaintext -proto .pb/terrarium/module/module.proto 10.43.191.121:8080 terrarium.module.Publisher/RegisterModuleDependencies

grpcurl -d '{"session_key": "01e1891f-1f7e-46b8-b2f5-8b5f4081934a", "container_image_references": ["image1", "image2"] }' -plaintext -proto .pb/terrarium/module/module.proto 10.43.191.121:8080 terrarium.module.Publisher/RegisterContainerDependencies


grpcurl -d '{"session_key": "01e1891f-1f7e-46b8-b2f5-8b5f4081934a", "zip_data_chunk": "VGhpcyBpcyBhIHRlc3QgZmlsZQo=" }' -plaintext -proto .pb/terrarium/module/module.proto 10.43.191.121:8080 terrarium.module.Publisher/UploadSourceZip
```
