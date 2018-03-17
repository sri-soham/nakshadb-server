package controller

import (
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"naksha"
	"naksha/helper"
	"naksha/service"
	"net/http"
	"strconv"
)

type ExportsController struct {
	*BaseController
	export_service service.ExportService
}

func (ec *ExportsController) Status(w http.ResponseWriter, r *naksha.Request) {
	id := r.UriParams["id"]
	result := ec.export_service.GetStatus(id)
	if result.IsSuccess() {
		if result.GetBoolData("remove_export_id") {
			id_int, _ := strconv.Atoi(id)
			helper.RemoveExportIDFromSession(id_int, r.Session)
			r.Session.Save(r.Request, w)
		}
	}
	json_resp := result.ForJsonResponse()
	ec.RenderJson(w, json_resp)
}

func (ec *ExportsController) Download(w http.ResponseWriter, r *naksha.Request) {
	exports_dir := ec.GetAppConfig().ExportsDir()
	result := ec.export_service.Download(r.UriParams["id"], exports_dir)
	if result.IsSuccess() {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "Sat, 01 Jan 1994 00:10:20 GMT")
		w.Header().Set("Content-Type", result.GetStringData("content_type"))
		attachment := fmt.Sprintf("attachment; filename=\"%v\"", result.GetStringData("filename"))
		w.Header().Set("Content-Disposition", attachment)
		http.ServeFile(w, r.Request, result.GetStringData("file_path"))
	} else {
		ec.RenderHtmlErrors(w, result.GetErrors())
	}
}

func (ec *ExportsController) Index(w http.ResponseWriter, r *naksha.Request) {
	result := ec.export_service.ExportsOfUser(r)
	if result.IsSuccess() {
		ec.RenderHtml(w, r.Session, "exports_index", result.GetData())
	} else {
		ec.RenderHtmlErrors(w, result.GetErrors())
	}
}

func (ec *ExportsController) Delete(w http.ResponseWriter, r *naksha.Request) {
	exports_dir := ec.GetAppConfig().ExportsDir()
	result := ec.export_service.Delete(exports_dir, r.UriParams["id"])
	res := result.ForJsonResponse()
	ec.RenderJson(w, res)
}

func (ec *ExportsController) AuthAuth(req *http.Request, sess *sessions.Session, path_parts []string) bool {
	var is_valid bool
	var user_id int

	tmp := sess.Values["user_id"]
	if tmp == nil {
		is_valid = false
	} else {
		user_id = tmp.(int)
		is_valid = user_id > 0
	}
	if is_valid {
		if len(path_parts) == 3 {
			is_valid = ec.export_service.ExportBelongsToUser(path_parts[1], user_id)
		}
	}

	return is_valid
}

func GetExportsHandlers(app_config *naksha.AppConfig, tmpl *template.Template, repository service.Repository) []Route {
	routes := make([]Route, 0)
	var route Route

	validators := make(map[string]string)
	validators["id"] = "^[\\d]+$"

	controller := ExportsController{
		&BaseController{tmpl: tmpl, app_config: app_config},
		repository.GetExportService(),
	}

	route = Route{GET, "/exports/{id}/status", &controller, controller.Status, validators}
	routes = append(routes, route)

	route = Route{GET, "/exports/{id}/download", &controller, controller.Download, validators}
	routes = append(routes, route)

	route = Route{GET, "/exports/index", &controller, controller.Index, validators}
	routes = append(routes, route)

	route = Route{POST, "/exports/{id}/delete", &controller, controller.Delete, validators}
	routes = append(routes, route)

	return routes
}
