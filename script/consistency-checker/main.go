package main

import "jojihouse-system-consistency-checker/database"

func main() {
	// DBへ接続
	database.ConnectPostgres()
	defer database.ClosePostgres()

	database.ConnectMongo()
	defer database.CloseMongo()
}
