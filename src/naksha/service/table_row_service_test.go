package service_test

import (
	"naksha/logger"
	"naksha/service"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	logs_dir := "../../../test"
	logger.InitLoggers(logs_dir)
	code := m.Run()
	os.Exit(code)
}

// Test names are prefixed with TestRowTable instead of TestTableRow as is the
// convention of naming tests of a service based on the name of the service.
// As per the convention, tests of service table_row_service should have prefix
// TestTableRow but these tests have prefix TestRowTable. This is because, tests
// for the table_service.go start with TestTable; this makes it difficult to run
// tests only of table_service with "-run TestTable" because that prefix will
// match TestTableRow too. Hence, TestRowTable instead of TestTableRow which will
// allow running the tests of table_service and table_row_service with the "-run"
// switch.
func TestRowTableData(t *testing.T) {
	map_dao := &TableRowDataImpl{&DaoImpl{}}
	service := service.MakeTableRowService(map_dao)
	uri_params := map[string]string{"table_id": "1", "page": "1"}
	values := url.Values{}
	values.Set("order_column", "name")
	values.Set("order_type", "asc")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Data(request)
	if !result.IsSuccess() {
		t.Errorf("Error returned instead of success")
	}
	tmp1, ok := result.GetDataByKey("rows")
	if !ok {
		t.Errorf("rows is not present")
	}
	_, ok = result.GetDataByKey("count")
	if !ok {
		t.Errorf("count is not present")
	}
	rows := tmp1.([]map[string]string)
	_, ok = rows[0]["the_geom_webmercator"]
	if ok {
		t.Errorf("the_geom_webmercator should not be included in the data")
	}
}

func TestRowTableUpdateInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	service := service.MakeTableRowService(map_dao)
	uri_params := map[string]string{"table_id": "124"}
	values := url.Values{}
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Update(request)
	if result.IsSuccess() {
		t.Errorf("Form validation failed")
	}
}

func TestRowTableUpdateNonGeom(t *testing.T) {
	map_dao := &TableRowUpdateNonGeomImpl{&DaoImpl{}, false}
	service := service.MakeTableRowService(map_dao)
	uri_params := map[string]string{"id": "1", "table_id": "1"}
	values := url.Values{}
	values.Set("column", "name")
	values.Set("value", "Some Name")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Update(request)
	if !result.IsSuccess() {
		t.Errorf("Update non geom field failed. %s", result.GetErrors())
	}
	_, ok := result.GetDataByKey("update_hash")
	if ok {
		t.Errorf("update_hash should not be present for non geom field")
	}
	if !map_dao.IsUpdateCalled() {
		t.Errorf("Update has not been called on dao class")
	}
}

func TestRowTableUpdateGeomUpdateStyle(t *testing.T) {
	map_dao := &TableRowUpdateGeomImpl{&DaoImpl{}, "unknown", false, false, false}
	service := service.MakeTableRowService(map_dao)
	uri_params := map[string]string{"id": "1", "table_id": "1"}
	values := url.Values{}
	values.Set("column", "the_geom")
	values.Set("value", "SRID=4326;POINT(78.4622620046139 17.411652235651)")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Update(request)

	// wait for a second to make sure that trs.updateGeometryTypeStyle is called
	time.Sleep(1 * time.Millisecond)
	if !result.IsSuccess() {
		t.Errorf("Update geom field failed. %s", result.GetErrors())
	}
	_, ok := result.GetDataByKey("update_hash")
	if !ok {
		t.Errorf("update_hash should be present for the_geom field")
	}
	if !map_dao.IsUpdateGeometryCalled() {
		t.Errorf("UpdateGeometry has not been called on dao class")
	}
	if !map_dao.IsFindWhereCalled() {
		t.Errorf("Geometry Type in layer table has not been queried for")
	}
	if !map_dao.IsFetchDataCalled() {
		t.Errorf("Geometry type has not been queried for on the table")
	}
}

func TestRowTableUpdateGeomNoStyle(t *testing.T) {
	map_dao := &TableRowUpdateGeomImpl{&DaoImpl{}, "polygon", false, false, false}
	service := service.MakeTableRowService(map_dao)
	uri_params := map[string]string{"id": "1", "table_id": "1"}
	values := url.Values{}
	values.Set("column", "the_geom")
	values.Set("value", "SRID=4326;POINT(78.4622620046139 17.411652235651)")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Update(request)

	// wait for a second to make sure that trs.updateGeometryTypeStyle is called
	time.Sleep(1 * time.Millisecond)
	if !result.IsSuccess() {
		t.Errorf("Update geom field failed. %s", result.GetErrors())
	}
	_, ok := result.GetDataByKey("update_hash")
	if !ok {
		t.Errorf("update_hash should be present for the_geom field")
	}
	if !map_dao.IsUpdateGeometryCalled() {
		t.Errorf("UpdateGeometry has not been called on dao class")
	}
	if !map_dao.IsFindWhereCalled() {
		t.Errorf("FindWhere has not been called on dao class")
	}
	if map_dao.IsFetchDataCalled() {
		t.Errorf("Style has been updated when it shouldn't be")
	}
}

func TestRowTableAddInvalidForm(t *testing.T) {
	map_dao := &DaoImpl{}
	service := service.MakeTableRowService(map_dao)
	uri_params := map[string]string{"table_id": "42"}
	values := url.Values{}
	values.Set("with_geometry", "1")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Add(request)
	if result.IsSuccess() {
		t.Errorf("Input validation failed")
	}
}

func TestRowTableAddWithoutGeometry(t *testing.T) {
	map_dao := &TableRowAddNonGeomImpl{&DaoImpl{}, false}
	service := service.MakeTableRowService(map_dao)
	uri_params := map[string]string{"table_id": "42"}
	values := url.Values{}
	values.Set("with_geometry", "0")
	values.Set("name", "Some Name")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Add(request)
	if !result.IsSuccess() {
		t.Errorf("Error returned")
	}
	_, ok := result.GetDataByKey("row")
	if !ok {
		t.Errorf("row should be present in data")
	}
	if !map_dao.IsInsertCalled() {
		t.Errorf("Insert has not been called")
	}
}

func TestRowTableAddGeometryNoStyle(t *testing.T) {
	map_dao := &TableRowAddGeomImpl{&DaoImpl{}, "polygon", false, false, false}
	service := service.MakeTableRowService(map_dao)
	uri_params := map[string]string{"table_id": "42"}
	values := url.Values{}
	values.Set("with_geometry", "1")
	values.Set("geometry", "SRID=4326;POINT(78.4622620046139 17.411652235651)")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Add(request)
	// wait for a second to make sure that trs.updateGeometryTypeStyle is called
	time.Sleep(100 * time.Millisecond)
	if !result.IsSuccess() {
		t.Errorf("Error returned")
	}
	_, ok := result.GetDataByKey("row")
	if !ok {
		t.Errorf("row should be present in data")
	}
	if !map_dao.insert_geometry_called {
		t.Errorf("InsertWithGeometry not called")
	}
	if !map_dao.find_where_called {
		t.Errorf("FindWhere not called")
	}
	if map_dao.fetch_data_called {
		t.Errorf("Geometry type queried for even when geometry type is set in mstr_layer")
	}
}

func TestRowTableAddGeometryUpdateStyle(t *testing.T) {
	map_dao := &TableRowAddGeomImpl{&DaoImpl{}, "unknown", false, false, false}
	service := service.MakeTableRowService(map_dao)
	uri_params := map[string]string{"table_id": "42"}
	values := url.Values{}
	values.Set("with_geometry", "1")
	values.Set("geometry", "SRID=4326;POINT(78.4622620046139 17.411652235651)")
	request := makeRequest(http.MethodPost, uri_params, values, true)
	result := service.Add(request)
	// wait for a second to make sure that trs.updateGeometryTypeStyle is called
	time.Sleep(100 * time.Millisecond)
	if !result.IsSuccess() {
		t.Errorf("Error returned")
	}
	_, ok := result.GetDataByKey("row")
	if !ok {
		t.Errorf("row should be present in data")
	}
	if !map_dao.insert_geometry_called {
		t.Errorf("InsertWithGeometry not called")
	}
	if !map_dao.find_where_called {
		t.Errorf("mstr_layer not queried for geometry_type")
	}
	if !map_dao.fetch_data_called {
		t.Errorf("Geometry type not queried for when geometry type is 'unknown'in mstr_layer")
	}
}

func TestRowTableDelete(t *testing.T) {
	map_dao := &TableRowDeleteImpl{&DaoImpl{}}
	service := service.MakeTableRowService(map_dao)
	result := service.Delete("323", "12")
	if !result.IsSuccess() {
		t.Errorf("Delete failed")
	}
	_, ok := result.GetDataByKey("update_hash")
	if !ok {
		t.Errorf("update_hash not present in data")
	}
}

func TestRowTableShow(t *testing.T) {
	map_dao := &TableRowShowImpl{&DaoImpl{}}
	service := service.MakeTableRowService(map_dao)
	result := service.Show("523", "92")
	if !result.IsSuccess() {
		t.Errorf("Could not get details")
	}
	_, ok := result.GetDataByKey("data")
	if !ok {
		t.Errorf("data is not present in result")
	}
}

func TestRowTableTableBelongsToUserFalse(t *testing.T) {
	map_dao := &DaoImpl{}
	service := service.MakeTableRowService(map_dao)
	belongs := service.TableBelongsToUser("143", 3)
	if belongs {
		t.Errorf("TableBelongsToUser: unauthorized user given access")
	}
}

func TestRowTableTableBelongsToUserTrue(t *testing.T) {
	map_dao := &TableRowBelongsImpl{&DaoImpl{}}
	service := service.MakeTableRowService(map_dao)
	belongs := service.TableBelongsToUser("1432", 2)
	if !belongs {
		t.Errorf("TableBelongsToUser: authorized user denied access")
	}
}
