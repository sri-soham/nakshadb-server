package controller

import (
	"github.com/gorilla/sessions"
	"html/template"
	"naksha"
	"naksha/helper"
	"naksha/service"
	"net/http"
	"strconv"
)

type TablesController struct {
	*BaseController
	table_service service.TableService
}

func (t *TablesController) setTableService(ts service.TableService) {
	t.table_service = ts
}

func (t *TablesController) New(w http.ResponseWriter, r *naksha.Request) {
	data := make(map[string]interface{})
	data["js"] = []string{"jquery.form.min.js", "new_table.js"}
	t.RenderHtml(w, r.Session, "tables_new", data)
}

func (t *TablesController) NewPost(w http.ResponseWriter, r *naksha.Request) {
	result := t.table_service.CreateTable(r, t.BaseController.GetAppConfig())
	if result.IsSuccess() {
		id, ok := result.GetDataByKey("id")
		if ok {
			helper.AddImportIDToSession(id.(int), r.Session)
			r.Session.Save(r.Request, w)
		}
	}
	json_resp := result.ForJsonResponse()
	t.RenderJson(w, json_resp)
}

func (t *TablesController) Status(w http.ResponseWriter, r *naksha.Request) {
	id := r.UriParams["id"]
	result := t.table_service.CheckStatus(id)
	if _, ok := result["remove_import_id"]; ok {
		id_int, _ := strconv.Atoi(id)
		helper.RemoveImportIDFromSession(id_int, r.Session)
		r.Session.Save(r.Request, w)
	}
	t.RenderJson(w, result)
}

func (t *TablesController) Show(w http.ResponseWriter, r *naksha.Request) {
	tiler_url := t.GetAppConfig().TilerDomain()
	result := t.table_service.GetDetails(r, tiler_url)
	if result.IsSuccess() {
		data := result.GetData()
		data["js"] = []string{"jquery.form.min.js", "table.js", "table_admin.js"}
		t.RenderHtml(w, r.Session, "tables_show", data)
	} else {
		t.RenderHtmlErrors(w, result.GetErrors())
	}
}

func (t *TablesController) Styles(w http.ResponseWriter, r *naksha.Request) {
	result := t.table_service.UpdateStyles(r)
	json_resp := result.ForJsonResponse()
	t.RenderJson(w, json_resp)
}

func (t *TablesController) Delete(w http.ResponseWriter, r *naksha.Request) {
	result := t.table_service.DeleteTable(r.UriParams["id"])
	json_resp := result.ForJsonResponse()
	t.RenderJson(w, json_resp)
}

func (t *TablesController) AddColumn(w http.ResponseWriter, r *naksha.Request) {
	result := t.table_service.AddColumn(r)
	json_resp := result.ForJsonResponse()
	t.RenderJson(w, json_resp)
}

func (t *TablesController) DeleteColumn(w http.ResponseWriter, r *naksha.Request) {
	result := t.table_service.DeleteColumn(r)
	json_resp := result.ForJsonResponse()
	t.RenderJson(w, json_resp)
}

func (t *TablesController) Infowindow(w http.ResponseWriter, r *naksha.Request) {
	result := t.table_service.Infowindow(r)
	json_resp := result.ForJsonResponse()
	t.RenderJson(w, json_resp)
}

func (t *TablesController) Export(w http.ResponseWriter, r *naksha.Request) {
	result := t.table_service.Export(r, t.BaseController.GetAppConfig())
	if result.IsSuccess() {
		id := result.GetIntData("id")
		helper.AddExportIDToSession(id, r.Session)
		r.Session.Save(r.Request, w)
	}
	json_resp := result.ForJsonResponse()
	t.RenderJson(w, json_resp)
}

func (t *TablesController) ApiAccess(w http.ResponseWriter, r *naksha.Request) {
	result := t.table_service.ApiAccess(r, t.GetAppConfig().DBApiUser())
	json_resp := result.ForJsonResponse()
	t.RenderJson(w, json_resp)
}

func (t *TablesController) AuthAuth(req *http.Request, sess *sessions.Session, path_parts []string) bool {
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
			is_valid = t.table_service.TableBelongsToUser(path_parts[1], user_id)
		}
	}

	return is_valid
}

func GetTablesHandlers(app_config *naksha.AppConfig, tmpl *template.Template, repository service.Repository) []Route {
	routes := make([]Route, 0)
	var route Route
	empty_validators := make(map[string]string)
	validators := make(map[string]string)
	validators["id"] = "^[\\d]+$"

	controller := TablesController{&BaseController{tmpl: tmpl, app_config: app_config}, nil}
	controller.setTableService(repository.GetTableService())

	route = Route{GET, "/tables/new", &controller, controller.New, empty_validators}
	routes = append(routes, route)

	route = Route{POST, "/tables/new", &controller, controller.NewPost, empty_validators}
	routes = append(routes, route)

	route = Route{GET, "/tables/{id}/status", &controller, controller.Status, validators}
	routes = append(routes, route)

	route = Route{GET, "/tables/{id}/show", &controller, controller.Show, validators}
	routes = append(routes, route)

	route = Route{POST, "/tables/{id}/styles", &controller, controller.Styles, validators}
	routes = append(routes, route)

	route = Route{POST, "/tables/{id}/delete", &controller, controller.Delete, validators}
	routes = append(routes, route)

	route = Route{POST, "/tables/{id}/add_column", &controller, controller.AddColumn, validators}
	routes = append(routes, route)

	route = Route{POST, "/tables/{id}/delete_column", &controller, controller.DeleteColumn, validators}
	routes = append(routes, route)

	route = Route{POST, "/tables/{id}/infowindow", &controller, controller.Infowindow, validators}
	routes = append(routes, route)

	route = Route{POST, "/tables/{id}/export", &controller, controller.Export, validators}
	routes = append(routes, route)

	route = Route{POST, "/tables/{id}/api_access", &controller, controller.ApiAccess, validators}
	routes = append(routes, route)

	return routes
}
