package users

import (
	"fmt"
	"github.com/shawnzxx/bookstore_users-api/infrastructure/mysql/users_db"
	"github.com/shawnzxx/bookstore_users-api/utils/date_utils"
	"github.com/shawnzxx/bookstore_users-api/utils/errors"
	"strings"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?,?,?,?)"
)

var dbUsers = make(map[int64]*User)

func (user *User) Get() *errors.RestErr {
	result := dbUsers[user.Id]
	if result == nil {
		return errors.NewBadRequestError(fmt.Sprintf("user %d not found", user.Id))
	}
	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.DbContext.Prepare(queryInsertUser)
	if err  != nil{
		return errors.NewInternalServerError(err.Error())
	}
	//as long as we didn't detect any error then defer close connection
	//this will be called when statement no need anymore
	defer stmt.Close()
	user.DateCreated = date_utils.GetNowString()
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil{
		if strings.Contains(err.Error(), indexUniqueEmail){
			return errors.NewBadRequestError(
				fmt.Sprintf("email %s alreadt excists", user.Email))
		}
		return errors.NewInternalServerError(
			fmt.Sprintf("error when tyring to save user: %s", err.Error()))
	}
	userId, err := insertResult.LastInsertId()
	if err != nil{
		return errors.NewInternalServerError(
			fmt.Sprintf("error when tyring to save user: %s", err.Error()))
	}
	user.Id = userId
	return nil
}
