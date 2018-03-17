package validator

import "net/http"

type UserProfileValidator struct {
	BaseValidator
}

func (upv *UserProfileValidator) Validate() []string {
	upv.Required("key", "Key")
	upv.MaxLength("value", "Value", 128)
	upv.InArray("key", "Key", []string{"google_maps_key", "bing_maps_key"})

	return upv.GetErrors()
}

func MakeUserProfileValidator(request *http.Request) UserProfileValidator {
	errors := make([]string, 0)
	return UserProfileValidator{BaseValidator{request, errors}}
}
