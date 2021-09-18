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
	MysqlUsersUsername = "MYSQL_USERS_USERNAME"
	MysqlUsersPassword = "MYSQL_USERS_PASSWORD"
	MysqlUsersHost     = "MYSQL_USERS_HOST"
	MysqlUsersSchema   = "MYSQL_USERS_SCHEMA"
)

//DbConn - export connected DbConn object
var (
	DbContext *sql.DB
	username  = os.Getenv(MysqlUsersUsername)
	password  = os.Getenv(MysqlUsersPassword)
	host      = os.Getenv(MysqlUsersHost)
	schema    = os.Getenv(MysqlUsersSchema)
)

//SetupDatabase - connect to the db
func SetupDatabase() {
	var err error
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, schema)
	DbContext, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Println(err.Error())
	}
	//Retry Database Connect With Docker and Go
	//http://www.matthiassommer.it/programming/docker-compose-retry-database-connect-with-docker-and-go/
	retryCount := 30
	for {
		err := DbContext.Ping()
		if err != nil {
			if retryCount == 0 {
				log.Fatalf("Not able to establish connection to host %s database %s", host, schema)
			}
			log.Printf(fmt.Sprintf("Could not connect to database. Wait 2 seconds. %d retries left...", retryCount))
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	// test ping db after connect to the database
	//if err = DbContext.Ping(); err != nil {
	//	panic(err)
	//}
	DbContext.SetMaxOpenConns(3)
	DbContext.SetMaxIdleConns(3)
	DbContext.SetConnMaxLifetime(60 * time.Second)

	log.Println("database successfully configured")
}
