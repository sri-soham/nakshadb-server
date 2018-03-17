package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"naksha"
	"naksha/db"
	"naksha/helper"
	"naksha/logger"
	"naksha/validator"
	"strconv"
	"time"
)

type UserService interface {
	Login(req *naksha.Request) helper.Result
	ChangePassword(req *naksha.Request) helper.Result
	ProfilePost(req *naksha.Request) helper.Result
	UserDetails(req *naksha.Request) helper.Result
}

type UserServiceImpl struct {
	map_user_dao db.Dao
}

func (us *UserServiceImpl) Login(req *naksha.Request) helper.Result {
	var id int
	var name string
	var schema_name string
	var result helper.Result

	validator := validator.MakeUserLoginValidator(req.Request)
	errors := validator.Validate()
	if len(errors) == 0 {
		where := make(map[string]interface{})
		where["username"] = req.PostFormValue("username")
		user_details, err := us.map_user_dao.FindWhere(db.MasterUser(), where)
		if err != nil {
			return helper.ErrorResultError(err)
		} else {
			login_attempts, err := strconv.Atoi(user_details["login_attempts"])
			if err != nil {
				login_attempts = 0
			}
			if login_attempts >= 5 {
				start_time, _ := time.Parse("2006-01-02 15:04:05", user_details["last_login_time"])
				diff := time.Now().UTC().Sub(start_time)
				if diff.Minutes() > 15 {
					values := make(map[string]interface{})
					values["login_attempts"] = 0
					_, err = us.map_user_dao.Update(db.MasterUser(), values, where)
				} else {
					errors = append(errors, "Account locked. Please try after some time")
				}
			}

			if len(errors) == 0 {
				hashed := []byte(user_details["password"])
				password := []byte(req.PostFormValue("password"))
				err := bcrypt.CompareHashAndPassword(hashed, password)
				if err != nil {
					errors = append(errors, "Password mismatch")
					login_attempts += 1
					values := make(map[string]interface{})
					values["login_attempts"] = login_attempts
					values["last_login_time"] = time.Now().UTC().Format(time.RFC3339)
					_, err = us.map_user_dao.Update(db.MasterUser(), values, where)
					if err != nil {
						logger.LogMessage(fmt.Sprintf("Could not update login-attempts for %s. login-attempts = %d", where["username"], login_attempts))
					}
				} else {
					id, _ = strconv.Atoi(user_details["id"])
					name = user_details["name"]
					schema_name = user_details["schema_name"]
				}
			}

			if len(errors) == 0 {
				values := make(map[string]interface{})
				values["login_attempts"] = 0
				values["last_login_time"] = time.Now().UTC().Format(time.RFC3339)
				values["last_login_ip"] = req.Request.RemoteAddr
				_, err = us.map_user_dao.Update(db.MasterUser(), values, where)
				if err != nil {
					logger.LogMessage(fmt.Sprintf("Could not update login time for user %s", user_details["username"]))
				}
				result = helper.MakeSuccessResult()
				result.AddToData("id", id)
				result.AddToData("name", name)
				result.AddToData("schema_name", schema_name)
			} else {
				result = helper.MakeErrorResult()
				result.SetErrors(errors)
			}
		}
	} else {
		result = helper.MakeErrorResult()
		result.SetErrors(errors)
	}

	return result
}

func (us *UserServiceImpl) ChangePassword(req *naksha.Request) helper.Result {
	var result helper.Result

	validator := validator.MakeChangePasswordValidator(req.Request)
	errors := validator.Validate()
	if len(errors) == 0 {
		user_id := helper.UserIDFromSession(req.Session)
		user_details, err := us.map_user_dao.Find(db.MasterUser(), user_id)
		if err != nil {
			return helper.ErrorResultError(err)
		}

		hashed := []byte(user_details["password"])
		password := []byte(req.PostFormValue("current_password"))
		err = bcrypt.CompareHashAndPassword(hashed, password)
		if err != nil {
			return helper.ErrorResultString("Password mismatch")
		}

		new_password := []byte(req.PostFormValue("new_password"))
		new_hash, err := bcrypt.GenerateFromPassword(new_password, bcrypt.DefaultCost)
		if err != nil {
			return helper.ErrorResultString("Could not generate hash")
		}

		values := make(map[string]interface{})
		values["password"] = new_hash
		where := make(map[string]interface{})
		where["id"] = user_id
		_, err = us.map_user_dao.Update(db.MasterUser(), values, where)
		if err != nil {
			result = helper.ErrorResultString("Could not update password")
		} else {
			result = helper.MakeSuccessResult()
			result.AddToData("message", "Password Updated")
		}

	} else {
		result = helper.MakeErrorResult()
		result.SetErrors(errors)
	}

	return result
}

func (us *UserServiceImpl) ProfilePost(req *naksha.Request) helper.Result {
	var result helper.Result

	validator := validator.MakeUserProfileValidator(req.Request)
	errs := validator.Validate()
	if len(errs) > 0 {
		result = helper.MakeErrorResult()
		result.SetErrors(errs)
		return result
	}
	key := req.PostFormValue("key")
	value := req.PostFormValue("value")
	user_id := helper.UserIDFromSession(req.Session)
	values := map[string]interface{}{key: value}
	where := map[string]interface{}{"id": user_id}
	_, err := us.map_user_dao.Update(db.MasterUser(), values, where)
	if err != nil {
		result = helper.ErrorResultString("Query Failed")
	} else {
		result = helper.MakeSuccessResult()
	}

	return result
}

func (us *UserServiceImpl) UserDetails(req *naksha.Request) helper.Result {
	var result helper.Result

	user_id := helper.UserIDFromSession(req.Session)
	user_details, err := us.map_user_dao.Find(db.MasterUser(), user_id)
	if err != nil {
		return helper.ErrorResultError(err)
	}

	result = helper.MakeSuccessResult()
	result.AddToData("user_details", user_details)
	return result
}

func MakeUserService(map_user_dao db.Dao) UserService {
	return &UserServiceImpl{map_user_dao}
}
