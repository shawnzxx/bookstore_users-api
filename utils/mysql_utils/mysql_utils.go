package mysql_utils

import (
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/shawnzxx/bookstore_users-api/utils/errors"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	// try to convert error into sqlError
	sqlErr, ok := err.(*mysql.MySQLError)
	// if it is not sql error
	if !ok {
		// check if error message contain below sentense
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError("no record matching with given id")
		}
		// other unknown error
		return errors.NewInternalServerError("error parsing database response")
	}
	// if it is sql error
	switch sqlErr.Number {
	case 1062:
		return errors.NewBadRequestError("duplicated data")
	}
	return errors.NewInternalServerError("error processing request")
}
