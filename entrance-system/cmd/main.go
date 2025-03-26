package main

import (
	"jojihouse-entrance-system/internal/database"
)

func main() {
	database.ConnectPostgres()
	defer database.ClosePostgres()

	// test(database.PostgresDB)

	database.ConnectMongo()
	defer database.CloseMongo()
}
