package users

import (
	"strings"

	"github.com/shawnzxx/bookstore_utils-go/rest_errors"
)

const (
	StatusActive = "active"
)

// domain DTO object
type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

// Users define as array of User, later in other function can user Users directly
type Users []User

func (user *User) Validate() *rest_errors.RestErr {
	user.FirstName = strings.TrimSpace(strings.ToLower(user.FirstName))
	user.LastName = strings.TrimSpace(strings.ToLower(user.LastName))
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return rest_errors.NewBadRequestError("invalid email address")
	}
	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return rest_errors.NewBadRequestError("invalid password")
	}
	return nil
}
