package users_db

import (
	"database/sql"
	"log"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

//DbConn - export connected DbConn object
var dbContext *sql.DB

//SetupDatabase - connect to the db
func SetupDatabase() {
	var err error
	dbContext, err = sql.Open("mysql", "root:Passw0rd123!@tcp(127.0.0.1:3306)/users_db?charset=utf8")
	if err != nil {
		panic(err)
	}
	if err = dbContext.Ping(); err != nil {
		panic(err)
	}

	dbContext.SetMaxOpenConns(3)
	dbContext.SetMaxIdleConns(3)
	dbContext.SetConnMaxLifetime(60 * time.Second)

	log.Println("database successfully configured")
}
