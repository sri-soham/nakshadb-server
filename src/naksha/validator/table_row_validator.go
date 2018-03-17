package validator

import (
	"net/http"
	"strings"
)

type TableRowValidator struct {
	BaseValidator
}

func (truv *TableRowValidator) ValidateUpdate() []string {
	truv.Required("column", "Column")
	if len(truv.GetErrors()) == 0 {
		column := truv.BaseValidator.PostFormValue("column")
		value := truv.BaseValidator.PostFormValue("value")
		if column == "the_geom" {
			if len(value) > 0 {
				err := truv.geomEwktString(value)
				if len(err) > 0 {
					truv.BaseValidator.AppendError(err)
				}
			}
		}
	}

	return truv.GetErrors()
}

func (truv *TableRowValidator) ValidateAdd() []string {
	with_geometry := truv.BaseValidator.PostFormValue("with_geometry")
	geometry := truv.BaseValidator.PostFormValue("geometry")
	if with_geometry == "1" {
		if len(geometry) == 0 {
			truv.BaseValidator.AppendError("Geometry is required")
		} else {
			err := truv.geomEwktString(geometry)
			if len(err) > 0 {
				truv.BaseValidator.AppendError(err)
			}
		}
	}

	return truv.GetErrors()
}

func (truv *TableRowValidator) geomEwktString(geom_str string) string {
	err := ""

	parts := strings.Split(geom_str, ";")
	if len(parts) == 2 {
		if parts[0] == "SRID=4326" {
			geom_type := strings.Split(parts[1], "(")[0]
			switch geom_type {
			case "POLYGON", "MULTIPOLYGON":
			case "POINT", "MULTIPOINT":
			case "LINESTRING", "MULTILINESTRING":
			default:
				err = "Invalid geometry type"
			}
		} else {
			err = "Invalid SRID value"
		}
	} else {
		err = "Invalid geometry ewkt string"
	}

	return err
}

func MakeTableRowValidator(request *http.Request) TableRowValidator {
	errors := make([]string, 0)
	return TableRowValidator{BaseValidator{request, errors}}
}
