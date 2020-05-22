package mysql_utils

import (
	"strings"

	"github.com/go-sql-driver/mysql"

	builtInError "errors"

	errors "github.com/sovernut/bookstore_utils-go/rest_errors"
)

const (
	ErrorNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return errors.NewNotFoundError("no record found by given id")
		}
		return errors.NewInternalServerError("error while trying to ParseError", builtInError.New("error parsing db response"))

	}

	switch sqlErr.Number {
	case 1062:
		return errors.NewBadRequestError("invalid data")
	}

	return errors.NewInternalServerError("error while trying to ParseError", builtInError.New("error processing request"))

}
