module github.com/terrariumcloud/terrarium/tools/cli

go 1.21.4

require (
	github.com/spf13/cobra v1.8.0
	github.com/terrariumcloud/terrarium v0.0.69
	google.golang.org/grpc v1.59.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231030173426-d783a09b4405 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace github.com/terrariumcloud/terrarium/pkg/terrarium/module v0.0.69 => ../../pkg/terrarium/module
