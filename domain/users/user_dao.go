package users

import (
	"github.com/shawnzxx/bookstore_users-api/infrastructure/mysql/users_db"
	"github.com/shawnzxx/bookstore_users-api/utils/date_utils"
	"github.com/shawnzxx/bookstore_users-api/utils/errors"
	"github.com/shawnzxx/bookstore_users-api/utils/mysql_utils"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?,?,?,?)"
	queryGetUser     = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.DbContext.Prepare(queryGetUser)
	if err != nil {
		return errors.NewBadRequestError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); getErr != nil {
		return mysql_utils.ParseError(getErr)
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.DbContext.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	//as long as we didn't detect any error then defer close connection
	//this will be called when statement no need anymore
	defer stmt.Close()
	user.DateCreated = date_utils.GetNowString()
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	user.Id = userId
	return nil
}
