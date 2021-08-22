package users

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shawnzxx/bookstore_users-api/infrastructure/mysql/users_db"
	"github.com/shawnzxx/bookstore_users-api/utils/mysql_utils"
	"github.com/shawnzxx/bookstore_utils-go/logger"
	"github.com/shawnzxx/bookstore_utils-go/rest_errors"
)

// database access layer
const (
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?,?,?,?,?,?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=? AND status=?;"
	queryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=? AND status=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?;"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?;"
)

func (user *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.DbContext.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id, StatusActive)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("no records found")
		}
		logger.Error("error when trying to get user by id", getErr)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	return nil
}

func (user *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.DbContext.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	//as long as we didn't detect any error then defer close connection
	//this will be called when statement no need anymore
	defer stmt.Close()
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.DbContext.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id, StatusActive)
	if updateErr != nil {
		logger.Error("error when trying to update user by id", updateErr)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}
	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.DbContext.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	defer stmt.Close()

	_, deleteErr := stmt.Exec(user.Id)
	if deleteErr != nil {
		logger.Error("error when trying to delete user by id", deleteErr)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := users_db.DbContext.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	//whennever you work with RDB, you must rememeber to close both the statement you have created and the rows you have opened
	defer rows.Close() //after we detect there is no error then we defer the close when scan finished

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to scan user row into user struct", err)
			return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := users_db.DbContext.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	return nil
}
