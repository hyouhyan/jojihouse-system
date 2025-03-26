package main

import (
	"jojihouse-entrance-system/internal/database"
)

func main() {
	database.Connect()
	defer database.Close()

	test(database.DB)
}
