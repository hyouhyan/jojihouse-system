package database

import (
	"fmt"
	"jojihouse-system/internal/config"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var PostgresDB *sqlx.DB

func ConnectPostgres() {
	config.Env_load()

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

	log.Println("Successfully connected to postgres!")

	PostgresDB = db
}

func ClosePostgres() {
	if PostgresDB != nil {
		PostgresDB.Close()
		log.Println("Successfully closed postgres!")
	}
}
