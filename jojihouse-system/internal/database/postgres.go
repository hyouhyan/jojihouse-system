package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var PostgresDB *sqlx.DB

func ConnectPostgres() {
	dst := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	// 10回までリトライ
	for i := 0; i < 10; i++ {
		if i > 0 {
			log.Println("Retrying connection to Postgres...")
			// 5秒待機
			time.Sleep(5 * time.Second)
		}

		// postgreSQLに接続
		db, err := sqlx.Open("postgres", dst)
		if err != nil {
			log.Printf("Failed to open a DB connection: %v", err)
			continue
		}

		// データベースに接続できるか確認
		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping: %v", err)
			continue
		}

		// 接続成功で代入
		PostgresDB = db
		break
	}

	// nilの場合は失敗
	if PostgresDB == nil {
		log.Fatal("Failed to connect to Postgres after multiple attempts")
	}

	log.Println("Successfully connected to postgres!")
}

func ClosePostgres() {
	if PostgresDB != nil {
		PostgresDB.Close()
		log.Println("Successfully closed postgres!")
	}
}
