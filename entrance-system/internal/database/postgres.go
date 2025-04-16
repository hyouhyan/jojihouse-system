package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}
}

var PostgresDB *sqlx.DB

func ConnectPostgres() {
	Env_load()
	dst := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	// postgreSQLに接続
	db, err := sqlx.Open("postgres", dst)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}

	// データベースに接続できるか確認
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping: %v", err)
	}

	fmt.Println("Successfully connected to postgres!")

	PostgresDB = db
}

func ClosePostgres() {
	if PostgresDB != nil {
		PostgresDB.Close()
		fmt.Println("Successfully closed postgres!")
	}
}
