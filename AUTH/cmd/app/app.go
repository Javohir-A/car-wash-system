package app

import (
	"auth-service/config"
	"auth-service/genprotos/auth"
	"auth-service/internal/service"

	"net"

	"google.golang.org/grpc"
)

func Run(cnf *config.ServerConfig, autService *service.AuthServiceImpl) error {

	listener, err := net.Listen("tcp", cnf.Host+":"+cnf.Port)
	if err != nil {
		return err
	}
	
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, autService)
	auth.RegisterUserManagementServiceServer(grpcServer, autService)

	return grpcServer.Serve(listener)
}
