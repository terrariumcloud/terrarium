package main

import (
	"context"
	"terrarium-grpc-gateway/internal/services"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
)

type DependencyService struct {
	services.UnimplementedDependencyResolverServer
	db dynamodbiface.DynamoDBAPI
}

func (s *DependencyService) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	return nil, nil
}

func (s *DependencyService) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	return nil, nil
}

func (s *DependencyService) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server services.DependencyResolver_RetrieveContainerDependenciesServer) error {
	return nil
}

func (s *DependencyService) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server services.DependencyResolver_RetrieveModuleDependenciesServer) error {
	return nil
}
