#!/bin/bash

rm -rf ./vendor terrarium
go fmt ./...
go mod tidy
go mod vendor
go build -o terrarium -mod vendor
