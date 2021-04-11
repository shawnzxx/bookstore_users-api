package users_db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

const (
	mysqlUsersUsername = "mysql_users_username"
	mysqlUsersPassword = "mysql_users_password"
	mysqlUsersHost     = "mysql_users_host"
	mysqlUsersSchema   = "mysql_users_schema"
)

//DbConn - export connected DbConn object
var (
	DbContext *sql.DB
	username  = os.Getenv(mysqlUsersUsername)
	password  = os.Getenv(mysqlUsersPassword)
	host      = os.Getenv(mysqlUsersHost)
	schema    = os.Getenv(mysqlUsersSchema)
)

//SetupDatabase - connect to the db
func SetupDatabase() {
	var err error
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, schema)
	DbContext, err = sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	// test ping db after connect to the database
	if err = DbContext.Ping(); err != nil {
		panic(err)
	}
	DbContext.SetMaxOpenConns(3)
	DbContext.SetMaxIdleConns(3)
	DbContext.SetConnMaxLifetime(60 * time.Second)

	log.Println("database successfully configured")
}
