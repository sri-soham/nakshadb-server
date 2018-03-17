package controller

import (
	"github.com/gorilla/sessions"
	"html/template"
	"naksha"
	"naksha/service"
	"net/http"
)

type IndexController struct {
	*BaseController
	user_service  service.UserService
	table_service service.TableService
}

func (i *IndexController) setUserService(us service.UserService) {
	i.user_service = us
}

func (i *IndexController) setTableService(ts service.TableService) {
	i.table_service = ts
}

func (i *IndexController) Dashboard(w http.ResponseWriter, r *naksha.Request) {
	result := i.table_service.UserTables(r)
	if result.IsSuccess() {
		data := result.GetData()
		data["js"] = []string{"dashboard.js"}
		i.RenderHtml(w, r.Session, "index_dashboard", data)
	} else {
		i.InternalServerError(w, result.GetError())
	}
}

func (i *IndexController) Login(w http.ResponseWriter, r *naksha.Request) {
	i.RenderHtmlNoData(w, r.Session, "index_index")
}

func (i *IndexController) LoginPost(w http.ResponseWriter, r *naksha.Request) {
	result := i.user_service.Login(r)
	if result.IsSuccess() {
		r.Session.Values["user_id"] = result.GetIntData("id")
		r.Session.Values["name"] = result.GetStringData("name")
		r.Session.Values["schema_name"] = result.GetStringData("schema_name")
		r.Session.Save(r.Request, w)
		http.Redirect(w, r.Request, "/dashboard", http.StatusSeeOther)
	} else {
		if result.IsFatalError() {
			i.InternalServerError(w, result.GetError())
		} else {
			data := make(map[string]interface{})
			data["errors"] = result.GetErrors()
			i.RenderHtml(w, r.Session, "index_index", data)
		}
	}
}

func (i *IndexController) Profile(w http.ResponseWriter, r *naksha.Request) {
	res := i.user_service.UserDetails(r)
	if res.IsSuccess() {
		data := res.GetData()
		data["js"] = []string{"jquery.form.min.js", "user.js"}
		i.RenderHtml(w, r.Session, "index_profile", res.GetData())
	} else {
		i.RenderHtmlErrors(w, res.GetErrors())
	}
}

func (i *IndexController) ChangePasswordPost(w http.ResponseWriter, r *naksha.Request) {
	result := i.user_service.ChangePassword(r)
	data := result.ForJsonResponse()
	i.RenderJson(w, data)
}

func (i *IndexController) ProfilePost(w http.ResponseWriter, r *naksha.Request) {
	result := i.user_service.ProfilePost(r)
	data := result.ForJsonResponse()
	i.RenderJson(w, data)
}

func (i *IndexController) Logout(w http.ResponseWriter, r *naksha.Request) {
	keys := make([]interface{}, 0, len(r.Session.Values))
	for k, _ := range r.Session.Values {
		keys = append(keys, k)
	}
	for _, k := range keys {
		delete(r.Session.Values, k)
	}
	r.Session.Save(r.Request, w)
	http.Redirect(w, r.Request, "/", http.StatusSeeOther)
}

func (i *IndexController) AuthAuth(req *http.Request, sess *sessions.Session, path_parts []string) bool {
	var is_valid bool

	switch path_parts[0] {
	case "dashboard", "logout":
		tmp := sess.Values["user_id"]
		if tmp == nil {
			is_valid = false
		} else {
			user_id := tmp.(int)
			is_valid = user_id > 0
		}
	default:
		is_valid = true
	}

	return is_valid
}

func GetIndexHandlers(app_config *naksha.AppConfig, tmpl *template.Template, repository service.Repository) []Route {
	routes := make([]Route, 0)
	var route Route

	validators := make(map[string]string)
	validators["map_hash"] = "^[[:alnum:]]{64}$"
	controller := IndexController{&BaseController{tmpl: tmpl, app_config: app_config}, nil, nil}
	controller.setUserService(repository.GetUserService())
	controller.setTableService(repository.GetTableService())

	route = Route{GET, "/", &controller, controller.Login, validators}
	routes = append(routes, route)
	route = Route{POST, "/", &controller, controller.LoginPost, validators}
	routes = append(routes, route)
	route = Route{GET, "/logout", &controller, controller.Logout, validators}
	routes = append(routes, route)
	route = Route{GET, "/dashboard", &controller, controller.Dashboard, validators}
	routes = append(routes, route)
	route = Route{GET, "/profile", &controller, controller.Profile, validators}
	routes = append(routes, route)
	route = Route{POST, "/change_password", &controller, controller.ChangePasswordPost, validators}
	routes = append(routes, route)
	route = Route{POST, "/profile", &controller, controller.ProfilePost, validators}
	routes = append(routes, route)

	return routes
}
