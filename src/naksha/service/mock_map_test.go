package service_test

import (
	"database/sql"
	"errors"
	"naksha/db"
	"naksha/helper"
)

type MapAddNoLayerUserImpl struct {
	*DaoImpl
}

func (ma *MapAddNoLayerUserImpl) FetchData(query string, params []interface{}) db.SelectResult {
	rows := make([]map[string]string, 0)
	row := make(map[string]string)
	row["cnt"] = "0"
	rows = append(rows, row)
	columns := []string{"cnt"}

	return db.SelectResult{nil, columns, rows}
}

type MapAddImpl struct {
	*DaoImpl
}

func (ma *MapAddImpl) FetchData(query string, params []interface{}) db.SelectResult {
	rows := make([]map[string]string, 0)
	row := make(map[string]string)
	row["cnt"] = "1"
	rows = append(rows, row)
	columns := []string{"cnt"}

	return db.SelectResult{nil, columns, rows}
}

func (ma *MapAddImpl) TxTransaction(func(dao db.Dao) helper.Result) helper.Result {
	result := helper.MakeSuccessResult()
	result.AddToData("redir_url", "/maps/120/show")

	return result
}

type MapDetailsImpl struct {
	*DaoImpl
}

func (md *MapDetailsImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := make(map[string]string)
	row["id"] = "12"
	row["name"] = "Some Map"

	return row, nil
}

func (md *MapDetailsImpl) FetchData(query string, params []interface{}) db.SelectResult {
	rows := make([]map[string]string, 0)
	row := map[string]string{
		"table_name": "Table 1",
		"layer_id":   "12",
	}
	rows = append(rows, row)

	row = map[string]string{
		"table_name": "Table 2",
		"layer_id":   "23",
	}
	rows = append(rows, row)

	return db.SelectResult{nil, []string{"table_name", "layer_id"}, rows}
}

type MapUserMapsImpl struct {
	*DaoImpl
}

func (mu *MapUserMapsImpl) SelectSelect(sq db.SelectQuery) db.SelectResult {
	rows := make([]map[string]string, 0)
	row := map[string]string{"id": "1", "name": "Map 1"}
	rows = append(rows, row)
	row = map[string]string{"id": "2", "name": "Map 2"}
	rows = append(rows, row)
	row = map[string]string{"id": "3", "name": "Map 3"}
	rows = append(rows, row)
	row = map[string]string{"id": "4", "name": "Map 4"}
	rows = append(rows, row)

	return db.SelectResult{nil, []string{"id", "name"}, rows}
}

func (mu *MapUserMapsImpl) CountWhere(table string, where map[string]interface{}) (int, error) {
	return 4, nil
}

type MapShowMapImpl struct {
	*DaoImpl
}

func (mu *MapShowMapImpl) FindWhere(table string, where map[string]interface{}) (map[string]string, error) {
	row := map[string]string{"id": "102", "user_id": "1", "name": "First Map", "hash": "sdklf3234lkasfd-34", "base_layer": "o-osm"}
	return row, nil
}

func (mu *MapShowMapImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := map[string]string{"id": "10", "name": "SomeOne", "username": "someone", "password": "slkfja;lkweroiqwre", "table_seq": "2", "google_maps_key": "ioewrl23490", "bing_maps_key": "8lsdfl324asf"}
	return row, nil
}

func (mu *MapShowMapImpl) FetchData(query string, params []interface{}) db.SelectResult {
	var columns []string
	rows := make([]map[string]string, 0)
	if len(params) == 0 { // extents query
		row := map[string]string{"table_name": "table_1", "xtnt": "ST_BOUNDS(10 10)"}
		rows = append(rows, row)
		columns = []string{"table_name", "xtnt"}
	} else { // tables and layers query
		row := map[string]string{"name": "Table 1", "table_name": "table_1", "layer_hash": "lwersfdkjf", "update_hash": "wiolsflsdf", "infowindow": "{fields: [\"name\"]}", "layer_id": "123"}
		rows = append(rows, row)
		columns = []string{"name", "table_name", "layer_hash", "update_hash", "infowindow", "layer_id"}
	}

	return db.SelectResult{nil, columns, rows}
}

type MapDeleteImpl struct {
	*DaoImpl
}

func (mu *MapDeleteImpl) TxTransaction(func(dao db.Dao) helper.Result) helper.Result {
	result := helper.MakeSuccessResult()
	result.AddToData("redir_url", "/maps/index")

	return result
}

type MapBaseLayerUpdateImpl struct {
	*DaoImpl
}

func (mb *MapBaseLayerUpdateImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := map[string]string{"id": "1", "name": "SomeUser", "username": "someone", "password": "asdlkjadf", "table_seq": "10", "google_maps_key": "googl", "bing_maps_key": "bing"}

	return row, nil
}

type MapBaseLayerUpdateErrorImpl struct {
	*DaoImpl
}

func (mb *MapBaseLayerUpdateErrorImpl) Update(table string, values map[string]interface{}, where map[string]interface{}) (sql.Result, error) {
	return nil, errors.New("Some update error")
}

type MapSearchTablesImpl struct {
	*DaoImpl
}

func (mq *MapSearchTablesImpl) FetchData(query string, params []interface{}) db.SelectResult {
	rows := make([]map[string]string, 0)
	row := map[string]string{"value": "Table 1", "layer_id": "1"}
	rows = append(rows, row)
	row = map[string]string{"value": "Table 2", "layer_id": "39"}
	rows = append(rows, row)
	columns := []string{"value", "layer_id"}

	return db.SelectResult{nil, columns, rows}
}

type MapAddLayerImpl struct {
	*DaoImpl
}

func (ma *MapAddLayerImpl) FetchData(query string, params []interface{}) db.SelectResult {
	rows := make([]map[string]string, 0)
	row := map[string]string{"cnt": "1"}
	rows = append(rows, row)

	return db.SelectResult{nil, []string{"*"}, rows}
}

type MapDeleteLayerImpl struct {
	*DaoImpl
}

func (md *MapDeleteLayerImpl) CountWhere(table string, where map[string]interface{}) (int, error) {
	return 2, nil
}

type MapBelongsToUserImpl struct {
	*DaoImpl
}

func (mb *MapBelongsToUserImpl) CountWhere(table string, where map[string]interface{}) (int, error) {
	return 1, nil
}
