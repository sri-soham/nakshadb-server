package service_test

import (
	"naksha"
	"naksha/importer"
	"naksha/logger"
	"naksha/service"
	"net/http"
	"net/url"
	"testing"
)

func TestTableUserTables(t *testing.T) {
	map_dao := &DaoImpl{}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("page", "1")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.UserTables(request)
	if !result.IsSuccess() {
		t.Errorf("Could not fetch user tables")
	}
	_, ok := result.GetDataByKey("tables")
	if !ok {
		t.Errorf("tables not present in result")
	}
	_, ok = result.GetDataByKey("pagination_links")
	if !ok {
		t.Errorf("pagination_links not present in result")
	}
	_, ok = result.GetDataByKey("pagination_text")
	if !ok {
		t.Errorf("pagination_text not present in result")
	}
}

func TestTableCreateEmptyInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	request.Request.URL.RawQuery = "is_empty=1"
	app_config := naksha.MakeAppConfigTest(make(map[string]interface{}))
	result := service.CreateTable(request, &app_config)
	if result.IsSuccess() {
		t.Errorf("CreateEmpyTable: input validation failed")
	}
}

func TestTableCreateEmptySuccess(t *testing.T) {
	logger.InitLoggers("../../../test")
	map_dao := &TableEmptyTableDaoImpl{&DaoImpl{}}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("name", "SomeName")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	request.Request.URL.RawQuery = "is_empty=1"
	app_config := naksha.MakeAppConfigTest(make(map[string]interface{}))
	result := service.CreateTable(request, &app_config)
	if !result.IsSuccess() {
		t.Errorf("Valid request failed")
	}
	_, ok := result.GetDataByKey("url")
	if !ok {
		t.Errorf("url is not present in data")
	}
}

func TestTableCreateWithUploadInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	app_config := naksha.MakeAppConfigTest(make(map[string]interface{}))
	result := service.CreateTable(request, &app_config)
	if result.IsSuccess() {
		t.Errorf("input validation failed")
	}
}

func TestTableCreateWithUploadSuccess(t *testing.T) {
	logger.InitLoggers("../../../test")
	map_dao := &TableEmptyTableDaoImpl{&DaoImpl{}}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("name", "SomeName")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	app_config := naksha.MakeAppConfigTest(make(map[string]interface{}))
	result := service.CreateTable(request, &app_config)
	if !result.IsSuccess() {
		t.Errorf("Valid request failed")
	}
	_, ok := result.GetDataByKey("id")
	if !ok {
		t.Errorf("id is not present in data")
	}
}

func TestTableCheckStatusDBError(t *testing.T) {
	map_dao := &TableCheckStatusDBErrorDaoImpl{&DaoImpl{}}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	result := service.CheckStatus("10")
	t_status, _ := result["status"]
	t_err, _ := result["errors"]
	status := t_status.(string)
	err := t_err.(string)
	if status != "error" {
		t.Errorf("Expected (error), found(%s)", status)
	}
	if err != "Query Failed" {
		t.Errorf("Expected (Query Failed), found(%s)", err)
	}
}

func TestTableCheckStatusNoRows(t *testing.T) {
	map_dao := &TableCheckStatusNoRowsDaoImpl{&DaoImpl{}}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	result := service.CheckStatus("10")
	t_status, _ := result["status"]
	t_err, _ := result["errors"]
	status := t_status.(string)
	err := t_err.(string)
	if status != "error" {
		t.Errorf("Expected (error), found(%s)", status)
	}
	if err != "No such record" {
		t.Errorf("Expected (No such record), found(%s)", err)
	}
}

func TestTableCheckStatusReady(t *testing.T) {
	map_dao := &TableCheckStatusDaoImpl{&DaoImpl{}, importer.READY}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	result := service.CheckStatus("10")
	t_status, _ := result["status"]
	status := t_status.(string)
	if status != "success" {
		t.Errorf("Expected (success), found(%s)", status)
	}

	t_import_name, _ := result["import_name"]
	import_name := t_import_name.(string)
	if import_name != "Some Import" {
		t.Errorf("Expected (Some Import), found(%s)", import_name)
	}

	t_import_status, _ := result["import_status"]
	import_status := t_import_status.(string)
	if import_status != "success" {
		t.Errorf("Expected: (success), found(%s)", import_status)
	}

	t_table_url, _ := result["table_url"]
	table_url := t_table_url.(string)
	expected_table_url := "/tables/10/show"
	if table_url != expected_table_url {
		t.Errorf("Expected: (%s), found(%s)", expected_table_url, table_url)
	}

	t_remove_id := result["remove_import_id"]
	remove_id := t_remove_id.(string)
	if remove_id != "1" {
		t.Errorf("Expected (1), found(%s)", remove_id)
	}
}

func TestTableCheckStatusError(t *testing.T) {
	map_dao := &TableCheckStatusDaoImpl{&DaoImpl{}, importer.ERROR}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	result := service.CheckStatus("10")
	t_status, _ := result["status"]
	status := t_status.(string)
	if status != "success" {
		t.Errorf("Expected (success), found(%s)", status)
	}

	t_import_name, _ := result["import_name"]
	import_name := t_import_name.(string)
	if import_name != "Some Import" {
		t.Errorf("Expected (Some Import), found(%s)", import_name)
	}

	t_import_status, _ := result["import_status"]
	import_status := t_import_status.(string)
	if import_status != "error" {
		t.Errorf("Expected: (error), found(%s)", import_status)
	}

	t_err, _ := result["errors"]
	err := t_err.(string)
	if err != "Import failed" {
		t.Errorf("Expected (Import failed), found(%s)", err)
	}
}

func TestTableCheckStatusImporting(t *testing.T) {
	map_dao := &TableCheckStatusDaoImpl{&DaoImpl{}, importer.UPLOADED}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	result := service.CheckStatus("10")
	t_status, _ := result["status"]
	status := t_status.(string)
	if status != "success" {
		t.Errorf("Expected (success), found(%s)", status)
	}

	t_import_name, _ := result["import_name"]
	import_name := t_import_name.(string)
	if import_name != "Some Import" {
		t.Errorf("Expected (Some Import), found(%s)", import_name)
	}

	t_import_status, _ := result["import_status"]
	import_status := t_import_status.(string)
	if import_status != "importing" {
		t.Errorf("Expected: (importing), found(%s)", import_status)
	}
}

func TestTableGetDetails(t *testing.T) {
	map_dao := &TableGetDetailsDaoImpl{&DaoImpl{}}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := map[string]string{"id": "2"}
	values := url.Values{}
	request := makeRequest(http.MethodGet, uri_params, values, true)
	result := service.GetDetails(request, "http://tiler.my/")
	if !result.IsSuccess() {
		t.Errorf("Request failed")
	}

	_, ok := result.GetDataByKey("table_details")
	if !ok {
		t.Errorf("table_details not present in data")
	}

	_, ok = result.GetDataByKey("columns")
	if !ok {
		t.Errorf("columns not present in data")
	}

	tmp, ok := result.GetDataByKey("url")
	if !ok {
		t.Errorf("url not present in data")
	}
	expected_url := "/table_rows/2/"
	url := tmp.(string)
	if url != expected_url {
		t.Errorf("Url: expected(%s), found(%s)", expected_url, url)
	}

	tmp, ok = result.GetDataByKey("map_url")
	if !ok {
		t.Errorf("map_url not present in data")
	}

	expected_map_url := "http://tiler.my/lyr/hash123-[ts]/{z}/{x}/{y}.png"
	map_url := tmp.(string)
	if map_url != expected_map_url {
		t.Errorf("map-url: expected(%s), found(%s)", expected_map_url, map_url)
	}

	_, ok = result.GetDataByKey("extent")
	if !ok {
		t.Errorf("extent is not present in data")
	}

	tmp, ok = result.GetDataByKey("layer_id")
	if !ok {
		t.Errorf("layer_id is not present in data")
	}
	layer_id := tmp.(string)
	if layer_id != "3" {
		t.Errorf("layer_id mismatch: expected(3), found(%s)", layer_id)
	}

	_, ok = result.GetDataByKey("geometry_type")
	if !ok {
		t.Errorf("geometry_type is not present in data")
	}

	_, ok = result.GetDataByKey("style")
	if !ok {
		t.Errorf("style is not present in data")
	}

	_, ok = result.GetDataByKey("infowindow")
	if !ok {
		t.Errorf("infowindow is not present in data")
	}

	_, ok = result.GetDataByKey("update_hash")
	if !ok {
		t.Errorf("update_hash is not present in data")
	}

	_, ok = result.GetDataByKey("export_formats")
	if !ok {
		t.Errorf("export_formats is not present in data")
	}

	tmp, ok = result.GetDataByKey("tables_url")
	if !ok {
		t.Errorf("tables_url is not present in data")
	}
	expected_tables_url := "/tables/2"
	tables_url := tmp.(string)
	if tables_url != expected_tables_url {
		t.Errorf("tables_url: expected(%s), found(%s)", expected_tables_url, tables_url)
	}
}

func TestTableUpdateStylesInvalidForm(t *testing.T) {
	map_dao := &TableUpdateStylesDaoImpl{&DaoImpl{}}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := map[string]string{"id": "10"}
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.UpdateStyles(request)
	if result.IsSuccess() {
		t.Errorf("Validation failed")
	}
}

func TestTableUpdateStylesSuccess(t *testing.T) {
	map_dao := &TableUpdateStylesDaoImpl{&DaoImpl{}}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := map[string]string{"id": "10"}
	values := url.Values{}
	values.Set("stroke", "#ffffff")
	values.Set("stroke_opacity", "0.75")
	values.Set("stroke_width", "2.5")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.UpdateStyles(request)
	if !result.IsSuccess() {
		t.Errorf("Request failed")
	}
	_, ok := result.GetDataByKey("update_hash")
	if !ok {
		t.Errorf("update_hash is not present in data")
	}
}

func TestTableDeleteTable(t *testing.T) {
	map_dao := &TableDeleteColumnDaoImpl{&DaoImpl{}}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	result := service.DeleteTable("10")
	if !result.IsSuccess() {
		t.Errorf("Delete table failed")
	}
}

func TestTableAddColumnInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := map[string]string{"id": "10"}
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.AddColumn(request)
	if result.IsSuccess() {
		t.Errorf("validation failed")
	}
}

func TestTableAddColumnSuccess(t *testing.T) {
	map_dao := &TableAddColumnDaoImpl{&DaoImpl{}}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := map[string]string{"id": "10"}
	values := url.Values{}
	values.Set("name", "tbl_city")
	values.Set("data_type", "1")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.AddColumn(request)
	if !result.IsSuccess() {
		t.Errorf("Could not add column")
	}
}

func TestTableDeleteColumnNoValue(t *testing.T) {
	map_dao := &DaoImpl{}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.DeleteColumn(request)
	if result.IsSuccess() {
		t.Errorf("Request succeeding when column name is not given")
	}
}

func TestTableDeleteColumnReservedColumn(t *testing.T) {
	map_dao := &DaoImpl{}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("column_name", "naksha_id")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.DeleteColumn(request)
	if result.IsSuccess() {
		t.Errorf("Request succeeding when reserved column is being deleted")
	}
}

func TestTableDeleteColumnSuccess(t *testing.T) {
	map_dao := &TableDeleteColumnDaoImpl{&DaoImpl{}}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("column_name", "name")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.DeleteColumn(request)
	if !result.IsSuccess() {
		t.Errorf("DeleteColumn request failed")
	}
}

func TestTableInfoWindowProtectedColumns(t *testing.T) {
	map_dao := &DaoImpl{}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Add("columns", "name")
	values.Add("columns", "the_geom")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Infowindow(request)
	if result.IsSuccess() {
		t.Errorf("Protected columns included in infowindow")
	}
}

func TestTableInfoWindowSuccess(t *testing.T) {
	map_dao := &DaoImpl{}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := map[string]string{"id": "10"}
	values := url.Values{}
	values.Add("columns", "name")
	values.Add("columns", "stop_count")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Infowindow(request)
	if !result.IsSuccess() {
		t.Errorf("Modify infowindow request failed")
	}
}

// For valid format, tests are in the exporter package
func TestTableExportInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	app_config := naksha.MakeAppConfigTest(make(map[string]interface{}))
	result := service.Export(request, &app_config)
	if result.IsSuccess() {
		t.Errorf("Format validation failed")
	}
}

func TestTableApiAccessNoValue(t *testing.T) {
	map_dao := &DaoImpl{}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.ApiAccess(request, "dbuser2")
	if result.IsSuccess() {
		t.Errorf("Request succeeded when no value given for api_access")
	}
}

func TestTableApiAccessInvalidValue(t *testing.T) {
	map_dao := &DaoImpl{}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("api_access", "a")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.ApiAccess(request, "dbuser2")
	if result.IsSuccess() {
		t.Errorf("Request succeeded when invalid value given for api_access")
	}
}

func TestTableApiAccessSuccess(t *testing.T) {
	map_dao := &TableAddColumnDaoImpl{&DaoImpl{}}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	uri_params := map[string]string{"id": "10"}
	values := url.Values{}
	values.Set("api_access", "t")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.ApiAccess(request, "dbuser2")
	if !result.IsSuccess() {
		t.Errorf("Valid request failing")
	}
}

func TestTableBelongsToUserNo(t *testing.T) {
	map_dao := &TableBelongsToUserDaoImpl{&DaoImpl{}, 0}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	belongs := service.TableBelongsToUser("10", 2)
	if belongs {
		t.Errorf("Unauthorized user being given access")
	}
}

func TestTableBelongsToUserYes(t *testing.T) {
	map_dao := &TableBelongsToUserDaoImpl{&DaoImpl{}, 1}
	importer := &MockImporter{}
	service := service.MakeTableServiceWithImporter(map_dao, importer)
	belongs := service.TableBelongsToUser("10", 2)
	if !belongs {
		t.Errorf("Authorized not given access")
	}
}
