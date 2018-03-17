package controller

import (
	"net/http"
	"regexp"
	"strings"
)

type Route struct {
	method     string
	uri        string
	controller IController
	handler    Handler
	validators map[string]string
}

func (r *Route) IsMatch(request *http.Request) (map[string]string, bool) {
	var matches bool
	var count1 int
	var count2 int
	var req_part string
	var uri_part string
	var i int

	uri_params := make(map[string]string)

	if request.Method != r.method {
		return uri_params, false
	}

	req_parts := strings.Split(request.URL.Path, "/")
	uri_parts := strings.Split(r.uri, "/")
	// for a string like "/users/show", uri_parts[0] will be "" empty string.
	req_parts = req_parts[1:]
	uri_parts = uri_parts[1:]
	if len(req_parts) == len(uri_parts) {
		count1 = len(uri_parts)
		count2 = 0
		for i, req_part = range req_parts {
			uri_part = uri_parts[i]
			if uri_part == req_part {
				count2++
			} else {
				if strings.HasPrefix(uri_part, "{") && strings.HasSuffix(uri_part, "}") {
					frag := strings.TrimLeft(uri_part, "{")
					frag = strings.TrimRight(frag, "}")
					validator, ok := r.validators[frag]
					if ok {
						match, _ := regexp.MatchString(validator, req_part)
						if match {
							uri_params[frag] = req_part
							count2++
						}
					} else {
						uri_params[frag] = req_part
						count2++
					}
				}
			}
		}
		matches = (count1 == count2)
	} else {
		matches = false
	}

	return uri_params, matches
}

func (r *Route) GetController() IController {
	return r.controller
}

func (r *Route) GetHandler() Handler {
	return r.handler
}
