package validator

import (
	"net/http"
	"net/url"
	"testing"
)

func getAddColumnValidatorValues() url.Values {
	values := url.Values{}
	values.Set("name", "census_wards")
	values.Set("data_type", "1")

	return values
}

func TestAddColumnDataType(t *testing.T) {
	testAddColumnDataType(t, "data_type", "0", 1)
	testAddColumnDataType(t, "data_type", "5", 1)
	testAddColumnDataType(t, "name", "1map_type", 1)
	testAddColumnDataType(t, "name", "_map_type", 1)
	testAddColumnDataType(t, "name", "Maptype", 1)
	testAddColumnDataType(t, "name", "mapType", 1)
	testAddColumnDataType(t, "name", "maptypE", 1)
	testAddColumnDataType(t, "data_type", "1", 0)
	testAddColumnDataType(t, "data_type", "4", 0)
	testAddColumnDataType(t, "name", "c_123_wards", 0)
	testAddColumnDataType(t, "name", "census_wards_1", 0)
}

func testAddColumnDataType(t *testing.T, col string, val string, expected_len int) {
	values := getAddColumnValidatorValues()
	values.Set(col, val)
	request := makeRequest(http.MethodPost, values)
	validator := MakeAddColumnValidator(request)
	errs := validator.Validate()
	err_len := len(errs)
	if err_len != expected_len {
		t.Errorf("Error count: expected(1), actual(%d)", err_len)
	}
}
