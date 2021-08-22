package mysql_utils

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/shawnzxx/bookstore_utils-go/rest_errors"
)

const (
	ErrorNoRows = "no rows in result set"
)

func ParseError(err error) *rest_errors.RestErr {
	// try to convert error into sqlError
	sqlErr, ok := err.(*mysql.MySQLError)
	// if it is not sql error
	if !ok {
		// check if error message contain below sentense
		if strings.Contains(err.Error(), ErrorNoRows) {
			return rest_errors.NewNotFoundError("no record matching with given id")
		}
		// other unknown error
		return rest_errors.NewInternalServerError("error parsing database response", errors.New("database error"))
	}
	// if it is sql error
	switch sqlErr.Number {
	case 1062:
		return rest_errors.NewBadRequestError("duplicated data")
	}
	return rest_errors.NewInternalServerError("error processing request", errors.New("database error"))
}
