package validator

import (
	"net/http"
	"net/url"
	"testing"
)

func getTableRowUpdateValues() url.Values {
	values := url.Values{}
	values.Set("column", "the_geom")
	values.Set("value", "SRID=4326;POINT(78.65957737 17.45620948)")

	return values
}

func TestTableRowUpdateValidate(t *testing.T) {
	testTableRowUpdate(t, getTableRowUpdateValues(), "column", "", 1)
	testTableRowUpdate(t, getTableRowUpdateValues(), "value", "", 0)
	testTableRowUpdate(t, getTableRowUpdateValues(), "column", "something", 0)
	testTableRowUpdate(t, getTableRowUpdateValues(), "value", "POINT(78.65957737 17.45620948)", 1)
	testTableRowUpdate(t, getTableRowUpdateValues(), "value", "SRID=3857;POINT(78.65957737 17.45620948)", 1)
	testTableRowUpdate(t, getTableRowUpdateValues(), "value", "SRID=4326;POINT(78.65957737 17.45620948)", 0)
	testTableRowUpdate(t, getTableRowUpdateValues(), "value", "SRID=4326;MULTIPOINT(78.65957737 17.45620948)", 0)
	testTableRowUpdate(t, getTableRowUpdateValues(), "value", "SRID=4326;MULTILINESTRING(78.65957737 17.45620948)", 0)
	testTableRowUpdate(t, getTableRowUpdateValues(), "value", "SRID=4326;MULTIPOLYGON(78.65957737 17.45620948)", 0)
	testTableRowUpdate(t, getTableRowUpdateValues(), "value", "SRID=4326;LINESTRING(78.65957737 17.45620948)", 0)
	testTableRowUpdate(t, getTableRowUpdateValues(), "value", "SRID=4326;POLYGON(78.65957737 17.45620948)", 0)
}

func TestTableRowAddValidate(t *testing.T) {
	values := url.Values{}
	values.Set("with_geometry", "0")
	testTableRowAdd(t, values, false, 0)

	values = url.Values{}
	values.Set("with_geometry", "1")
	testTableRowAdd(t, values, true, 1)

	values = url.Values{}
	values.Set("with_geometry", "1")
	values.Set("geometry", "SRID=4326;POINT(78.65957737 17.45620948)")
	testTableRowAdd(t, values, true, 0)
}

func testTableRowUpdate(t *testing.T, values url.Values, field string, value string, expected_err_len int) {
	values.Set(field, value)
	request := makeRequest(http.MethodPost, values)
	validator := MakeTableRowValidator(request)
	errs := validator.ValidateUpdate()
	err_len := len(errs)
	if err_len != expected_err_len {
		t.Errorf("Error count: expected(%d), found(%d). Field: %s, value: %s", expected_err_len, err_len, field, value)
	}
}

func testTableRowAdd(t *testing.T, values url.Values, with_geometry bool, expected_err_len int) {
	request := makeRequest(http.MethodPost, values)
	validator := MakeTableRowValidator(request)
	errs := validator.ValidateAdd()
	err_len := len(errs)
	if err_len != expected_err_len {
		t.Errorf("Error count: with-geometry = (%v); count: expected(%d), found(%d). errs = %v", with_geometry, expected_err_len, err_len, errs)
	}
}
