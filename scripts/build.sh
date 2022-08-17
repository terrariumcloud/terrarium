#!/bin/bash

go fmt ./...
go mod tidy
go build -o terrarium -mod vendor
