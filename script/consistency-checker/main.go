package main

import (
	"fmt"
	"jojihouse-system-consistency-checker/database"
	"jojihouse-system-consistency-checker/model"

	"github.com/jmoiron/sqlx"
)

func main() {
	// DBへ接続
	database.ConnectPostgres()
	defer database.ClosePostgres()

	database.ConnectMongo()
	defer database.CloseMongo()

	users, err := GetAllUsers(database.PostgresDB)
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		return
	}

	for _, user := range users {
		fmt.Println(*user.Name)
	}
}

func GetAllUsers(db *sqlx.DB) ([]model.User, error) {
	var users []model.User
	err := db.Select(&users, "SELECT * FROM users ORDER BY number")

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}


