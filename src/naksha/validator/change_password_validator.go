package validator

import (
	"net/http"
)

type ChangePasswordValidator struct {
	BaseValidator
}

func (cpv *ChangePasswordValidator) Validate() []string {
	cpv.Required("current_password", "Current Password")
	cpv.Required("new_password", "New Password")
	cpv.MinLength("new_password", "New Password", 8)
	cpv.Required("confirm_password", "Confirm New Password")
	cpv.MinLength("confirm_password", "Confirm Password", 8)
	cpv.Equals("new_password", "New Password", "confirm_password", "Confirm New Passwod")

	return cpv.GetErrors()
}

func MakeChangePasswordValidator(request *http.Request) ChangePasswordValidator {
	errors := make([]string, 0)
	return ChangePasswordValidator{BaseValidator{request, errors}}
}
