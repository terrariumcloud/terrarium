module terrarium-grpc-gateway

go 1.18

require (
	github.com/terrariumcloud/terrarium-grpc-gateway v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/aws/aws-sdk-go v1.44.44
	github.com/google/uuid v1.3.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace github.com/terrariumcloud/terrarium-grpc-gateway => ./
