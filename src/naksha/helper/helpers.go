package helper

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"math/rand"
	"strings"
	"time"
)

const (
	EXPORTS_SESSION_KEY = "exports"
	IMPORTS_SESSION_KEY = "imports"
)

type Infowindow struct {
	Fields []string
}

func RandomString(of_length int) string {
	var j int
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	char_len := len(chars)
	seed := rand.NewSource(time.Now().UnixNano())
	rsource := rand.New(seed)
	rand_str := ""
	for i := 0; i < of_length; i++ {
		j = rsource.Intn(char_len)
		rand_str += string(chars[j])
	}

	return rand_str
}

func RandomKey(of_length int) string {
	var j int

	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	chars3 := chars + chars + chars
	charlen := len(chars3)
	seed := rand.NewSource(time.Now().UnixNano())
	rsource := rand.New(seed)
	rand_str := ""
	for i := 0; i < of_length; i++ {
		j = rsource.Intn(charlen)
		rand_str += string(chars3[j])
	}

	return rand_str
}

func RandomSchemaName(of_length int) string {
	var j int

	chars_first := "abcdefghijklmnopqrstuvwxyz"
	first_len := len(chars_first)
	chars_rest := "abcdefghijklmnopqrstuvwxyz0123456789"
	rest_len := len(chars_rest)
	seed := rand.NewSource(time.Now().UnixNano())
	rsource := rand.New(seed)
	rand_str := ""
	j = rsource.Intn(first_len)
	rand_str += string(chars_first[j])
	for i := 0; i < of_length-1; i++ {
		j = rsource.Intn(rest_len)
		rand_str += string(chars_rest[j])
	}

	return rand_str
}

// Assumes that user_id is present in session
func UserIDFromSession(sess *sessions.Session) int {
	user_id_str := sess.Values["user_id"]
	user_id := user_id_str.(int)

	return user_id
}

func SchemaNameFromSession(sess *sessions.Session) string {
	return sess.Values["schema_name"].(string)
}

func SchemaTableFromDetails(details map[string]string) string {
	return details["schema_name"] + "." + details["table_name"]
}

func HashForMapUrl() string {
	return fmt.Sprintf("%v_%v", RandomString(32), time.Now().Unix())
}

func StringToInfowindow(infowindow_str string) (Infowindow, error) {
	var iw Infowindow
	err := json.Unmarshal([]byte(infowindow_str), &iw)

	return iw, err
}

func InfowindowToString(ifw Infowindow) string {
	return "{\"fields\":[\"" + strings.Join(ifw.Fields, "\", \"") + "\"]}"
}

func IsGoogleMapsBaseLayer(base_layer string) bool {
	return string(base_layer[0]) == "g"
}

func IsBingMapsBaseLayer(base_layer string) bool {
	return string(base_layer[0]) == "b"
}

func IsYandexMapsBaseLayer(base_layer string) bool {
	return string(base_layer[0]) == "y"
}

func AddExportIDToSession(id int, sess *sessions.Session) {
	addElementToSessionArray(EXPORTS_SESSION_KEY, id, sess)
}

func RemoveExportIDFromSession(id int, sess *sessions.Session) {
	removeElementFromSessionArray(EXPORTS_SESSION_KEY, id, sess)
}

func GetExportIDsFromSession(sess *sessions.Session) []int {
	return elementsFromSessionArray(EXPORTS_SESSION_KEY, sess)
}

func AddImportIDToSession(id int, sess *sessions.Session) {
	addElementToSessionArray(IMPORTS_SESSION_KEY, id, sess)
}

func RemoveImportIDFromSession(id int, sess *sessions.Session) {
	removeElementFromSessionArray(IMPORTS_SESSION_KEY, id, sess)
}

func GetImportIDsFromSession(sess *sessions.Session) []int {
	return elementsFromSessionArray(IMPORTS_SESSION_KEY, sess)
}

func FormatTimestamp(t time.Time) string {
	tmp := fmt.Sprintf("%s", t)
	tmp = tmp[0:19]
	tmp = strings.Replace(tmp, "T", " ", 1)

	return tmp
}

func addElementToSessionArray(key string, id int, sess *sessions.Session) {
	tmp := sess.Values[key]
	if tmp == nil {
		values := []int{id}
		sess.Values[key] = values
	} else {
		values := tmp.([]int)
		values = append(values, id)
		sess.Values[key] = values
	}
}

func removeElementFromSessionArray(key string, id int, sess *sessions.Session) {
	tmp := sess.Values[key]
	if tmp == nil {
		return
	}
	old_values := tmp.([]int)
	new_values := make([]int, 0)
	for _, v := range old_values {
		if v != id {
			new_values = append(new_values, v)
		}
	}
	sess.Values[key] = new_values
}

func elementsFromSessionArray(key string, sess *sessions.Session) []int {
	tmp := sess.Values[key]
	if tmp == nil {
		return []int{}
	} else {
		return tmp.([]int)
	}
}
