package users

import (
	"fmt"

	"github.com/shawnzxx/bookstore_users-api/utils/date"
	"github.com/shawnzxx/bookstore_users-api/utils/errors"
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
	current := dbUsers[user.Id]
	if current != nil {
		if current.Email == user.Email {
			return errors.NewBadRequestError(fmt.Sprintf("email %s is already registered", user.Email))
		}
		return errors.NewBadRequestError(fmt.Sprintf("user %d is already excist", user.Id))
	}

	user.DateCreated = date.GetNowString()

	dbUsers[user.Id] = user
	return nil
}
