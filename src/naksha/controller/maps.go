package controller

import (
	"github.com/gorilla/sessions"
	"html/template"
	"naksha"
	"naksha/service"
	"net/http"
)

type MapsController struct {
	*BaseController
	map_service service.MapService
}

func (mc *MapsController) Index(w http.ResponseWriter, r *naksha.Request) {
	result := mc.map_service.UserMaps(r)
	data := result.GetData()
	mc.RenderHtml(w, r.Session, "maps_index", data)
}

func (mc *MapsController) NewPost(w http.ResponseWriter, r *naksha.Request) {
	result := mc.map_service.Add(r)
	data := result.ForJsonResponse()
	mc.RenderJson(w, data)
}

func (mc *MapsController) Show(w http.ResponseWriter, r *naksha.Request) {
	result := mc.map_service.GetDetails(r)
	data := result.GetData()
	data["js"] = []string{"jquery.form.min.js", "map_admin.js"}
	mc.RenderHtml(w, r.Session, "maps_show", data)
}

func (mc *MapsController) EditPost(w http.ResponseWriter, r *naksha.Request) {
	result := mc.map_service.Update(r)
	data := result.ForJsonResponse()
	mc.RenderJson(w, data)
}

func (mc *MapsController) DeletePost(w http.ResponseWriter, r *naksha.Request) {
	result := mc.map_service.Delete(r)
	data := result.ForJsonResponse()
	mc.RenderJson(w, data)
}

func (mc *MapsController) BaseLayerPost(w http.ResponseWriter, r *naksha.Request) {
	result := mc.map_service.BaseLayerUpdate(r)
	data := result.ForJsonResponse()
	mc.RenderJson(w, data)
}

func (mc *MapsController) SearchTables(w http.ResponseWriter, r *naksha.Request) {
	result := mc.map_service.SearchTables(r)
	data := result.ForJsonResponse()
	mc.RenderJson(w, data)
}

func (mc *MapsController) AddLayer(w http.ResponseWriter, r *naksha.Request) {
	result := mc.map_service.AddLayer(r)
	data := result.ForJsonResponse()
	mc.RenderJson(w, data)
}

func (mc *MapsController) DeleteLayer(w http.ResponseWriter, r *naksha.Request) {
	result := mc.map_service.DeleteLayer(r)
	data := result.ForJsonResponse()
	mc.RenderJson(w, data)
}

func (mc *MapsController) UpdateHash(w http.ResponseWriter, r *naksha.Request) {
	result := mc.map_service.UpdateHash(r)
	data := result.ForJsonResponse()
	mc.RenderJson(w, data)
}

func (mc *MapsController) AuthAuth(req *http.Request, sess *sessions.Session, path_parts []string) bool {
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
			is_valid = mc.map_service.MapBelongsToUser(path_parts[1], user_id)
		}
	}

	return is_valid
}

func GetMapsHandlers(app_config *naksha.AppConfig, tmpl *template.Template, repository service.Repository) []Route {
	routes := make([]Route, 0)
	var route Route

	validators := make(map[string]string)
	validators["id"] = "^[\\d]+$"

	controller := MapsController{
		&BaseController{tmpl: tmpl, app_config: app_config},
		repository.GetMapService(),
	}

	route = Route{GET, "/maps/index", &controller, controller.Index, validators}
	routes = append(routes, route)

	route = Route{POST, "/maps/new", &controller, controller.NewPost, validators}
	routes = append(routes, route)

	route = Route{GET, "/maps/{id}/show", &controller, controller.Show, validators}
	routes = append(routes, route)

	route = Route{POST, "/maps/{id}/edit", &controller, controller.EditPost, validators}
	routes = append(routes, route)

	route = Route{POST, "/maps/{id}/delete", &controller, controller.DeletePost, validators}
	routes = append(routes, route)

	route = Route{POST, "/maps/{id}/base_layer", &controller, controller.BaseLayerPost, validators}
	routes = append(routes, route)

	route = Route{GET, "/maps/{id}/search_tables", &controller, controller.SearchTables, validators}
	routes = append(routes, route)

	route = Route{POST, "/maps/{id}/add_layer", &controller, controller.AddLayer, validators}
	routes = append(routes, route)

	route = Route{POST, "/maps/{id}/delete_layer", &controller, controller.DeleteLayer, validators}
	routes = append(routes, route)

	route = Route{POST, "/maps/{id}/hash", &controller, controller.UpdateHash, validators}
	routes = append(routes, route)

	return routes
}
