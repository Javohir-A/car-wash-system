package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"booking-service/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(config *config.Config) (*mongo.Database, error) {

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		config.MongoConfig.User,
		config.MongoConfig.Password,
		config.MongoConfig.Host,
		config.MongoConfig.Port)
	// uri := fmt.Sprintf("mongodb://%s:%s",
	// 	config.MongoConfig.User,
	// 	config.MongoConfig.Password,
	// 	config.MongoConfig.Host,
	// 	config.MongoConfig.Port)
	// log.Println("Mongo started on localhost. Note: if deployment fails uncomment first connection code.")
	fmt.Println(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(config.MongoConfig.DBname)

	log.Printf("--------------------------- Connected to the database %s --------------------------------\n", config.MongoConfig.DBname)

	return db, nil
}
