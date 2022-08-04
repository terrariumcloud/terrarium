#!/bin/bash

rm -rf ./vendor terrarium
go fmt ./...
go mod vendor
go build -o terrarium -mod vendor
