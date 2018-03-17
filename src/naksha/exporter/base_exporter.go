package exporter

import (
	"errors"
	"fmt"
	"naksha"
	"naksha/db"
	"naksha/helper"
	"naksha/logger"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	ST_IN_QUEUE = 0
	ST_SUCCESS  = 10
	ST_ERROR    = -10

	SHAPE_FILE   = 10
	CSV_FILE     = 20
	GEOJSON_FILE = 30
	KML_FILE     = 40
)

type BaseExporter struct {
	map_user_dao db.Dao
	user_id      int
	table_id     string
	schema_name  string
	table_name   string
	id           int
	filename     string
	hash         string
}

func (be *BaseExporter) updateStatus(status int) {
	values := make(map[string]interface{})
	values["status"] = status
	values["updated_at"] = time.Now().Format(time.RFC3339)
	where := map[string]interface{}{"id": be.id}
	_, err := be.map_user_dao.Update(db.MasterExport(), values, where)
	if err != nil {
		msg := fmt.Sprintf("Could not update status to %v", status)
		be.logMessage(msg)
	}
}

func (be *BaseExporter) start(ext string) error {
	where := map[string]interface{}{"id": be.table_id}
	sres := be.map_user_dao.SelectWhere(db.MasterTable(), []string{"*"}, where)
	if sres.Error != nil {
		be.logMessage(sres.Error.Error())
		return sres.Error
	}
	filename := sres.Rows[0]["name"]
	filename = strings.TrimSpace(filename)
	filename = strings.Replace(filename, " ", "_", -1)
	filename = strings.ToLower(filename)
	be.filename = filename
	be.schema_name = sres.Rows[0]["schema_name"]
	be.table_name = sres.Rows[0]["table_name"]
	be.hash = fmt.Sprintf("%v%s%v", be.user_id, helper.RandomString(64), be.table_id)

	values := make(map[string]interface{})
	values["user_id"] = be.user_id
	values["table_id"] = be.table_id
	values["status"] = ST_IN_QUEUE
	values["filename"] = filename
	values["hash"] = be.hash
	values["extension"] = ext
	t := time.Now().Format(time.RFC3339)
	values["created_at"] = t
	values["updated_at"] = t
	id, err := be.map_user_dao.Insert(db.MasterExport(), values, "id")
	if err == nil {
		be.id = id
	} else {
		be.logMessage(err.Error())
	}

	return err
}

func (be *BaseExporter) createDirectory(exports_dir string) (string, error) {
	directory := fmt.Sprintf("%v/%v", exports_dir, be.hash)
	err := os.Mkdir(directory, 0775)
	if err != nil {
		msg := fmt.Sprintf("Could not create directory - %v; %s", directory, err)
		be.logMessage(msg)
	}

	return directory, err
}

func (be *BaseExporter) getColumns() ([]string, error) {
	columns := make([]string, 0)
	sres := be.map_user_dao.SelectWhere(
		"information_schema.columns",
		[]string{"column_name"},
		map[string]interface{}{"table_schema": be.schema_name, "table_name": be.table_name},
	)
	if sres.Error != nil {
		be.logMessage(sres.Error.Error())
		return columns, sres.Error
	}
	if len(sres.Rows) == 0 {
		be.logMessage("No such table")
		return columns, errors.New("No such table")
	}

	for _, row := range sres.Rows {
		if row["column_name"] != "the_geom_webmercator" {
			columns = append(columns, row["column_name"])
		}
	}

	return columns, nil
}

func (be *BaseExporter) schemaTable() string {
	return be.schema_name + "." + be.table_name
}

func (be *BaseExporter) queryFromColumns(columns []string) string {
	columns_str := strings.Join(columns, ", ")
	return fmt.Sprintf("SELECT %v FROM %v", columns_str, be.schemaTable())
}

func (be *BaseExporter) getFilePath(directory, filename, extension string) string {
	return directory + "/" + filename + extension
}

func (be *BaseExporter) toShapeFile(directory, extension, db_conn_str, query string) {
	file_path := be.getFilePath(directory, be.filename, extension)
	err := be.exportWithOgr2Ogr(file_path, db_conn_str, "ESRI Shapefile", query)
	if err != nil {
		return
	}

	shp_file := be.filename + ".shp"
	shx_file := be.filename + ".shx"
	prj_file := be.filename + ".prj"
	dbf_file := be.filename + ".dbf"
	zip_file := be.filename + ".zip"

	cmd := "zip"
	args := []string{zip_file, shp_file, shx_file, prj_file, dbf_file}
	cfg := exec.Command(cmd, args...)
	cfg.Dir = directory
	_, err = cfg.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("Could not create zip file %v. %s", directory+zip_file, err)
		be.logMessage(msg)
		be.updateStatus(ST_ERROR)
	} else {
		be.updateStatus(ST_SUCCESS)
	}
}

func (be *BaseExporter) toGeoJsonFile(directory, extension, db_conn_str, query string) {
	file_path := be.getFilePath(directory, be.filename, extension)
	err := be.exportWithOgr2Ogr(file_path, db_conn_str, "GeoJSON", query)
	if err == nil {
		be.updateStatus(ST_SUCCESS)
	}

}

func (be *BaseExporter) toKmlFile(directory, extension, db_conn_str, query string) {
	file_path := be.getFilePath(directory, be.filename, extension)
	err := be.exportWithOgr2Ogr(file_path, db_conn_str, "KML", query)
	if err == nil {
		be.updateStatus(ST_SUCCESS)
	}
}

func (be *BaseExporter) toCsvFile(directory, extension, db_conn_str, query string) {
	cmd := "psql"
	psql_out_path := "/tmp/" + be.filename + extension
	psql_cmd := fmt.Sprintf("\\COPY (%v) TO '%v' DELIMITER ',' CSV HEADER", query, psql_out_path)
	be.logMessage(psql_cmd)
	args := []string{db_conn_str, "-c", psql_cmd}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("output = %s. error = %s", string(out), err)
		be.logMessage(msg)
		be.updateStatus(ST_ERROR)
		return
	}
	file_path := be.getFilePath(directory, be.filename, extension)
	err = os.Rename(psql_out_path, file_path)
	if err != nil {
		msg := fmt.Sprintf("Could not move %v to %v", psql_out_path, file_path)
		be.logMessage(msg)
		be.updateStatus(ST_ERROR)
	} else {
		be.updateStatus(ST_SUCCESS)
	}
}

func (be *BaseExporter) exportWithOgr2Ogr(file_path string, db_conn_str string, format_name string, query string) error {
	cmd := "ogr2ogr"
	args := []string{"-f", format_name, file_path, db_conn_str, "-sql", query}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("%v. output = %v. error = %s", file_path, string(out), err)
		be.logMessage(msg)
		be.updateStatus(ST_ERROR)
	}

	return err
}

func (be *BaseExporter) logMessage(msg string) {
	err_msg := fmt.Sprintf("User: %v. Table: %v. Export: %v. %s", be.user_id, be.table_id, be.id, msg)
	logger.LogMessage(err_msg)
}

func GetAvailableFormats() map[int]string {
	export_formats := make(map[int]string)
	export_formats[SHAPE_FILE] = "ESRI Shape File"
	export_formats[CSV_FILE] = "CSV File"
	export_formats[GEOJSON_FILE] = "GeoJson File"
	export_formats[KML_FILE] = "KML File"

	return export_formats
}

func IsValidFormat(frmt string) bool {
	is_valid := false
	format, err := strconv.ParseInt(frmt, 10, 16)
	if err == nil {
		switch format {
		case SHAPE_FILE, CSV_FILE, GEOJSON_FILE, KML_FILE:
			is_valid = true
		}
	}

	return is_valid
}

func GetTypeFromExtension(extension string) string {
	var tstr string
	switch extension {
	case ".shp":
		tstr = "ESRI Shapefile"
	case ".csv":
		tstr = "CSV"
	case ".geojson":
		tstr = "GeoJSON"
	case ".kml":
		tstr = "KML"
	default:
		tstr = "Unknown"
	}

	return tstr
}

func ExportedFileDir(exports_dir, hash string) string {
	return fmt.Sprintf("%v/%v", exports_dir, hash)
}

func Export(app_config *naksha.AppConfig, map_user_dao db.Dao, user_id int, table_id string, frmt_str string) (int, error) {
	var extension string

	frmt, _ := strconv.Atoi(frmt_str)
	switch frmt {
	case SHAPE_FILE:
		extension = ".shp"
	case CSV_FILE:
		extension = ".csv"
	case GEOJSON_FILE:
		extension = ".geojson"
	case KML_FILE:
		extension = ".kml"
	default:
		logger.LogMessage("Invalid export file format: " + frmt_str)
		return 0, errors.New("Invalid export file format")
	}

	exporter := BaseExporter{map_user_dao, user_id, table_id, "", "", 0, "", ""}
	err := exporter.start(extension)
	if err != nil {
		return 0, err
	}

	directory, err := exporter.createDirectory(app_config.ExportsDir())
	if err != nil {
		return 0, err
	}

	columns, err := exporter.getColumns()
	if err != nil {
		return 0, err
	}
	query := exporter.queryFromColumns(columns)

	switch frmt {
	case SHAPE_FILE:
		pg_conn_string := db.Ogr2ogrConnectionString(*app_config)
		go exporter.toShapeFile(directory, extension, pg_conn_string, query)
	case CSV_FILE:
		pg_conn_string := db.PsqlConnectionString(*app_config)
		go exporter.toCsvFile(directory, extension, pg_conn_string, query)
	case GEOJSON_FILE:
		pg_conn_string := db.Ogr2ogrConnectionString(*app_config)
		go exporter.toGeoJsonFile(directory, extension, pg_conn_string, query)
	case KML_FILE:
		pg_conn_string := db.Ogr2ogrConnectionString(*app_config)
		go exporter.toKmlFile(directory, extension, pg_conn_string, query)
	}

	return exporter.id, nil
}
