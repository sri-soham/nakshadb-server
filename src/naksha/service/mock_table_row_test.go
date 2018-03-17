package service_test

import (
	"database/sql"
	"naksha/db"
)

func tableRowTableName() (map[string]string, error) {
	row := map[string]string{"table_name": "user_1_table_2"}

	return row, nil
}

type TableRowDataImpl struct {
	*DaoImpl
}

func (trd *TableRowDataImpl) Find(table string, id interface{}) (map[string]string, error) {
	return tableRowTableName()
}

func (trd *TableRowDataImpl) SelectSelect(sq db.SelectQuery) db.SelectResult {
	rows := make([]map[string]string, 0)
	row := map[string]string{
		"id":                   "1",
		"name":                 "Stop 1",
		"the_geom":             "SRID=4326;POINT(18 20)",
		"bus_numbers":          "10,20,30",
		"naksha_id":            "1",
		"the_geom_webmercator": "POINT(118, 120)",
	}
	rows = append(rows, row)

	row = map[string]string{
		"id":                   "2",
		"name":                 "Stop 2",
		"the_geom":             "SRID=4326;POINT(19 21)",
		"bus_numbers":          "11,21,31",
		"naksha_id":            "2",
		"the_geom_webmercator": "POINT(119, 121)",
	}
	rows = append(rows, row)

	columns := []string{"id", "name", "the_geom", "bus_numbers", "naksha_id", "the_geom_webmercator"}

	return db.SelectResult{nil, columns, rows}
}

func (trd *TableRowDataImpl) CountAll(table string) (int, error) {
	return 2, nil
}

type TableRowUpdateNonGeomImpl struct {
	*DaoImpl
	update_called bool
}

func (tu *TableRowUpdateNonGeomImpl) Find(table string, id interface{}) (map[string]string, error) {
	return tableRowTableName()
}

func (tu *TableRowUpdateNonGeomImpl) Update(table string, values map[string]interface{}, where map[string]interface{}) (sql.Result, error) {
	tu.update_called = true
	return nil, nil
}

func (tu *TableRowUpdateNonGeomImpl) IsUpdateCalled() bool {
	return tu.update_called
}

// if count is 1, then styles will be updated, else no.
// if geom_type us 'unknown', then styles will be updated, else no.
type TableRowUpdateGeomImpl struct {
	*DaoImpl
	geom_type              string
	update_geometry_called bool
	find_where_called      bool
	fetch_data_called      bool
}

func (tru *TableRowUpdateGeomImpl) Find(table string, id interface{}) (map[string]string, error) {
	return tableRowTableName()
}

func (tru *TableRowUpdateGeomImpl) UpdateGeometry(table string, value string, where map[string]interface{}) (sql.Result, error) {
	tru.update_geometry_called = true
	return nil, nil
}

func (tru *TableRowUpdateGeomImpl) FindWhere(table string, where map[string]interface{}) (map[string]string, error) {
	tru.find_where_called = true
	row := map[string]string{"geometry_type": tru.geom_type}
	return row, nil
}

func (tru *TableRowUpdateGeomImpl) FetchData(query string, params []interface{}) db.SelectResult {
	tru.fetch_data_called = true
	rows := make([]map[string]string, 0)
	row := map[string]string{"geom_type": "ST_MULTIPOLYGON"}
	rows = append(rows, row)

	return db.SelectResult{nil, []string{"geom_type"}, rows}
}

func (tru *TableRowUpdateGeomImpl) IsUpdateGeometryCalled() bool {
	return tru.update_geometry_called
}

func (tru *TableRowUpdateGeomImpl) IsFindWhereCalled() bool {
	return tru.find_where_called
}

func (tru *TableRowUpdateGeomImpl) IsFetchDataCalled() bool {
	return tru.fetch_data_called
}

type TableRowAddNonGeomImpl struct {
	*DaoImpl
	insert_called bool
}

func (tra *TableRowAddNonGeomImpl) Find(table string, id interface{}) (map[string]string, error) {
	return tableRowTableName()
}

func (tra *TableRowAddNonGeomImpl) Insert(table string, values map[string]interface{}, auto_inc_col string) (int, error) {
	tra.insert_called = true
	return 123, nil
}

func (tra *TableRowAddNonGeomImpl) IsInsertCalled() bool {
	return tra.insert_called
}

type TableRowAddGeomImpl struct {
	*DaoImpl
	geom_type              string
	insert_geometry_called bool
	find_where_called      bool
	fetch_data_called      bool
}

func (tra *TableRowAddGeomImpl) Find(table string, id interface{}) (map[string]string, error) {
	return tableRowTableName()
}

func (tra *TableRowAddGeomImpl) InsertWithGeometry(table string, values map[string]interface{}, auto_inc_col string) (int, error) {
	tra.insert_geometry_called = true
	return 4, nil
}

func (tra *TableRowAddGeomImpl) FindWhere(table string, where map[string]interface{}) (map[string]string, error) {
	tra.find_where_called = true
	row := map[string]string{"geometry_type": tra.geom_type}
	return row, nil
}

func (tra *TableRowAddGeomImpl) FetchData(query string, params []interface{}) db.SelectResult {
	tra.fetch_data_called = true
	rows := make([]map[string]string, 0)
	row := map[string]string{"geom_type": "ST_MULTILINESTRING"}
	rows = append(rows, row)

	return db.SelectResult{nil, []string{"geom_type"}, rows}
}

type TableRowDeleteImpl struct {
	*DaoImpl
}

func (trd *TableRowDeleteImpl) Find(table string, id interface{}) (map[string]string, error) {
	return tableRowTableName()
}

type TableRowShowImpl struct {
	*DaoImpl
}

func (trs *TableRowShowImpl) Find(table string, id interface{}) (map[string]string, error) {
	return tableRowTableName()
}

func (trs *TableRowShowImpl) FindWhere(table string, where map[string]interface{}) (map[string]string, error) {
	row := map[string]string{"infowindow": `{"fields":["name", "bus_numbers"]}`}
	return row, nil
}

func (trs *TableRowShowImpl) SelectWhere(table string, columns []string, where map[string]interface{}) db.SelectResult {
	rows := make([]map[string]string, 0)
	row := map[string]string{"name": "Abids", "bus_numbers": "8,5,7"}
	rows = append(rows, row)
	row = map[string]string{"name": "Liberty", "bus_numbers": "8,5,6,7"}
	rows = append(rows, row)
	return db.SelectResult{nil, []string{"name", "bus_numbers"}, rows}
}

type TableRowBelongsImpl struct {
	*DaoImpl
}

func (trb *TableRowBelongsImpl) CountWhere(table string, where map[string]interface{}) (int, error) {
	return 1, nil
}
