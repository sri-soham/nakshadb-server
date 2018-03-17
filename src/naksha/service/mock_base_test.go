package service_test

import (
	"database/sql"
	"naksha/db"
	"naksha/helper"
)

type DaoImpl struct{}

func (i *DaoImpl) FetchData(query string, params []interface{}) db.SelectResult {
	return db.SelectResult{nil, make([]string, 0), make([]map[string]string, 0)}
}

func (i *DaoImpl) SelectAll(table string, columns []string) db.SelectResult {
	return db.SelectResult{nil, make([]string, 0), make([]map[string]string, 0)}
}

func (i *DaoImpl) SelectWhere(table string, columns []string, where map[string]interface{}) db.SelectResult {
	return db.SelectResult{nil, make([]string, 0), make([]map[string]string, 0)}
}

func (i *DaoImpl) SelectSelect(sq db.SelectQuery) db.SelectResult {
	return db.SelectResult{nil, make([]string, 0), make([]map[string]string, 0)}
}

func (i *DaoImpl) Find(table string, id interface{}) (map[string]string, error) {
	return make(map[string]string), nil
}

func (i *DaoImpl) FindWhere(table string, where map[string]interface{}) (map[string]string, error) {
	return make(map[string]string), nil
}

func (i *DaoImpl) SelectOne(table string, columns []string, where map[string]interface{}) (map[string]string, error) {
	return make(map[string]string), nil
}

func (i *DaoImpl) ModifyData(query string, params []interface{}) (sql.Result, error) {
	return nil, nil
}

func (i *DaoImpl) Insert(table string, values map[string]interface{}, auto_increment_column string) (int, error) {
	return 0, nil
}

func (i *DaoImpl) InsertWithGeometry(table string, values map[string]interface{}, auto_increment_column string) (int, error) {
	return 0, nil
}

func (i *DaoImpl) Update(table string, values map[string]interface{}, where map[string]interface{}) (sql.Result, error) {
	return nil, nil
}

func (i *DaoImpl) UpdateGeometry(table string, value string, where map[string]interface{}) (sql.Result, error) {
	return nil, nil
}

func (i *DaoImpl) Delete(table string, where map[string]interface{}) error {
	return nil
}

func (i *DaoImpl) CountAll(table string) (int, error) {
	return 0, nil
}

func (i *DaoImpl) CountWhere(table string, where map[string]interface{}) (int, error) {
	return 0, nil
}

func (i *DaoImpl) Exec(query string) (sql.Result, error) {
	return nil, nil
}

func (i *DaoImpl) StoredProcNoResult(query string, params []interface{}) error {
	return nil
}

func (i *DaoImpl) TxTransaction(tx_func func(db.Dao) helper.Result) helper.Result {
	return helper.MakeErrorResult()
}
