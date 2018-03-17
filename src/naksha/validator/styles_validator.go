package validator

import (
	"fmt"
	"net/http"
	"strconv"
)

type StylesValidator struct {
	BaseValidator
	geometry_type string
}

func (sv *StylesValidator) Validate() []string {
	switch sv.geometry_type {
	case "polygon":
		sv.validateColor("fill", "Fill Color")
		sv.validateOpacity("fill_opacity", "Fill Opacity")
		sv.validateColor("stroke", "Stroke Color")
		sv.validateOpacity("stroke_opacity", "Stroke Opacity")
		sv.validateStrokeWidth("stroke_width", "Stroke Width")
	case "linestring":
		sv.validateColor("stroke", "Stroke Color")
		sv.validateOpacity("stroke_opacity", "Stroke Opacity")
		sv.validateStrokeWidth("stroke_width", "Stroke Width")
	case "point":
		sv.validateColor("fill", "Fill Color")
		sv.validateOpacity("fill_opacity", "Fill Opacity")
		sv.validateColor("stroke", "Stroke Color")
		sv.validateOpacity("stroke_opacity", "Stroke Opacity")
		sv.validateStrokeWidth("stroke_width", "Stroke Width")
		sv.validateWidthHeight("width", "Width")
		sv.validateWidthHeight("height", "Height")
	default:
		sv.AppendError("Invalid geometry type")
	}

	return sv.GetErrors()
}

func (sv *StylesValidator) validateColor(field, label string) {
	value := sv.PostFormValue(field)
	length := len(value)
	if length == 0 {
		sv.AppendError(label + " is required")
	} else if length == 7 {
		c := fmt.Sprintf("%c", value[0])
		if c == "#" {
			for i := 1; i < 7; i++ {
				c1 := fmt.Sprintf("%c", value[i])
				switch c1 {
				case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f":
				default:
					sv.AppendError(label + " invalid value")
				}
			}
		} else {
			sv.AppendError(label + ": invalid format")
		}
	} else {
		sv.AppendError(label + ": should not be longer than 7 characters")
	}
}

func (sv *StylesValidator) validateOpacity(field, label string) {
	value := sv.PostFormValue(field)
	length := len(value)
	if length == 0 {
		sv.AppendError(label + " is required")
	} else if length <= 4 {
		i, err := strconv.ParseFloat(value, 32)
		if err == nil {
			if i < 0.0 || i > 1.0 {
				sv.AppendError(label + ": invalid value")
			}
		} else {
			sv.AppendError(label + ": invalid value")
		}
	} else {
		sv.AppendError(label + " should not be longer than 4 characters")
	}
}

func (sv *StylesValidator) validateStrokeWidth(field, label string) {
	value := sv.PostFormValue(field)
	length := len(value)
	if length == 0 {
		sv.AppendError(label + " is required")
	} else {
		i, err := strconv.ParseFloat(value, 32)
		if err == nil {
			if i < 0.0 {
				sv.AppendError(label + ": invalid value")
			}
		} else {
			sv.AppendError(label + ": invalid value")
		}
	}
}

func (sv *StylesValidator) validateWidthHeight(field, label string) {
	value := sv.PostFormValue(field)
	length := len(value)
	if length == 0 {
		sv.AppendError(label + ": is required")
	} else {
		i, err := strconv.Atoi(value)
		if err == nil {
			if i <= 0 {
				sv.AppendError(label + ": should be greater than 0")
			}
		} else {
			sv.AppendError(label + ": invalid value")
		}
	}
}

func MakeStylesValidator(request *http.Request, geometry_type string) StylesValidator {
	errors := make([]string, 0)
	return StylesValidator{BaseValidator{request, errors}, geometry_type}
}
