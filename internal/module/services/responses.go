package services

import terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

var (
	// Gateway service
	FailedToConnectToRegistrar          = Error("Unable to connect to Registrar service.")
	FailedToExecuteRegister             = Error("Failed to execute Register module.")
	FailedToConnectToVersionManager     = Error("Unable to connect to module Version service.")
	UnknownVersionManagerAction         = Error("Unknown Version manager action requested.")
	ArchiveUploaded                     = Ok("Archive uploaded successfully.")
	ArchiveUploadFailed                 = Error("Upload failed.")
	FailedToConnectToDependencyResolver = Error("Unable to connect to module Dependency service.")

	// Registrar service
	ModuleRegistered    = Ok("Module registered successfully.")
	ModuleNotRegistered = Error("Failed to register module.")
	MarshalModuleError  = Error("Failed to marshal module.")

	// VersionManager service
	SessionKeyNotRemoved = Error("Failed to remove session key.")
	VersionPublished     = Ok("Version published")
	VersionAborted       = Ok("Version aborted.")
	NotFound             = Error("Could not find item.")
	FailedToMarshal      = Error("Failed to marshal record.")
	FailedToUnmarshal    = Error("Failed to unmarshal record.")

	// DependencyResolver service
	RegisterModuleDependenciesFailed    = Error("Failed to register module dependencies.")
	ModuleDependenciesRegistered        = Ok("Module dependencies successfully registered.")
	RegisterContainerDependenciesFailed = Error("Failed to register container dependencies.")
	ContainerDependenciesRegistered     = Ok("Container dependencies successfully registered.")

	// Storage service
	ZipUploaded     = Ok("Source zip uploaded successfully.")
	ZipUploadFailed = Error("Failed to upload source zip.")
)

func Ok(message string) *terrarium.TransactionStatusResponse {
	return &terrarium.TransactionStatusResponse{
		Status:        terrarium.Status_OK,
		StatusMessage: message,
	}
}

func Error(message string) *terrarium.TransactionStatusResponse {
	return &terrarium.TransactionStatusResponse{
		Status:        terrarium.Status_UNKNOWN_ERROR,
		StatusMessage: message,
	}
}
