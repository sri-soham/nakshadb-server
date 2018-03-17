package validator

import (
	"fmt"
	"net/http"
	"regexp"
)

type BaseValidator struct {
	request *http.Request
	errors  []string
}

func (bv *BaseValidator) PostFormValue(field string) string {
	return bv.request.PostFormValue(field)
}

func (bv *BaseValidator) AppendError(err string) {
	bv.errors = append(bv.errors, err)
}

func (bv *BaseValidator) Required(field string, label string) {
	val := bv.request.PostFormValue(field)
	if len(val) == 0 {
		bv.errors = append(bv.errors, label+" is required")
	}
}

func (bv *BaseValidator) RequiredMaxLength(field string, label string, max_len int) {
	bv.Required(field, label)
	bv.MaxLength(field, label, max_len)
}

func (bv *BaseValidator) Digits(field string, label string) {
	val := bv.request.PostFormValue(field)
	pattern := "^[\\d]+$"
	if len(val) > 0 {
		matched, _ := regexp.MatchString(pattern, val)
		if !matched {
			bv.errors = append(bv.errors, label+": only digits are allowed")
		}
	}
}

func (bv *BaseValidator) MaxLength(field, label string, max_len int) {
	val := bv.request.PostFormValue(field)
	if len(val) > max_len {
		bv.errors = append(bv.errors, fmt.Sprintf("%v should not have more than %v characters", label, max_len))
	}
}

func (bv *BaseValidator) MinLength(field, label string, min_len int) {
	val := bv.request.PostFormValue(field)
	if len(val) < min_len {
		bv.errors = append(bv.errors, fmt.Sprintf("%v should have at least %v characters", label, min_len))
	}
}

func (bv *BaseValidator) Length(field, label string, length int) {
	val := bv.request.PostFormValue(field)
	if len(val) != length {
		bv.errors = append(bv.errors, fmt.Sprintf("%v: max. allowed length is %v", label, length))
	}
}

func (bv *BaseValidator) Equals(field1, label1, field2, label2 string) {
	val1 := bv.request.PostFormValue(field1)
	val2 := bv.request.PostFormValue(field2)
	if val1 != val2 {
		bv.errors = append(bv.errors, fmt.Sprintf("%v must be same as %v", label1, label2))
	}
}

func (bv *BaseValidator) InArray(field string, label string, allowed_values []string) {
	value := bv.request.PostFormValue(field)
	if len(value) > 0 {
		in_array := false
		for _, v := range allowed_values {
			if v == value {
				in_array = true
				break
			}
		}
		if !in_array {
			bv.errors = append(bv.errors, fmt.Sprintf("%v does not have allowed value", label))
		}
	}
}

func (bv *BaseValidator) GetErrors() []string {
	return bv.errors
}
