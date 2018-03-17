package validator

import (
	"net/http"
	"regexp"
	"strconv"
)

type AddColumnValidator struct {
	BaseValidator
}

func (acv *AddColumnValidator) Validate() []string {
	acv.RequiredMaxLength("name", "Name", 60)
	acv.Required("data_type", "Data Type")
	acv.Digits("data_type", "Data Type")

	errors := acv.GetErrors()
	if len(errors) == 0 {
		data_type, _ := strconv.Atoi(acv.PostFormValue("data_type"))
		if data_type < 1 || data_type > 4 {
			acv.AppendError("Invalid value for data type")
		}
		name := acv.PostFormValue("name")
		ok, _ := regexp.MatchString("^[a-z].[a-z0-9_]{1,62}$", name)
		if !ok {
			msg := "Invalid name. Only lower case alphabets, digits, and " +
				"underscores are allowed. Name must begin with lower case alphabet"
			acv.AppendError(msg)
		}
	}

	return acv.GetErrors()
}

func MakeAddColumnValidator(request *http.Request) AddColumnValidator {
	errors := make([]string, 0)
	return AddColumnValidator{BaseValidator{request, errors}}
}
