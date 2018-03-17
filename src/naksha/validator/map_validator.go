package validator

import (
	"naksha/helper"
	"net/http"
)

type MapValidator struct {
	BaseValidator
}

func (mv *MapValidator) Validate() []string {
	mv.RequiredMaxLength("name", "Name", 64)
	mv.Required("layer", "Layer")
	mv.Digits("layer", "Layer")

	return mv.GetErrors()
}

func (mv *MapValidator) ValidateEdit() []string {
	mv.RequiredMaxLength("name", "Name", 64)

	return mv.GetErrors()
}

func (mv *MapValidator) ValidateBaseLayer(base_layers map[string]string, user map[string]string) []string {
	errs := make([]string, 0)
	base_layer := mv.PostFormValue("base_layer")
	if len(base_layer) == 0 {
		errs = append(errs, "Base layer is required")
	} else {
		_, ok := base_layers[base_layer]
		if !ok {
			errs = append(errs, "Invalid value for base layer")
		} else {
			if helper.IsGoogleMapsBaseLayer(base_layer) && len(user["google_maps_key"]) == 0 {
				errs = append(errs, "Please update google maps key in profile page")
			}
			if helper.IsBingMapsBaseLayer(base_layer) && len(user["bing_maps_key"]) == 0 {
				errs = append(errs, "Please update bing maps key in profile page")
			}
		}
	}

	return errs
}

func MakeMapValidator(request *http.Request) MapValidator {
	errs := make([]string, 0)
	return MapValidator{BaseValidator{request, errs}}
}
