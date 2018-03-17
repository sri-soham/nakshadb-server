package lib

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var store *sessions.FilesystemStore

const SESSION_KEY = "naksha"

func InitSessionStore(path string, auth_key []byte, enc_key []byte) {
	store = sessions.NewFilesystemStore(path, auth_key, enc_key)
}

func GetSessions(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, SESSION_KEY)
}
