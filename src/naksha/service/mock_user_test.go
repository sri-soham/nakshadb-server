package service_test

import (
	"errors"
)

type UserServiceLoginInvalidPasswordImpl struct {
	*DaoImpl
}

func (us *UserServiceLoginInvalidPasswordImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := make(map[string]string)
	row["id"] = "1"
	row["name"] = "Tester"
	row["username"] = "test"
	row["password"] = "SomethingLongNotKnown"

	return row, nil
}

type UserServiceLoginInvalidUsernameImpl struct {
	*DaoImpl
}

func (us *UserServiceLoginInvalidUsernameImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := make(map[string]string)

	return row, errors.New("No such record")
}

type UserServiceLoginValidImpl struct {
	*DaoImpl
}

func (us *UserServiceLoginValidImpl) FindWhere(table string, where map[string]interface{}) (map[string]string, error) {
	row := make(map[string]string)
	row["id"] = "1"
	row["name"] = "Tester"
	row["username"] = "test"
	row["password"] = "$2a$10$470XQGXP/mStul31B6pAAOkeDYDTiX8nxY7AFQHyB6u6HtyxJB.ky"

	return row, nil
}

type UserServiceChangePasswordWrongCurrentPasswordImpl struct {
	*DaoImpl
}

func (us *UserServiceChangePasswordWrongCurrentPasswordImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := make(map[string]string)
	row["id"] = "1"
	row["name"] = "Tester"
	row["username"] = "test"
	row["password"] = "$2a$10$470XQGXP/mStul31B6pAAOkeDYDTiX8nxY7AFQHyB6u6HtyxJB.ky"

	return row, nil
}

type UserServiceChangePasswordValidImpl struct {
	*DaoImpl
}

func (us *UserServiceChangePasswordValidImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := make(map[string]string)
	row["id"] = "1"
	row["name"] = "Tester"
	row["username"] = "test"
	row["password"] = "$2a$10$470XQGXP/mStul31B6pAAOkeDYDTiX8nxY7AFQHyB6u6HtyxJB.ky"

	return row, nil
}
