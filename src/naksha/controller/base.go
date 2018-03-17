package controller

import (
	//    "bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"naksha"
	"naksha/helper"
	"naksha/logger"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	GET  = "GET"
	POST = "POST"
)

const ERROR_HTML = `<!DOCTYPE html>
<html>
  <head>
    <title>Naksha: Error</title>
    <style type="text/css">
      #errors-div {
        width: 500px;
        margin: 50px auto 0px auto;
        padding: 10px;
        border: 1px #ccc solid;
        color: #000;
        background-color: #fff;
        font-size: 16px;
      }
      #errors-title {
        font-size: 24px;
        color: #f00;
        text-align: center;
        margin-bottom: 10px;
      }
    </style>
  </head>
  <body>
    <div id="errors-div">
      <p id="errors-title">Error(s):</p>
      {errors}
    </div>
  </body>
</html>`

type BaseController struct {
	tmpl       *template.Template
	app_config *naksha.AppConfig
}

func (c *BaseController) RenderHtmlString(w http.ResponseWriter, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html)
}

func (c *BaseController) RenderHtml(w http.ResponseWriter, sess *sessions.Session, tmpl_name string, data map[string]interface{}) {
	c.addSessionVarsToTemplateData(sess, data)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := c.tmpl.ExecuteTemplate(w, tmpl_name, data)
	if err != nil {
		c.InternalServerError(w, err)
	}
}

func (c *BaseController) RenderHtmlNoData(w http.ResponseWriter, sess *sessions.Session, tmpl_name string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := make(map[string]interface{})
	c.addSessionVarsToTemplateData(sess, data)
	err := c.tmpl.ExecuteTemplate(w, tmpl_name, data)
	if err != nil {
		c.InternalServerError(w, err)
	}
}

func (c *BaseController) RenderHtmlErrors(w http.ResponseWriter, errors []string) {
	err_str := strings.Join(errors, "<br />")
	html_str := strings.Replace(ERROR_HTML, "{errors}", err_str, -1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html_str)
}

func (c *BaseController) RenderText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, text)
}

func (c *BaseController) RenderJson(w http.ResponseWriter, data map[string]interface{}) {
	js, err := json.Marshal(data)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (c *BaseController) RenderJsonP(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	js, err := json.Marshal(data)
	if err == nil {
		callback := r.URL.Query().Get("callback")
		if len(callback) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else {
			//n := bytes.IndexByte(js, 0)
			//s := string(js[:n])
			s := string(js[:])
			s = callback + "(" + s + ")"
			w.Header().Set("Content-Type", "application/javascript")
			fmt.Fprintf(w, s)
		}
	}
}

func (c *BaseController) GetTemplate() *template.Template {
	return c.tmpl
}

func (c *BaseController) GetAppConfig() *naksha.AppConfig {
	return c.app_config
}

func (c *BaseController) NotFoundError(w http.ResponseWriter, err error) {
	logger.LogMessage(err)
	http.Error(w, "404 - Not Found", http.StatusNotFound)
}

func (c *BaseController) InternalServerError(w http.ResponseWriter, err error) {
	logger.LogMessage(err)
	http.Error(w, "503 - Internal Server Error", http.StatusInternalServerError)
}

func (c *BaseController) AuthAuth(req *http.Request, sess *sessions.Session, path_parts []string) bool {
	return false
}

func (c *BaseController) addSessionVarsToTemplateData(sess *sessions.Session, data map[string]interface{}) {
	if tmp, ok := sess.Values["user_id"]; ok {
		data["sess_user_id"] = tmp
	}

	exports, ok := sess.Values[helper.EXPORTS_SESSION_KEY]
	if !ok || len(exports.([]int)) == 0 {
		data["sess_exports"] = make([]int, 0)
	} else {
		data["sess_exports"] = exports
	}
	imports, ok := sess.Values[helper.IMPORTS_SESSION_KEY]
	if !ok || len(imports.([]int)) == 0 {
		data["sess_imports"] = make([]int, 0)
	} else {
		data["sess_imports"] = imports
	}
}

type IController interface {
	AuthAuth(req *http.Request, sess *sessions.Session, path_parts []string) bool
}

type Handler func(http.ResponseWriter, *naksha.Request)

// This will be used during development. On production, urls starting with /assets/
// should be handled by the web server.
func RenderAsset(w http.ResponseWriter, r *http.Request, assets_path string, parts []string) {
	var content_type string

	asset_path_parts := parts[2:]
	full_path := assets_path + strings.Join(asset_path_parts, "/")
	if _, err := os.Stat(full_path); !os.IsNotExist(err) {
		ext := filepath.Ext(parts[len(parts)-1])
		ext = strings.ToLower(ext)
		ext = strings.TrimLeft(ext, ".")

		switch ext {
		case "css":
			content_type = "text/css"
		case "Js":
			content_type = "application/javascript"
		case "jpeg", "jpg", "jpe":
			content_type = "image/jpeg"
		case "png":
			content_type = "image/png"
		case "gif":
			content_type = "image/gif"
		default:
			content_type = "text/plain"
		}

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "Sat, 01 Jan 1994 00:10:20 GMT")
		w.Header().Set("Content-Type", content_type)
		http.ServeFile(w, r, full_path)
	} else {
		http.NotFound(w, r)
	}
}
