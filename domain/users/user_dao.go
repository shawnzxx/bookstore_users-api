package users

import (
	"fmt"
	"strings"

	"github.com/shawnzxx/bookstore_users-api/infrastructure/mysql/users_db"
	"github.com/shawnzxx/bookstore_users-api/utils/date_utils"
	"github.com/shawnzxx/bookstore_users-api/utils/errors"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	errorNoRows      = "no rows in result set"
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
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
		}
		return errors.NewInternalServerError(
			fmt.Sprintf("error when trying to get user %d: %s", user.Id, err.Error()))
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
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequestError(
				fmt.Sprintf("email %s alreadt excists", user.Email))
		}
		return errors.NewInternalServerError(
			fmt.Sprintf("error when tyring to save user: %s", err.Error()))
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(
			fmt.Sprintf("error when tyring to save user: %s", err.Error()))
	}
	user.Id = userId
	return nil
}
