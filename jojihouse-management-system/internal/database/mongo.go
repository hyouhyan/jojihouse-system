package database

import (
	"context"
	"fmt"
	"jojihouse-management-system/internal/config"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database
var mongoClient *mongo.Client

func ConnectMongo() {
	config.Env_load()

	dst := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		os.Getenv("MONGO_USER"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
	)
	// MongoDBに接続
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dst))
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}

	// データベースに接続できるか確認
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping: %v", err)
	}

	log.Println("Successfully connected to mongo!")

	mongoClient = client
	MongoDB = client.Database(os.Getenv("MONGO_DB"))
}

func CloseMongo() {
	if MongoDB != nil {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Failed to close: %v", err)
		}
		log.Println("Successfully closed mongo!")
	}
}
