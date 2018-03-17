package controller

import (
	"github.com/gorilla/sessions"
	"html/template"
	"naksha"
	"naksha/service"
	"net/http"
)

type TableRowsController struct {
	*BaseController
	table_row_service service.TableRowService
}

func (trc *TableRowsController) setTableRowService(trs service.TableRowService) {
	trc.table_row_service = trs
}

func (trc *TableRowsController) Data(w http.ResponseWriter, r *naksha.Request) {
	var data map[string]interface{}

	result := trc.table_row_service.Data(r)
	if result.IsSuccess() {
		data = result.GetData()
		data["status"] = "success"
	} else {
		data["status"] = "error"
		data["errors"] = result.GetErrors()
	}

	trc.RenderJson(w, data)
}

func (trc *TableRowsController) Update(w http.ResponseWriter, r *naksha.Request) {
	data := make(map[string]interface{})

	result := trc.table_row_service.Update(r)
	if result.IsSuccess() {
		data = result.GetData()
		data["status"] = "success"
	} else {
		data["status"] = "error"
		data["errors"] = result.GetErrors()
	}

	trc.RenderJson(w, data)
}

func (trc *TableRowsController) Add(w http.ResponseWriter, r *naksha.Request) {
	result := trc.table_row_service.Add(r)
	data := result.ForJsonResponse()
	trc.RenderJson(w, data)
}

func (trc *TableRowsController) Delete(w http.ResponseWriter, r *naksha.Request) {
	result := trc.table_row_service.Delete(r.UriParams["table_id"], r.UriParams["id"])
	data := result.ForJsonResponse()
	trc.RenderJson(w, data)
}

func (trc *TableRowsController) Show(w http.ResponseWriter, r *naksha.Request) {
	result := trc.table_row_service.Show(r.UriParams["table_id"], r.UriParams["id"])
	data := result.ForJsonResponse()
	trc.RenderJson(w, data)
}

func (trc *TableRowsController) AuthAuth(req *http.Request, sess *sessions.Session, path_parts []string) bool {
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
		is_valid = trc.table_row_service.TableBelongsToUser(path_parts[1], user_id)
	}

	return is_valid
}

func GetTableRowsHandlers(app_config *naksha.AppConfig, tmpl *template.Template, repository service.Repository) []Route {
	routes := make([]Route, 0)
	var route Route

	validators := make(map[string]string)
	//validators["table_name"] = "^user_([[:digit:]]+)_table_([[:alnum:]]+)$"
	validators["table_id"] = "^[\\d]+$"
	validators["page"] = "^[\\d]+$"
	validators["id"] = "^[\\d]+$"

	controller := TableRowsController{&BaseController{tmpl: tmpl, app_config: app_config}, nil}
	controller.setTableRowService(repository.GetTableRowService())

	route = Route{GET, "/table_rows/{table_id}/data/{page}", &controller, controller.Data, validators}
	routes = append(routes, route)

	route = Route{POST, "/table_rows/{table_id}/update/{id}", &controller, controller.Update, validators}
	routes = append(routes, route)

	route = Route{POST, "/table_rows/{table_id}/add", &controller, controller.Add, validators}
	routes = append(routes, route)

	route = Route{POST, "/table_rows/{table_id}/delete/{id}", &controller, controller.Delete, validators}
	routes = append(routes, route)

	route = Route{GET, "/table_rows/{table_id}/show/{id}", &controller, controller.Show, validators}
	routes = append(routes, route)

	return routes
}
