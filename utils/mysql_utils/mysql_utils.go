package mysql_utils

import (
	"strings"

	"github.com/go-sql-driver/mysql"

	errors "github.com/sovernut/bookstore_users-api/utils/error"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError("no record found by given id")
		}
		return errors.NewInternalServerError("error parsing db response")

	}

	switch sqlErr.Number {
	case 1062:
		return errors.NewBadRequestError("invalid data")
	}

	return errors.NewInternalServerError("error procesing request")
}
