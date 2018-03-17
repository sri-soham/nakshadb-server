package validator

import (
	"net/http"
	"net/url"
	"testing"
)

func TestStyleValidateColor(t *testing.T) {
	testStyleValidateColor(t, "", 1)
	testStyleValidateColor(t, "123456", 1)
	testStyleValidateColor(t, "12345678", 1)
	testStyleValidateColor(t, "123456#", 1)
	testStyleValidateColor(t, "123#456", 1)
	testStyleValidateColor(t, "#A23456", 1)
	testStyleValidateColor(t, "#23A456", 1)
	testStyleValidateColor(t, "#23456A", 1)
	testStyleValidateColor(t, "#g23456", 1)
	testStyleValidateColor(t, "#23g456", 1)
	testStyleValidateColor(t, "#23456g", 1)
	testStyleValidateColor(t, "#ffff00", 0)
	testStyleValidateColor(t, "#abcdef", 0)
	testStyleValidateColor(t, "#123456", 0)
	testStyleValidateColor(t, "#7890ab", 0)
}

func TestStyleValidateOpacity(t *testing.T) {
	testStyleValidateOpacity(t, "", 1)
	testStyleValidateOpacity(t, "ab", 1)
	testStyleValidateOpacity(t, "EB", 1)
	testStyleValidateOpacity(t, "10.20", 1)
	testStyleValidateOpacity(t, "-0.01", 1)
	testStyleValidateOpacity(t, "1.01", 1)
	testStyleValidateOpacity(t, "0.01", 0)
	testStyleValidateOpacity(t, "1.00", 0)
	testStyleValidateOpacity(t, "0.50", 0)
}

func TestStyleValidateStrokeWidth(t *testing.T) {
	testStyleValidateStrokeWidth(t, "", 1)
	testStyleValidateStrokeWidth(t, "xya", 1)
	testStyleValidateStrokeWidth(t, "Xya", 1)
	testStyleValidateStrokeWidth(t, "ABD", 1)
	testStyleValidateStrokeWidth(t, "-0.50", 1)
	testStyleValidateStrokeWidth(t, "0.50", 0)
	testStyleValidateStrokeWidth(t, "10", 0)
	testStyleValidateStrokeWidth(t, "20.75", 0)
}

func TestStyleValidateWidthHeight(t *testing.T) {
	testStyleValidateWidthHeight(t, "", 1)
	testStyleValidateWidthHeight(t, "gh", 1)
	testStyleValidateWidthHeight(t, "UD", 1)
	testStyleValidateWidthHeight(t, "UalD", 1)
	testStyleValidateWidthHeight(t, "-1", 1)
	testStyleValidateWidthHeight(t, "5", 0)
	testStyleValidateWidthHeight(t, "20", 0)
}

func TestStyleValidate(t *testing.T) {
	values := url.Values{}
	request := makeRequest(http.MethodPost, values)
	validator := MakeStylesValidator(request, "something")
	errs := validator.Validate()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), found(%d)", len(errs))
	}
}

func TestStyleValidatePoint(t *testing.T) {
	testStyleValidate(t, "point", getStyleValidatePointValues(), "fill", 1)
	testStyleValidate(t, "point", getStyleValidatePointValues(), "fill_opacity", 1)
	testStyleValidate(t, "point", getStyleValidatePointValues(), "stroke", 1)
	testStyleValidate(t, "point", getStyleValidatePointValues(), "stroke_opacity", 1)
	testStyleValidate(t, "point", getStyleValidatePointValues(), "stroke_width", 1)
	testStyleValidate(t, "point", getStyleValidatePointValues(), "width", 1)
	testStyleValidate(t, "point", getStyleValidatePointValues(), "height", 1)
	testStyleValidate(t, "point", getStyleValidatePointValues(), "nothing", 0)
}

func TestStyleValidateLinestring(t *testing.T) {
	testStyleValidate(t, "linestring", getStyleValidateLinestringValues(), "stroke", 1)
	testStyleValidate(t, "linestring", getStyleValidateLinestringValues(), "stroke_opacity", 1)
	testStyleValidate(t, "linestring", getStyleValidateLinestringValues(), "stroke_width", 1)
	testStyleValidate(t, "linestring", getStyleValidateLinestringValues(), "nothing", 0)
}

func TestStyleValidatePolygon(t *testing.T) {
	testStyleValidate(t, "polygon", getStyleValidatePolygonValues(), "fill", 1)
	testStyleValidate(t, "polygon", getStyleValidatePolygonValues(), "fill_opacity", 1)
	testStyleValidate(t, "polygon", getStyleValidatePolygonValues(), "stroke", 1)
	testStyleValidate(t, "polygon", getStyleValidatePolygonValues(), "stroke_opacity", 1)
	testStyleValidate(t, "polygon", getStyleValidatePolygonValues(), "stroke_width", 1)
	testStyleValidate(t, "polygon", getStyleValidatePolygonValues(), "nothing", 0)
}

func testStyleValidateColor(t *testing.T, color string, expected_err_len int) {
	values := url.Values{}
	values.Set("fill", color)
	request := makeRequest(http.MethodPost, values)
	validator := MakeStylesValidator(request, "point")
	validator.validateColor("fill", "Fill Color")
	errs := validator.GetErrors()
	err_len := len(errs)
	if err_len != expected_err_len {
		t.Errorf("Error count: expected(%d), found(%d)", expected_err_len, err_len)
	}
}

func getStyleValidatePointValues() url.Values {
	values := url.Values{}
	values.Set("fill", "#ffffff")
	values.Set("fill_opacity", "0.5")
	values.Set("stroke", "#dddddd")
	values.Set("stroke_opacity", "0.75")
	values.Set("stroke_width", "1.0")
	values.Set("width", "10")
	values.Set("height", "10")

	return values
}

func getStyleValidateLinestringValues() url.Values {
	values := url.Values{}
	values.Set("stroke", "#dddddd")
	values.Set("stroke_opacity", "0.75")
	values.Set("stroke_width", "1.0")

	return values
}

func getStyleValidatePolygonValues() url.Values {
	values := url.Values{}
	values.Set("fill", "#ffffff")
	values.Set("fill_opacity", "0.5")
	values.Set("stroke", "#dddddd")
	values.Set("stroke_opacity", "0.75")
	values.Set("stroke_width", "1.0")

	return values
}

func testStyleValidateOpacity(t *testing.T, opacity string, expected_err_len int) {
	values := url.Values{}
	values.Set("fill_opacity", opacity)
	request := makeRequest(http.MethodPost, values)
	validator := MakeStylesValidator(request, "point")
	validator.validateOpacity("fill_opacity", "Fill Opacity")
	errs := validator.GetErrors()
	err_len := len(errs)
	if err_len != expected_err_len {
		t.Errorf("Error count: expected(%d), found(%d)", expected_err_len, err_len)
	}
}

func testStyleValidateStrokeWidth(t *testing.T, swidth string, expected_err_len int) {
	values := url.Values{}
	values.Set("stroke_width", swidth)
	request := makeRequest(http.MethodPost, values)
	validator := MakeStylesValidator(request, "point")
	validator.validateStrokeWidth("stroke_width", "Stroke Width")
	errs := validator.GetErrors()
	err_len := len(errs)
	if err_len != expected_err_len {
		t.Errorf("Error count: expected(%d), found(%d)", expected_err_len, err_len)
	}
}

func testStyleValidateWidthHeight(t *testing.T, width string, expected_err_len int) {
	values := url.Values{}
	values.Set("width", width)
	request := makeRequest(http.MethodPost, values)
	validator := MakeStylesValidator(request, "point")
	validator.validateStrokeWidth("width", "Width")
	errs := validator.GetErrors()
	err_len := len(errs)
	if err_len != expected_err_len {
		t.Errorf("Error count: expected(%d), found(%d). Width: %s", expected_err_len, err_len, width)
	}
}

func testStyleValidate(t *testing.T, geometry_type string, values url.Values, to_remove_col string, expected_err_len int) {
	values.Set(to_remove_col, "")
	request := makeRequest(http.MethodPost, values)
	validator := MakeStylesValidator(request, geometry_type)
	errs := validator.Validate()
	err_len := len(errs)
	if err_len != expected_err_len {
		t.Errorf("Error count: expected(%d), found(%d). Column: %v", expected_err_len, err_len, to_remove_col)
	}
}
