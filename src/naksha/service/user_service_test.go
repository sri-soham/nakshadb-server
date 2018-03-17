package service_test

import (
	"errors"
	"github.com/gorilla/sessions"
	"naksha"
	"naksha/db"
	"naksha/service"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var store = sessions.NewCookieStore([]byte("ascii-unicode-uint-nibble"))

func makeErrorSelectResult(err_msg string) db.SelectResult {
	return db.SelectResult{errors.New(err_msg), make([]string, 0), make([]map[string]string, 0)}
}

func makeRequest(method string, uri_params map[string]string, values url.Values, authenticated bool) *naksha.Request {
	request := httptest.NewRequest(method, "/abc", strings.NewReader(values.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	session, _ := store.Get(request, "sess")
	if authenticated {
		session.Values["user_id"] = 1
		session.Values["schema_name"] = "abcdef"
	}
	naksha_request := &naksha.Request{request, session, uri_params}

	return naksha_request
}

func TestUserLoginInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	service := service.MakeUserService(map_dao)

	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, false)

	result := service.Login(request)
	if result.IsSuccess() {
		t.Errorf("UserLoginInvalidForm. validation failed")
	}
}

func TestUserLoginInvalidPassword(t *testing.T) {
	map_dao := &UserServiceLoginInvalidPasswordImpl{&DaoImpl{}}
	service := service.MakeUserService(map_dao)

	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("username", "test")
	values.Set("password", "something")
	request := makeRequest(http.MethodPost, uri_params, values, false)
	result := service.Login(request)
	if result.IsSuccess() {
		t.Errorf("UserLoginInvalidPassword. Wrong password is getting authenticated")
	}
}

func TestUserLoginInvalidUsername(t *testing.T) {
	map_dao := &UserServiceLoginInvalidUsernameImpl{&DaoImpl{}}
	service := service.MakeUserService(map_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("username", "someuser")
	values.Set("password", "test1234")
	request := makeRequest(http.MethodPost, uri_params, values, false)
	result := service.Login(request)
	if result.IsSuccess() {
		t.Errorf("UserLoginInvalidUsername. Non existent username is getting authenticated")
	}
}

func TestUserLoginValid(t *testing.T) {
	map_dao := &UserServiceLoginValidImpl{&DaoImpl{}}
	service := service.MakeUserService(map_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("username", "test")
	values.Set("password", "test1234")
	request := makeRequest(http.MethodPost, uri_params, values, false)
	result := service.Login(request)
	if !result.IsSuccess() {
		t.Errorf("LoginValid. Valid user/pass not authenticated")
	}
	user_id := result.GetIntData("id")
	if user_id != 1 {
		t.Errorf("UserLoginValid: user id is invalid or not set. Expected(1), found(%d)", user_id)
	}
	name := result.GetStringData("name")
	if name != "Tester" {
		t.Errorf("UserLoginValid: name is invalid or not set. Expected(Tester), found(%v)", name)
	}
}

func TestUserChangePasswordInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	service := service.MakeUserService(map_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.ChangePassword(request)
	if result.IsSuccess() {
		t.Errorf("UserChangePasswordInvalidForm: validation failed")
	}
}

func TestUserChangePasswordWrongCurrentPassword(t *testing.T) {
	map_dao := &UserServiceChangePasswordWrongCurrentPasswordImpl{&DaoImpl{}}
	service := service.MakeUserService(map_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("current_password", "testabcd")
	values.Set("new_password", "abcdefgh")
	values.Set("confirm_password", "abcdefgh")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.ChangePassword(request)
	if result.IsSuccess() {
		t.Errorf("UserChangePasswordWrongCurrentPassword: not checking for current password")
	}
}

func TestUserChangePasswordValid(t *testing.T) {
	map_dao := &UserServiceChangePasswordValidImpl{&DaoImpl{}}
	service := service.MakeUserService(map_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("current_password", "test1234")
	values.Set("new_password", "abcdefgh")
	values.Set("confirm_password", "abcdefgh")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.ChangePassword(request)
	if !result.IsSuccess() {
		t.Errorf("UserChangePasswordValid: valid form and current password not working")
	}
}

func TestUserProfileInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	service := service.MakeUserService(map_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.ProfilePost(request)
	if result.IsSuccess() {
		t.Errorf("UserProfileInvalidForm: invalid form post is succeeding")
	}
}

func TestUserProfileValid(t *testing.T) {
	map_dao := &DaoImpl{}
	service := service.MakeUserService(map_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	bing_maps_key := "sdflk1234sdfljsdf"
	values.Set("key", "bing_maps_key")
	values.Set("value", bing_maps_key)
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.ProfilePost(request)
	if !result.IsSuccess() {
		t.Errorf("UserProfileValid: invalid form post is succeeding")
	}
}

func TestUserDetails(t *testing.T) {
	map_dao := &UserServiceChangePasswordValidImpl{&DaoImpl{}}
	service := service.MakeUserService(map_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodGet, uri_params, values, true)
	result := service.UserDetails(request)
	if result.IsSuccess() {
		v, ok := result.GetDataByKey("user_details")
		if ok {
			data := v.(map[string]string)
			if data["name"] != "Tester" {
				t.Errorf("UserDetails: other users details returned")
			}
		} else {
			t.Errorf("UserDetails: user_details key not set in the result")
		}
	} else {
		t.Errorf("UserDetails: Could not fetch details")
	}
}
