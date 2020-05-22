package users

import (
	"fmt"
	"strings"

	"github.com/sovernut/bookstore_users-api/datasources/mysql/users_db"
	"github.com/sovernut/bookstore_users-api/logger"
	"github.com/sovernut/bookstore_users-api/utils/mysql_utils"

	builtInError "errors"

	errors "github.com/sovernut/bookstore_utils-go/rest_errors"
)

const (
	queryInsertUser             = "INSERT INTO users(first_name,last_name,email,date_created,status,password) values (?,?,?,?,?,?);"
	queryGetUser                = "SELECT  id,first_name,last_name,email,date_created,status FROM users WHERE id=?;"
	queryUpdateUser             = "UPDATE users SET first_name=?,last_name=?,email=? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus       = "SELECT  id,first_name,last_name,email,date_created,status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status from users WHERE email = ? AND password = ? AND status = ?;"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("error when trying to prepare get user statement", builtInError.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if err = result.Scan(&user.Id, &user.FirstName,
		&user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		logger.Error("error when trying to get user by id", err)
		return errors.NewInternalServerError("error when trying to get user by id", builtInError.New("database error"))
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("error when trying to prepare save user statement", builtInError.New("database error"))

	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName,
		user.LastName, user.Email, user.DateCreated, user.Status, user.Password)

	if saveErr != nil {
		logger.Error("error when trying to prepare save user ", saveErr)
		return errors.NewInternalServerError("error when trying to prepare save user", builtInError.New("database error"))

	}

	userId, err := insertResult.LastInsertId()

	if err != nil {
		logger.Error("error when trying to get saved user id", err)
		return errors.NewInternalServerError("error when trying to get saved user id", builtInError.New("database error"))

	}

	user.Id = userId

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("error when trying to prepare update user statement", builtInError.New("database error"))

	}
	defer stmt.Close()

	_, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)

	if saveErr != nil {
		logger.Error("error when trying to update user", err)
		return errors.NewInternalServerError("error when trying to update user", builtInError.New("database error"))

	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("error when trying to prepare delete user statement", builtInError.New("database error"))

	}
	defer stmt.Close()

	if _, deleteErr := stmt.Exec(user.Id); deleteErr != nil {
		logger.Error("error when trying to delete user", deleteErr)
		return errors.NewInternalServerError("error when trying to delete user", builtInError.New("database error"))

	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user by status statement", err)
		return nil, errors.NewInternalServerError("error when trying to prepare find user by status statement", builtInError.New("database error"))

	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to query find user", err)
		return nil, errors.NewInternalServerError("error when trying to query find user", builtInError.New("database error"))

	}
	defer rows.Close() // must close after make sure that you get data.

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, errors.NewInternalServerError("error when trying to find user", builtInError.New("database error"))

		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status : %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user statement by email and password", err)
		return errors.NewInternalServerError("error when trying to prepare get user statement by email and password", builtInError.New("database error"))

	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if err = result.Scan(&user.Id, &user.FirstName,
		&user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		logger.Error("error when trying to get user by id by email and password", err)
		if strings.Contains(err.Error(), mysql_utils.ErrorNoRows) {
			return errors.NewNotFoundError("user not found")
		}
		return errors.NewInternalServerError("error when trying to get user by id by email and password", builtInError.New("database error"))

	}

	return nil
}
