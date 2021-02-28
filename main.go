package main

import (
	"github.com/shawnzxx/bookstore_users-api/app"
	"github.com/shawnzxx/bookstore_users-api/infrastructure/mysql/users_db"
)

func main() {
	users_db.SetupDatabase()
	app.StartApplication()
}
