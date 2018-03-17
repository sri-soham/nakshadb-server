package validator

import (
	"net/http"
)

type UserLoginValidator struct {
	BaseValidator
}

func (uv *UserLoginValidator) Validate() []string {
	uv.Required("username", "Username")
	uv.Required("password", "Password")

	return uv.GetErrors()
}

func MakeUserLoginValidator(request *http.Request) UserLoginValidator {
	errors := make([]string, 0)
	return UserLoginValidator{BaseValidator{request, errors}}
}
