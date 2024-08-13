package main

import (
	"booking-service/genproto/providers"
	"booking-service/genproto/services"
	"booking-service/internal/service"
	mongostorage "booking-service/internal/storage/mongo"
	"booking-service/pkg/mongo"
	"log"
	"log/slog"
	"net"

	logg "github.com/labstack/gommon/log"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.Default().Handler())
	lis, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		slog.Error(err.Error())
		return
	}

	grpcServer := grpc.NewServer()

	mongoClient := mongo.ConnectDB()
	mongoDb := mongoClient.Database("booking-service")

	providerMongoStorage := mongostorage.NewProviderMongoStorage(mongoDb.Collection("providers"))

	providerService := service.NewProviderService(logger, providerMongoStorage)

	serviceStorage := mongostorage.NewServicesMongoStorage(mongoDb.Collection("services"))
	servicesService := service.NewServicesService(logger, serviceStorage)

	services.RegisterServicesServer(grpcServer, servicesService)
	providers.RegisterProvidersServer(grpcServer, providerService)

	logg.Infof("serving on 50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
