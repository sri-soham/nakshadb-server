package validator

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func makeRequest(method string, values url.Values) *http.Request {
	request := httptest.NewRequest(method, "/abc", strings.NewReader(values.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request
}

func makeValidator(method string, values url.Values) BaseValidator {
	request := makeRequest(method, values)
	errs := make([]string, 0)
	validator := BaseValidator{request, errs}

	return validator
}

func TestRequired(t *testing.T) {
	values := url.Values{}
	validator := makeValidator(http.MethodPost, values)
	validator.Required("name", "Name")
	errs := validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("name", "something")
	validator = makeValidator(http.MethodPost, values)
	validator.Required("name", "Name")
	errs = validator.GetErrors()
	if len(errs) != 0 {
		t.Errorf("Error count: expected(0), actual(%d)", len(errs))
	}
}

func TestDigits(t *testing.T) {
	values := url.Values{}
	values.Set("id", "abc")
	validator := makeValidator(http.MethodPost, values)
	validator.Digits("id", "ID")
	errs := validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("id", "a34")
	validator = makeValidator(http.MethodPost, values)
	validator.Digits("id", "ID")
	errs = validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("id", "34b")
	validator = makeValidator(http.MethodPost, values)
	validator.Digits("id", "ID")
	errs = validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("id", "3b4")
	validator = makeValidator(http.MethodPost, values)
	validator.Digits("id", "ID")
	errs = validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("id", "43")
	validator = makeValidator(http.MethodPost, values)
	validator.Digits("id", "ID")
	errs = validator.GetErrors()
	if len(errs) != 0 {
		t.Errorf("Error count: expected(0), actual(%d)", len(errs))
	}
}

func TestMaxLength(t *testing.T) {
	values := url.Values{}
	values.Set("name", "Xavier")
	validator := makeValidator(http.MethodPost, values)
	validator.MaxLength("name", "Name", 4)
	errs := validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("name", "Xavier")
	validator = makeValidator(http.MethodPost, values)
	validator.MaxLength("name", "Name", 5)
	errs = validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("name", "Xavier")
	validator = makeValidator(http.MethodPost, values)
	validator.MaxLength("name", "Name", 6)
	errs = validator.GetErrors()
	if len(errs) != 0 {
		t.Errorf("Error count: expected(0), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("name", "Xavier")
	validator = makeValidator(http.MethodPost, values)
	validator.MaxLength("name", "Name", 7)
	errs = validator.GetErrors()
	if len(errs) != 0 {
		t.Errorf("Error count: expected(0), actual(%d)", len(errs))
	}
}

func TestMinLength(t *testing.T) {
	values := url.Values{}
	values.Set("name", "Xavier")
	validator := makeValidator(http.MethodPost, values)
	validator.MinLength("name", "Name", 8)
	errs := validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("name", "Xavier")
	validator = makeValidator(http.MethodPost, values)
	validator.MinLength("name", "Name", 7)
	errs = validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("name", "Xavier")
	validator = makeValidator(http.MethodPost, values)
	validator.MinLength("name", "Name", 6)
	errs = validator.GetErrors()
	if len(errs) != 0 {
		t.Errorf("Error count: expected(0), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("name", "Xavier")
	validator = makeValidator(http.MethodPost, values)
	validator.MinLength("name", "Name", 5)
	errs = validator.GetErrors()
	if len(errs) != 0 {
		t.Errorf("Error count: expected(0), actual(%d)", len(errs))
	}
}

func TestLength(t *testing.T) {
	values := url.Values{}
	values.Set("name", "Xavier")
	validator := makeValidator(http.MethodPost, values)
	validator.Length("name", "Name", 7)
	errs := validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("name", "Xavier")
	validator = makeValidator(http.MethodPost, values)
	validator.Length("name", "Name", 5)
	errs = validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("name", "Xavier")
	validator = makeValidator(http.MethodPost, values)
	validator.Length("name", "Name", 6)
	errs = validator.GetErrors()
	if len(errs) != 0 {
		t.Errorf("Error count: expected(0), actual(%d)", len(errs))
	}
}

func TestEquals(t *testing.T) {
	values := url.Values{}
	values.Set("email1", "name1@www.com")
	values.Set("email2", "name2@www.com")
	validator := makeValidator(http.MethodPost, values)
	validator.Equals("email1", "Email 1", "email2", "Email 2")
	errs := validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("email1", "name1@www.com")
	values.Set("email2", "name1@www.com")
	validator = makeValidator(http.MethodPost, values)
	validator.Equals("email1", "Email 1", "email2", "Email 2")
	errs = validator.GetErrors()
	if len(errs) != 0 {
		t.Errorf("Error count: expected(0), actual(%d)", len(errs))
	}
}

func TestInArray(t *testing.T) {
	values := url.Values{}
	values.Set("category", "10")
	validator := makeValidator(http.MethodPost, values)
	validator.InArray("category", "Category", []string{"1", "0", "20", "1", "01", "102", "210"})
	errs := validator.GetErrors()
	if len(errs) != 1 {
		t.Errorf("Error count: expected(1), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("category", "10")
	validator = makeValidator(http.MethodPost, values)
	validator.InArray("category", "Category", []string{"1", "0", "20", "1", "01", "102", "210", "10"})
	errs = validator.GetErrors()
	if len(errs) != 0 {
		t.Errorf("Error count: expected(0), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("category", "10")
	validator = makeValidator(http.MethodPost, values)
	validator.InArray("category", "Category", []string{"10", "1", "0", "20", "1", "01", "102", "210"})
	errs = validator.GetErrors()
	if len(errs) != 0 {
		t.Errorf("Error count: expected(0), actual(%d)", len(errs))
	}

	values = url.Values{}
	values.Set("category", "10")
	validator = makeValidator(http.MethodPost, values)
	validator.InArray("category", "Category", []string{"1", "0", "20", "10", "01", "102", "210"})
	errs = validator.GetErrors()
	if len(errs) != 0 {
		t.Errorf("Error count: expected(0), actual(%d)", len(errs))
	}
}
