/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	services2 "github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/module/services/dependency_manager"
	"github.com/terrariumcloud/terrarium/internal/module/services/gateway"
	"github.com/terrariumcloud/terrarium/internal/module/services/registrar"
	storage2 "github.com/terrariumcloud/terrarium/internal/module/services/storage"
	"github.com/terrariumcloud/terrarium/internal/module/services/tag_manager"
	"github.com/terrariumcloud/terrarium/internal/module/services/version_manager"
	"github.com/terrariumcloud/terrarium/internal/release/services/release"
	"github.com/terrariumcloud/terrarium/internal/restapi/browse"
	modulesv1 "github.com/terrariumcloud/terrarium/internal/restapi/modules/v1"
	"github.com/terrariumcloud/terrarium/internal/storage"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

const (
	allInOneInternalEndpoint    = "localhost:30001"
	allInOneGrpcGatewayEndpoint = "0.0.0.0:3001"
	allInOneHTTPEndpoint        = "0.0.0.0:8080"
)

type allInOneRestHandler struct {
	router *mux.Router
}

func (a allInOneRestHandler) GetHttpHandler(mountPath string) http.Handler {
	return a.router
}

// allInOneCmd represents the allInOne command
var allInOneCmd = &cobra.Command{
	Use:   "all-in-one",
	Short: "Runs all the services in a single command.",
	Long:  `This runs all the micro-services as part of a single process, useful for developing and for trying out Terrarium.`,
	Run: func(cmd *cobra.Command, args []string) {
		dependencyServiceServer := &dependency_manager.DependencyManagerService{
			Db:              storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
			ModuleTable:     dependency_manager.ModuleDependenciesTableName,
			ModuleSchema:    dependency_manager.GetDependenciesSchema(dependency_manager.ModuleDependenciesTableName),
			ContainerTable:  dependency_manager.ContainerDependenciesTableName,
			ContainerSchema: dependency_manager.GetDependenciesSchema(dependency_manager.ContainerDependenciesTableName),
		}

		registrarServiceServer := &registrar.RegistrarService{
			Db:     storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
			Table:  registrar.RegistrarTableName,
			Schema: registrar.GetModulesSchema(registrar.RegistrarTableName),
		}

		storageServiceServer := &storage2.StorageService{
			Client:     storage.NewS3Client(awsAccessKey, awsSecretKey, awsRegion),
			BucketName: storage2.BucketName,
			Region:     awsRegion,
		}

		tagManagerServer := &tag_manager.TagManagerService{
			Db:     storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
			Table:  tag_manager.TagTableName,
			Schema: tag_manager.GetTagsSchema(tag_manager.TagTableName),
		}

		releaseServiceServer := &release.ReleaseService{
			Db:     storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
			Table:  release.ReleaseTableName,
			Schema: release.GetReleaseSchema(release.ReleaseTableName),
		}

		versionManagerServer := &version_manager.VersionManagerService{
			Db:             storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
			Table:          version_manager.VersionsTableName,
			Schema:         version_manager.GetModuleVersionsSchema(version_manager.VersionsTableName),
			ReleaseService: release.NewPublisherGrpcClient(allInOneInternalEndpoint),
		}

		services := []services2.Service{
			dependencyServiceServer,
			registrarServiceServer,
			storageServiceServer,
			tagManagerServer,
			releaseServiceServer,
			versionManagerServer,
		}

		otelShutdown := initOpenTelemetry("all-in-one")
		defer otelShutdown()

		startAllInOneGrpcServices(services, allInOneInternalEndpoint)

		gatewayServer := gateway.New(registrar.NewRegistrarGrpcClient(allInOneInternalEndpoint),
			tag_manager.NewTagManagerGrpcClient(allInOneInternalEndpoint),
			version_manager.NewVersionManagerGrpcClient(allInOneInternalEndpoint),
			storage2.NewStorageGrpcClient(allInOneInternalEndpoint),
			dependency_manager.NewDependencyManagerGrpcClient(allInOneInternalEndpoint),
			release.NewPublisherGrpcClient(allInOneInternalEndpoint),
		)
		startAllInOneGrpcServices([]services2.Service{gatewayServer}, allInOneGrpcGatewayEndpoint)

		restAPIServer := browse.New(registrar.NewRegistrarGrpcClient(allInOneInternalEndpoint),
			version_manager.NewVersionManagerGrpcClient(allInOneInternalEndpoint),
			release.NewBrowseGrpcClient(allInOneInternalEndpoint))

		modulesAPIServer := modulesv1.New(version_manager.NewVersionManagerGrpcClient(allInOneInternalEndpoint), storage2.NewStorageGrpcClient(allInOneInternalEndpoint))

		router := mux.NewRouter()
		router.PathPrefix("/modules").Handler(modulesAPIServer.GetHttpHandler("/modules"))
		router.PathPrefix("/").Handler(restAPIServer.GetHttpHandler(""))

		endpoint = allInOneHTTPEndpoint
		startRESTAPIService("browse", "", allInOneRestHandler{router: router})

	},
}

func init() {
	rootCmd.AddCommand(allInOneCmd)
	allInOneCmd.Flags().StringVar(&storage2.BucketName, "storage-bucket", storage2.DefaultBucketName, "Module bucket name")
	allInOneCmd.Flags().StringVar(&version_manager.VersionsTableName, "version-table", version_manager.DefaultVersionsTableName, "Module versions table name")
	allInOneCmd.Flags().StringVar(&tag_manager.TagTableName, "tag-table", tag_manager.DefaultTagTableName, "Module tags table name")
	allInOneCmd.Flags().StringVar(&release.ReleaseTableName, "release-table", release.DefaultReleaseTableName, "Releases table name")
	allInOneCmd.Flags().StringVar(&registrar.RegistrarTableName, "registrar-table", registrar.DefaultRegistrarTableName, "Module Registrar table name")
	allInOneCmd.Flags().StringVar(&dependency_manager.ModuleDependenciesTableName, "module-dependencies-table", dependency_manager.DefaultModuleDependenciesTableName, "Module dependencies table name")
	allInOneCmd.Flags().StringVar(&dependency_manager.ContainerDependenciesTableName, "container-dependencies-table", dependency_manager.DefaultContainerDependenciesTableName, "Module container dependencies table name")
}

func startAllInOneGrpcServices(services []services2.Service, endpoint string) {
	listener, err := net.Listen("tcp4", endpoint)
	if err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	for _, service := range services {
		if err := service.RegisterWithServer(grpcServer); err != nil {
			log.Fatalf("Failed to start: %v", err)
		}
	}

	go func() {
		log.Printf("Listening at %s", endpoint)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed: %v", err)
		}
	}()
}
