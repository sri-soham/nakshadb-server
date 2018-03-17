package service_test

import (
	"errors"
	"fmt"
	"naksha"
	"naksha/db"
	"naksha/helper"
	"naksha/importer"
	"net/http"
)

type MockImporter struct {
}

func (mi *MockImporter) EmptyImport(table_name string, id int) helper.Result {
	return helper.MakeSuccessResult()
}

func (mi *MockImporter) HandleUpload(req *http.Request, workspace_dir string) importer.UploadedFile {
	return importer.UploadedFile{
		IsValid:      true,
		FileType:     0,
		OutfilePath:  "",
		ErrorMessage: "",
	}
}

func (mi *MockImporter) FileImport(app_config *naksha.AppConfig, uploaded_file importer.UploadedFile, table_name string, id int) {
}

type TableEmptyTableDaoImpl struct {
	*DaoImpl
}

func (tet *TableEmptyTableDaoImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := map[string]string{
		"id":        "1",
		"name":      "Tester",
		"table_seq": "10",
	}

	return row, nil
}

func (tet *TableEmptyTableDaoImpl) Insert(table string, values map[string]interface{}, id string) (int, error) {
	return 11, nil
}

type TableCheckStatusDBErrorDaoImpl struct {
	*DaoImpl
}

func (tc *TableCheckStatusDBErrorDaoImpl) Find(table string, id interface{}) (map[string]string, error) {
	err := errors.New("Query Failed")
	return make(map[string]string), err
}

type TableCheckStatusNoRowsDaoImpl struct {
	*DaoImpl
}

func (tc *TableCheckStatusNoRowsDaoImpl) Find(table string, id interface{}) (map[string]string, error) {
	err := errors.New("No such record")
	return make(map[string]string), err
}

type TableCheckStatusDaoImpl struct {
	*DaoImpl
	status int
}

func (tc *TableCheckStatusDaoImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := map[string]string{"name": "Some Import", "status": fmt.Sprintf("%d", tc.status)}

	return row, nil
}

type TableGetDetailsDaoImpl struct {
	*DaoImpl
}

func (tg *TableGetDetailsDaoImpl) Find(table string, id interface{}) (map[string]string, error) {
	return map[string]string{"id": "2", "table_name": "tbl_sdfl1", "name": "Some Table"}, nil
}

func (tg *TableGetDetailsDaoImpl) FindWhere(table string, where map[string]interface{}) (map[string]string, error) {
	row := map[string]string{
		"id":            "3",
		"style":         "Some Style",
		"geometry_type": "ST_POLYGON",
		"hash":          "hash123",
		"infowindow":    "{fields:[\"name\"]}",
		"update_hash":   "hash123update",
	}

	return row, nil
}

func (tg *TableGetDetailsDaoImpl) SelectWhere(table string, columns []string, where map[string]interface{}) db.SelectResult {
	rows := make([]map[string]string, 0)

	rows = append(rows, map[string]string{"column_name": "naksha_id"})
	rows = append(rows, map[string]string{"column_name": "name"})
	rows = append(rows, map[string]string{"column_name": "the_geom"})
	rows = append(rows, map[string]string{"column_name": "stop_count"})
	rows = append(rows, map[string]string{"column_name": "the_geom_webmercator"})
	rows = append(rows, map[string]string{"column_name": "created_at"})
	rows = append(rows, map[string]string{"column_name": "updated_at"})

	return db.SelectResult{nil, columns, rows}
}

func (tg *TableGetDetailsDaoImpl) SelectAll(table string, columns []string) db.SelectResult {
	rows := make([]map[string]string, 0)
	row := map[string]string{"xtnt": "alkdjfasldfs"}
	rows = append(rows, row)

	return db.SelectResult{nil, columns, rows}
}

type TableUpdateStylesDaoImpl struct {
	*DaoImpl
}

func (tus *TableUpdateStylesDaoImpl) FindWhere(table string, where map[string]interface{}) (map[string]string, error) {
	row := map[string]string{"geometry_type": "linestring"}

	return row, nil
}

type TableDeleteColumnDaoImpl struct {
	*DaoImpl
}

func (td *TableDeleteColumnDaoImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := map[string]string{"id": "10", "table_name": "tbl_ab12"}
	return row, nil
}

func (td *TableDeleteColumnDaoImpl) TxTransaction(tx_func func(db.Dao) helper.Result) helper.Result {
	return helper.MakeSuccessResult()
}

type TableAddColumnDaoImpl struct {
	*DaoImpl
}

func (ta *TableAddColumnDaoImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := map[string]string{"id": "10", "table_name": "tbl_ab12"}
	return row, nil
}

type TableBelongsToUserDaoImpl struct {
	*DaoImpl
	count int
}

func (tb *TableBelongsToUserDaoImpl) CountWhere(table string, where map[string]interface{}) (int, error) {
	return tb.count, nil
}
