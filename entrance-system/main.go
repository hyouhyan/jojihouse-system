package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	Env_load()
	// postgreSQLに接続
	db, err := sql.Open("postgres", "host="+os.Getenv("POSTGRES_HOST")+" port="+os.Getenv("POSTGRES_PORT")+" user="+os.Getenv("POSTGRES_USER")+" password="+os.Getenv("POSTGRES_PASSWORD")+" dbname="+os.Getenv("POSTGRES_DB")+" sslmode=disable")

	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	// データベースに接続できるか確認
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully connected!")
}
