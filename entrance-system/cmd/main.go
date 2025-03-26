package main

import (
	"jojihouse-entrance-system/internal/database"
)

func main() {
	database.ConnectPostgres()
	defer database.ClosePostgres()

	database.ConnectMongo()
	defer database.CloseMongo()

	test(database.PostgresDB, database.MongoDB)
}
