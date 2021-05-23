package services

import (
	"github.com/shawnzxx/bookstore_users-api/domain/users"
	"github.com/shawnzxx/bookstore_users-api/utils/date_utils"
	"github.com/shawnzxx/bookstore_users-api/utils/errors"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDateTimeString()
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUser(userId int64) (*users.User, *errors.RestErr) {
	var user = &users.User{
		Id: userId,
	}
	if err := user.Get(); err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {
	// get user from db by Id
	curUser, err := GetUser(user.Id)
	if err != nil {
		return nil, err
	}
	// validate pass in body
	if err := user.Validate(); err != nil {
		return nil, err
	}
	//if is partial only update the fileds not empty
	if isPartial {
		if user.FirstName != "" {
			curUser.FirstName = user.FirstName
		}
		if user.LastName != "" {
			curUser.LastName = user.LastName
		}
		if user.Email != "" {
			curUser.Email = user.Email
		}
	} else {
		curUser.FirstName = user.FirstName
		curUser.LastName = user.LastName
		curUser.Email = user.Email
	}

	if err := curUser.Update(); err != nil {
		return nil, err
	}
	return curUser, nil
}

func DeleteUser(userId int64) *errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func Search(status string) ([]users.User, *errors.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}
