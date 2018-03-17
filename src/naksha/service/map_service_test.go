package service_test

import (
	"naksha/service"
	"net/http"
	"net/url"
	"testing"
)

func TestMapAddInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Add(request)
	if result.IsSuccess() {
		t.Errorf("MapAddInvalid: success returned even when form is empty")
	}
}

func TestMapAddNoLayerUser(t *testing.T) {
	map_dao := &MapAddNoLayerUserImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("name", "City Map")
	values.Set("layer", "392")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Add(request)
	if result.IsSuccess() {
		t.Errorf("MapAddNoLayerUser: success returned when layer does not belong to user")
	}
}

func TestMapAdd(t *testing.T) {
	map_dao := &MapAddImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("name", "City Map")
	values.Set("layer", "24")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Add(request)
	if !result.IsSuccess() {
		t.Errorf("MapAdd: should be successful")
	}
	expected_redir_url := "/maps/120/show"
	redir_url := result.GetStringData("redir_url")
	if redir_url != expected_redir_url {
		t.Errorf("MapAdd: redir-url - expected(%v), found(%v)", expected_redir_url, redir_url)
	}
}

func TestMapDetails(t *testing.T) {
	map_dao := &MapDetailsImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := map[string]string{"id": "10"}
	values := url.Values{}
	request := makeRequest(http.MethodGet, uri_params, values, true)
	result := service.GetDetails(request)
	if !result.IsSuccess() {
		t.Errorf("MapDetails: should be successful")
	}
	_, ok := result.GetDataByKey("map_details")
	if !ok {
		t.Errorf("MapDetails: map_details is missing")
	}
	_, ok = result.GetDataByKey("base_layers")
	if !ok {
		t.Errorf("MapDetails: base_layers is missing")
	}
	_, ok = result.GetDataByKey("tables")
	if !ok {
		t.Errorf("MapDetails: tables is missing")
	}
}

func TestMapUpdateInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Update(request)
	if result.IsSuccess() {
		t.Errorf("MapUpdateInvalidForm: form validation not working")
	}
}

func TestMapUpdate(t *testing.T) {
	map_dao := &DaoImpl{}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	uri_params["id"] = "94"
	values := url.Values{}
	values.Set("name", "Updated Name")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Update(request)
	if !result.IsSuccess() {
		t.Errorf("MapUpdateInvalidForm: should be successful")
	}
}

func TestMapUserMaps(t *testing.T) {
	map_dao := &MapUserMapsImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodGet, uri_params, values, true)
	result := service.UserMaps(request)
	if !result.IsSuccess() {
		t.Errorf("MapUserMaps: expected true, found false")
	}
	_, ok := result.GetDataByKey("maps")
	if !ok {
		t.Errorf("MapUserMaps: maps is missing")
	}
	_, ok = result.GetDataByKey("pagination_links")
	if !ok {
		t.Errorf("MapUserMaps: pagination_links is missing")
	}
	_, ok = result.GetDataByKey("pagination_text")
	if !ok {
		t.Errorf("MapUserMaps: pagination_text is missing")
	}
}

func TestMapShowMap(t *testing.T) {
	map_dao := &MapShowMapImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := map[string]string{"map_hash": "1-lkwerldsf"}
	values := url.Values{}
	request := makeRequest(http.MethodGet, uri_params, values, true)
	result := service.ShowMap(request)
	if !result.IsSuccess() {
		t.Errorf("MapShowMap: expected true, found false")
	}
	_, ok := result.GetDataByKey("map_details")
	if !ok {
		t.Errorf("MapShowMap: map_details is missing")
	}

	_, ok = result.GetDataByKey("is_google_maps")
	if !ok {
		t.Errorf("MapShowMap: is_google_maps is missing")
	}

	_, ok = result.GetDataByKey("is_bing_maps")
	if !ok {
		t.Errorf("MapShowMap: is_bing_maps is missing")
	}

	_, ok = result.GetDataByKey("is_yandex_maps")
	if !ok {
		t.Errorf("MapShowMap: is_yandex_maps is missing")
	}

	_, ok = result.GetDataByKey("user_details")
	if !ok {
		t.Errorf("MapShowMap: user_details is missing")
	}

	_, ok = result.GetDataByKey("layer_data")
	if !ok {
		t.Errorf("MapShowMap: layer_data is missing")
	}

	_, ok = result.GetDataByKey("extents")
	if !ok {
		t.Errorf("MapShowMap: extents is missing")
	}
}

func TestMapDelete(t *testing.T) {
	map_dao := &MapDeleteImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := map[string]string{"id": "1"}
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Delete(request)
	if !result.IsSuccess() {
		t.Errorf("MapDelete: result should be success but is error")
	}
	expected_redir_url := "/maps/index"
	redir_url := result.GetStringData("redir_url")
	if redir_url != expected_redir_url {
		t.Errorf("MapDelete: redir_url - expected(%v), found(%v)", expected_redir_url, redir_url)
	}
}

func TestMapBaseLayerUpdateInvalidForm(t *testing.T) {
	map_dao := &MapBaseLayerUpdateImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.BaseLayerUpdate(request)
	if result.IsSuccess() {
		t.Errorf("MapBaseLayerUpdateInvalidForm: validation failed")
	}
}

func TestMapBaseLayerUpdateError(t *testing.T) {
	map_dao := &MapBaseLayerUpdateErrorImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := map[string]string{"id": "10"}
	values := url.Values{}
	values.Set("base_layer", "o-osm")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.BaseLayerUpdate(request)
	if result.IsSuccess() {
		t.Errorf("MapBaseLayerUpdateError: success returned even when update failed")
	}
}

func TestMapBaseLayerUpdate(t *testing.T) {
	map_dao := &MapBaseLayerUpdateImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := map[string]string{"id": "120"}
	values := url.Values{}
	values.Set("base_layer", "o-osm")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.BaseLayerUpdate(request)
	if !result.IsSuccess() {
		t.Errorf("MapBaseLayer: should be success but is error")
	}
}

func TestMapSearchTables(t *testing.T) {
	map_dao := &MapSearchTablesImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("table_name", "cens")
	request := makeRequest(http.MethodGet, uri_params, values, true)
	result := service.SearchTables(request)
	if !result.IsSuccess() {
		t.Errorf("MapSearchTables: error returned instead of success")
	}
	_, ok := result.GetDataByKey("tables")
	if !ok {
		t.Errorf("MapSearchTables: tables not present")
	}
}

func TestMapQueryData(t *testing.T) {
	map_dao := &DaoImpl{}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("query", "SELECT * from user_1_table_1")
	// 2017-10-27:
	// When method is http.MethodGet, form data is not being parsed.
	request := makeRequest(http.MethodPost, uri_params, values, false)
	result := service.QueryData(request)
	if !result.IsSuccess() {
		t.Errorf("MapQueryData: error returned instead of success. error: (%v), query: %s", result.GetErrors(), request.Request.FormValue("query"))
	}
	_, ok := result.GetDataByKey("data")
	if !ok {
		t.Errorf("MapQueryData: data is not present")
	}
}

func TestMapAddLayerInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.AddLayer(request)
	if result.IsSuccess() {
		t.Errorf("MapAddLayerInvalidForm: validation failed")
	}
}

func TestMapAddLayer(t *testing.T) {
	map_dao := &MapAddLayerImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	uri_params["id"] = "12"
	values := url.Values{}
	values.Set("layer_id", "21")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.AddLayer(request)
	if !result.IsSuccess() {
		t.Errorf("MapAddLayer: result should be success")
	}
}

func TestMapDeleteLayerInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.DeleteLayer(request)
	if result.IsSuccess() {
		t.Errorf("MapDeleteLayerInvalidForm: validation failed")
	}
}

func TestMapDeleteLayerOne(t *testing.T) {
	map_dao := &DaoImpl{}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := map[string]string{"id": "33"}
	values := url.Values{}
	values.Set("layer_id", "42")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.DeleteLayer(request)
	if result.IsSuccess() {
		t.Errorf("MapDeleteLayerOne: layer deleted even when there is only one layer for map")
	}
}

func TestMapDeleteLayer(t *testing.T) {
	map_dao := &MapDeleteLayerImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	uri_params := map[string]string{"id": "33"}
	values := url.Values{}
	values.Set("layer_id", "42")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.DeleteLayer(request)
	if !result.IsSuccess() {
		t.Errorf("MapDeleteLayer: delete layer failed")
	}
}

func TestMapBelongsToUserFalse(t *testing.T) {
	map_dao := &DaoImpl{}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	belongs := service.MapBelongsToUser("12", 21)
	if belongs {
		t.Errorf("MapBelongsToUserFalse: expected(false), found(%v)", belongs)
	}
}

func TestMapBelongsToUser(t *testing.T) {
	map_dao := &MapBelongsToUserImpl{&DaoImpl{}}
	api_dao := &DaoImpl{}
	service := service.MakeMapService(map_dao, api_dao)
	belongs := service.MapBelongsToUser("42", 32)
	if !belongs {
		t.Errorf("MapBelongsToUser: expected(true), found(%v)", belongs)
	}
}
