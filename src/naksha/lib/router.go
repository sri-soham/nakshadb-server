package lib

import (
	"html/template"
	"naksha"
	"naksha/controller"
	"naksha/service"
	"net/http"
	"strings"
)

type Router struct {
	routes      []controller.Route
	assets_path string
}

func (r *Router) ParseAndHandle(resp http.ResponseWriter, req *http.Request) {
	var route controller.Route
	var found bool
	var uri_params map[string]string

	parts := strings.Split(req.URL.Path[1:], "/")
	if len(req.URL.Path) > 8 && req.URL.Path[0:8] == "/assets/" {
		controller.RenderAsset(resp, req, r.assets_path, parts)
	} else {
		found = false
		for _, route = range r.routes {
			if uri_params, found = route.IsMatch(req); found {
				break
			}
		}
		if found {
			controller := route.GetController()
			handler := route.GetHandler()

			// get the sessions
			// https://github.com/gorilla/sessions/issues/16
			// http://www.gorillatoolkit.org/pkg/sessions#CookieStore.Get
			// if session is found but the auth_key and enc_key have changed
			// then error is returned. We don't have throw any error in that
			// scenario. New empty session will be generated and will be used
			// as required
			sess, _ := GetSessions(req)

			if !controller.AuthAuth(req, sess, parts) {
				ForbiddenRequest(resp)
				return
			}

			naksha_req := &naksha.Request{Request: req, Session: sess, UriParams: uri_params}
			handler(resp, naksha_req)

		} else {
			NotFound(resp)
			return
		}
	}
}

func (r *Router) setAssetsPath(path string) {
	path = strings.TrimSpace(path)
	path = strings.TrimRight(path, "/")
	path += "/"

	r.assets_path = path
}

func (r *Router) GetAssetsPath() string {
	return r.assets_path
}

func MakeRouter(app_config *naksha.AppConfig, tmpl *template.Template, repository service.Repository) *Router {
	router := Router{}
	var routes []controller.Route

	routes = controller.GetIndexHandlers(app_config, tmpl, repository)
	router.routes = append(router.routes, routes...)

	routes = controller.GetTablesHandlers(app_config, tmpl, repository)
	router.routes = append(router.routes, routes...)

	routes = controller.GetTableRowsHandlers(app_config, tmpl, repository)
	router.routes = append(router.routes, routes...)

	routes = controller.GetMapsHandlers(app_config, tmpl, repository)
	router.routes = append(router.routes, routes...)

	routes = controller.GetExportsHandlers(app_config, tmpl, repository)
	router.routes = append(router.routes, routes...)

	routes = controller.GetPublicMapsHandlers(app_config, tmpl, repository)
	router.routes = append(router.routes, routes...)

	router.setAssetsPath(app_config.AssetsDir())

	return &router
}

func NotFound(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusNotFound)
	resp.Write([]byte("404 - Not Found"))
}

func ForbiddenRequest(resp http.ResponseWriter) {
	http.Error(resp, "Forbidden Request", http.StatusForbidden)
}

func InternalServerError(resp http.ResponseWriter) {
	http.Error(resp, "Internal Server Error", http.StatusInternalServerError)
}
