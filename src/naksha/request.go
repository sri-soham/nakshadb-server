package naksha

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type Request struct {
	Request   *http.Request
	Session   *sessions.Session
	UriParams map[string]string
}

func (r *Request) PostFormValue(key string) string {
	return r.Request.PostFormValue(key)
}
