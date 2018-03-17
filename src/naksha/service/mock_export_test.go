package service_test

import (
	"errors"
	"fmt"
	"naksha/db"
	"naksha/exporter"
)

func getExportStatus(status int) (map[string]string, error) {
	row := make(map[string]string)
	row["id"] = "132"
	row["status"] = fmt.Sprintf("%d", status)
	row["filename"] = "somefile"

	return row, nil
}

func getExportDownloadResult(filename string) db.SelectResult {
	rows := make([]map[string]string, 0)
	row := make(map[string]string)
	row["id"] = "32"
	row["hash"] = "abcd1234"
	row["filename"] = filename
	row["extension"] = ".shp"
	rows = append(rows, row)
	columns := []string{"id", "hash", "filename", "extension"}

	return db.SelectResult{nil, columns, rows}
}

func getExportDownload(filename string) (map[string]string, error) {
	row := make(map[string]string)
	row["id"] = "32"
	row["hash"] = "abcd1234"
	row["filename"] = filename
	row["extension"] = ".shp"

	return row, nil
}

type ExportStatusInQueueImpl struct {
	*DaoImpl
}

func (es *ExportStatusInQueueImpl) Find(table string, id interface{}) (map[string]string, error) {
	return getExportStatus(exporter.ST_IN_QUEUE)
}

type ExportStatusErrorImpl struct {
	*DaoImpl
}

func (es *ExportStatusErrorImpl) Find(table string, id interface{}) (map[string]string, error) {
	return getExportStatus(exporter.ST_ERROR)
}

type ExportStatusSuccessImpl struct {
	*DaoImpl
}

func (es *ExportStatusSuccessImpl) Find(table string, id interface{}) (map[string]string, error) {
	return getExportStatus(exporter.ST_SUCCESS)
}

type ExportDownloadNoFileImpl struct {
	*DaoImpl
}

func (ed *ExportDownloadNoFileImpl) Find(table string, id interface{}) (map[string]string, error) {
	return getExportDownload("otherfile")
}

type ExportDownloadFileImpl struct {
	*DaoImpl
}

func (ed *ExportDownloadFileImpl) Find(table string, id interface{}) (map[string]string, error) {
	return getExportDownload("somefile")
}

type ExportOfUserImpl struct {
	*DaoImpl
}

func (ed *ExportOfUserImpl) SelectSelect(sq db.SelectQuery) db.SelectResult {
	rows := make([]map[string]string, 0)
	row := map[string]string{
		"id":         "1",
		"status":     fmt.Sprintf("%d", exporter.ST_SUCCESS),
		"filename":   "file1",
		"extension":  ".shp",
		"created_at": "2017-09-10 10:10:10",
		"updated_at": "2017-09-10 10:11:01",
	}
	rows = append(rows, row)

	row = map[string]string{
		"id":         "2",
		"status":     fmt.Sprintf("%d", exporter.ST_ERROR),
		"filename":   "file2",
		"extension":  ".csv",
		"created_at": "2017-09-10 10:40:10",
		"updated_at": "2017-09-10 10:40:50",
	}
	rows = append(rows, row)

	row = map[string]string{
		"id":         "3",
		"status":     fmt.Sprintf("%d", exporter.ST_IN_QUEUE),
		"filename":   "file3",
		"extension":  ".geojson",
		"created_at": "2017-09-10 10:52:00",
		"updated_at": "2017-09-10 10:52:32",
	}
	rows = append(rows, row)

	row = map[string]string{
		"id":         "4",
		"status":     fmt.Sprintf("%d", exporter.ST_SUCCESS),
		"filename":   "file4",
		"extension":  ".kml",
		"created_at": "2017-09-10 11:01:12",
		"updated_at": "2017-09-10 11:01:51",
	}
	rows = append(rows, row)

	columns := []string{"id", "status", "filename", "extension", "created_at", "updated_at"}

	return db.SelectResult{nil, columns, rows}
}

func (ed *ExportOfUserImpl) CountWhere(table string, where map[string]interface{}) (int, error) {
	return 4, nil
}

type ExportDeleteFailImpl struct {
	*DaoImpl
}

func (ed *ExportDeleteFailImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := make(map[string]string)
	row["id"] = "1"
	row["hash"] = "werlqwe"

	return row, nil
}

func (ed *ExportDeleteFailImpl) Delete(table string, where map[string]interface{}) error {
	return errors.New("Some error")
}

type ExportDeleteImpl struct {
	*DaoImpl
}

func (ed *ExportDeleteImpl) Find(table string, id interface{}) (map[string]string, error) {
	row := make(map[string]string)
	row["id"] = "1"
	row["hash"] = "werlqwe"

	return row, nil
}

type ExportBelongsToUserImpl struct {
	*DaoImpl
}

func (ed *ExportBelongsToUserImpl) CountWhere(table string, where map[string]interface{}) (int, error) {
	t_id, _ := where["id"]
	id := t_id.(string)
	t_user_id := where["user_id"]
	user_id := t_user_id.(int)

	if user_id == 1 && id == "10" {
		return 1, nil
	} else {
		return 0, nil
	}
}
