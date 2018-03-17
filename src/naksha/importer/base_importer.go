package importer

import (
	"fmt"
	"naksha"
	"naksha/db"
	"naksha/helper"
	"naksha/logger"
	"net/http"
	"os"
	"path/filepath"
)

const (
	UPLOADED = 10
	IMPORTED = 20
	UPDATED  = 30
	READY    = 40
	ERROR    = -10

	PRIMARY_KEY = "naksha_id"

	SHAPE_FILE   = 10
	CSV_FILE     = 20
	GEOJSON_FILE = 30
	KML_FILE     = 40
	INVALID_FILE = -10
)

type baseImporter struct {
	map_user_dao db.Dao
	table_name   string
	id           int
	outfile_path string
}

func (imp *baseImporter) updateStatus(status int) {
	values := make(map[string]interface{})
	where := make(map[string]interface{})
	values["status"] = status
	where["id"] = imp.id
	_, err := imp.map_user_dao.Update(db.MasterTable(), values, where)
	if err != nil {
		msg := fmt.Sprintf("%v: could not update status. id = %v, status = %v", imp.outfile_path, imp.id, status)
		logger.ImporterLog(msg)
	}
}

func (imp *baseImporter) logError(msg string, err error) {
	err_msg := fmt.Sprintf("%v : %s : %s: %s. error = %s", imp.id, imp.table_name, imp.outfile_path, msg, err)
	logger.ImporterLog(err_msg)
}

func (imp *baseImporter) logNotice(msg string) {
	notice := fmt.Sprintf("%v : %s : %s: %s", imp.id, imp.table_name, imp.outfile_path, msg)
	logger.ImporterLog(notice)
}

func (imp *baseImporter) getId() int {
	return imp.id
}

func (imp *baseImporter) getOutfilePath() string {
	return imp.outfile_path
}

func (imp *baseImporter) getTableName() string {
	return imp.table_name
}

func (imp *baseImporter) prepareTable() error {
	layer_hash := helper.HashForMapUrl()
	query := "SELECT public.naksha_prepare_table($1, $2, $3)"
	params := []interface{}{imp.id, imp.table_name, layer_hash}
	return imp.map_user_dao.StoredProcNoResult(query, params)
}

func (imp *baseImporter) deleteDirectory() error {
	ws_dir := filepath.Dir(imp.outfile_path)
	return os.RemoveAll(ws_dir)
}

type Importer interface {
	EmptyImport(table string, id int) helper.Result
	HandleUpload(req *http.Request, workspace_dir string) UploadedFile
	FileImport(app_config *naksha.AppConfig, uploaded_file UploadedFile, table_name string, id int)
}

type ImporterImpl struct {
	map_user_dao db.Dao
}

func (ii *ImporterImpl) EmptyImport(table string, id int) helper.Result {
	base_importer := baseImporter{ii.map_user_dao, table, id, ""}
	empty_importer := emptyImporter{base_importer}
	return empty_importer.doImport()
}

func (ii *ImporterImpl) HandleUpload(req *http.Request, workspace_dir string) UploadedFile {
	ifh := importFileHandler{req, workspace_dir, "", ""}
	return ifh.handle()
}

func (ii *ImporterImpl) FileImport(app_config *naksha.AppConfig, uploaded_file UploadedFile, table_name string, id int) {
	base_importer := baseImporter{ii.map_user_dao, table_name, id, uploaded_file.OutfilePath}
	switch uploaded_file.FileType {
	case SHAPE_FILE, GEOJSON_FILE, KML_FILE:
		conn_str := db.Ogr2ogrConnectionString(*app_config)
		ogr2ogr_importer := ogr2OgrImporter{base_importer, conn_str}
		ogr2ogr_importer.doImport()
	case CSV_FILE:
		conn_str := db.PsqlConnectionString(*app_config)
		csv_importer := csvFileImporter{base_importer, conn_str}
		csv_importer.doImport()
	}
}

func MakeImporter(dao db.Dao) ImporterImpl {
	return ImporterImpl{dao}
}
