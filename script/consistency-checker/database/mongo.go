package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database
var mongoClient *mongo.Client

func ConnectMongo() {
	dst := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		os.Getenv("MONGO_USER"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
	)

	// 10回までリトライ
	for i := 0; i < 10; i++ {
		if i > 0 {
			log.Println("Retrying connection to Mongo...")
			// 5秒待機
			time.Sleep(5 * time.Second)
		}

		// MongoDBに接続
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dst))
		if err != nil {
			log.Printf("Failed to open a DB connection: %v", err)
			continue
		}

		// データベースに接続できるか確認
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Printf("Failed to ping: %v", err)
			continue
		}

		// 接続成功で代入
		mongoClient = client
		break
	}

	// nilの場合は失敗
	if mongoClient == nil {
		log.Fatal("Failed to connect to Mongo after multiple attempts")
	}

	log.Println("Successfully connected to mongo!")

	// Mongoのデータベース取得
	MongoDB = mongoClient.Database(os.Getenv("MONGO_DB"))

	// nilの場合は失敗
	if MongoDB == nil {
		log.Fatal("Failed to get MongoDB instance: ", os.Getenv("MONGO_DB"))
	}

	log.Print("Successfully got MongoDB instance: ", MongoDB.Name())
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
