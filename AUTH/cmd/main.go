package main

import (
	"auth-service/cmd/app"
	"auth-service/cmd/repository"
	"auth-service/config"
	p "auth-service/pkg/postgres"
	r "auth-service/pkg/redis"
	"fmt"
	"log"

	"log/slog"
)

func main() {
	cnf := config.NewConfig()
	if err := cnf.Load(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cnf.Database)
	db, err := p.ConnectDB(cnf.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rClient := r.ConnectDB(&cnf.Redis)

	authService, err := repository.SetupAuthService(db, rClient, cnf, *slog.New(slog.Default().Handler()))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("grpc server is listening on %s", cnf.Server.Host+":"+cnf.Server.Port)

	if err := app.Run(&cnf.Server, authService); err != nil {
		log.Fatal(err)
	}
}
