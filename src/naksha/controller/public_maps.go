package controller

import (
	"github.com/gorilla/sessions"
	"html/template"
	"naksha"
	"naksha/service"
	"net/http"
)

type PublicMapsController struct {
	*BaseController
	map_service service.MapService
}

func (pmc *PublicMapsController) ShowMap(w http.ResponseWriter, r *naksha.Request) {
	result := pmc.map_service.ShowMap(r)
	if result.IsSuccess() {
		data := result.GetData()
		pmc.RenderHtml(w, r.Session, "index_map", data)
	} else {
		errs := result.GetErrors()
		pmc.RenderHtmlErrors(w, errs)
	}
}

func (pmc *PublicMapsController) QueryData(w http.ResponseWriter, r *naksha.Request) {
	result := pmc.map_service.QueryData(r)
	json_resp := result.ForJsonResponse()
	pmc.RenderJsonP(w, r.Request, json_resp)
}

func (pmc *PublicMapsController) LayerOfTable(w http.ResponseWriter, r *naksha.Request) {
	tiler_url := pmc.GetAppConfig().TilerDomain()
	result := pmc.map_service.LayerOfTable(r, tiler_url)
	json_resp := result.ForJsonResponse()
	pmc.RenderJsonP(w, r.Request, json_resp)
}

func (c *PublicMapsController) AuthAuth(req *http.Request, sess *sessions.Session, path_parts []string) bool {
	return true
}

func GetPublicMapsHandlers(app_config *naksha.AppConfig, tmpl *template.Template, repository service.Repository) []Route {
	routes := make([]Route, 0)
	var route Route

	validators := make(map[string]string)
	validators["map_hash"] = "^[0-9]+[-][A-Za-z0-9_-]{1,64}$"
	validators["schema"] = "^[[:alnum:]]{8}"
	validators["table"] = "^[[:alnum:]_]{1,63}$"
	controller := PublicMapsController{&BaseController{tmpl: tmpl, app_config: app_config}, repository.GetMapService()}

	route = Route{GET, "/p/m/{map_hash}", &controller, controller.ShowMap, validators}
	routes = append(routes, route)

	route = Route{GET, "/p/s", &controller, controller.QueryData, validators}
	routes = append(routes, route)

	route = Route{GET, "/p/l/{schema}/{table}", &controller, controller.LayerOfTable, validators}
	routes = append(routes, route)

	return routes
}
