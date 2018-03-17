package importer_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"naksha"
	"naksha/db"
	"naksha/helper"
	"naksha/importer"
	"naksha/logger"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
)

var map_user_dao db.Dao
var app_config naksha.AppConfig

const TEST_DIR = "../../../test"

func init() {
	var err error

	config_path := flag.String("config_path", "", "Full path to config File")
	flag.Parse()
	if len(*config_path) == 0 {
		log.Println("Please give config file path")
		os.Exit(1)
	}

	app_config, err = naksha.MakeConfig(*config_path)
	if err != nil {
		log.Println("Could not parse config file. Error = ", err)
		os.Exit(1)
	}
	map_user_dao, err = db.MapUserDao(app_config)
	if err != nil {
		log.Println("Could not connect to database. Error = ", err)
		os.Exit(1)
	}

	err = logger.InitLoggers(TEST_DIR)
	if err != nil {
		log.Println("Could not open log files. Error = ", err)
		os.Exit(1)
	}
}

func TestEmptyImport(t *testing.T) {
	table_name, id, err := addTableForTest()
	if err != nil {
		t.Errorf("Could not create table. Error = %v", err)
	}

	importer := importer.MakeImporter(map_user_dao)
	result := importer.EmptyImport(table_name, id)
	if !result.IsSuccess() {
		t.Errorf("Empty import failed. Errors = %s", result.GetErrors())
	}
}

func TestHandleUpload(t *testing.T) {
	req, err := createRequestForTest(TEST_DIR+"/empty_two.csv", "empty_two.csv")
	if err != nil {
		t.Errorf("Error = %v", err)
	}
	importer := importer.MakeImporter(map_user_dao)
	uploaded_file := importer.HandleUpload(req, TEST_DIR+"/imports")
	if !uploaded_file.IsValid {
		t.Errorf("Handle upload failed. Error = %s", uploaded_file.ErrorMessage)
	}
}

func TestHandleUploadCsvZip(t *testing.T) {
	req, err := createRequestForTest(TEST_DIR+"/empty_two.csv.zip", "empty_two.csv.zip")
	if err != nil {
		t.Errorf("Error = %v", err)
	}
	my_importer := importer.MakeImporter(map_user_dao)
	uploaded_file := my_importer.HandleUpload(req, TEST_DIR+"/imports")
	if !uploaded_file.IsValid {
		t.Errorf("Handle upload failed. Error = %s", uploaded_file.ErrorMessage)
	}
	if uploaded_file.FileType != importer.CSV_FILE {
		t.Errorf("FileType is not correct. Expected(%d), found(%d)", importer.CSV_FILE, uploaded_file.FileType)
	}
}

func TestHandleUploadShapeZip(t *testing.T) {
	req, err := createRequestForTest(TEST_DIR+"/empty_two.zip", "empty_two.zip")
	if err != nil {
		t.Errorf("Error = %v", err)
	}
	my_importer := importer.MakeImporter(map_user_dao)
	uploaded_file := my_importer.HandleUpload(req, TEST_DIR+"/imports")
	if !uploaded_file.IsValid {
		t.Errorf("Handle upload failed. Error = %s", uploaded_file.ErrorMessage)
	}
	if uploaded_file.FileType != importer.SHAPE_FILE {
		t.Errorf("FileType is not correct. Expected(%d), found(%d)", importer.SHAPE_FILE, uploaded_file.FileType)
	}
}

func TestFileImportShapeFile(t *testing.T) {
	fileImportTest(t, "empty_two.shp", importer.SHAPE_FILE)
}

func TestFileImportGeojsonFile(t *testing.T) {
	fileImportTest(t, "empty_two.geojson", importer.GEOJSON_FILE)
}

func TestFileImportKmlFile(t *testing.T) {
	fileImportTest(t, "empty_two.kml", importer.KML_FILE)
}

func TestFileImportCsvFile(t *testing.T) {
	fileImportTest(t, "empty_two.csv", importer.CSV_FILE)
}

func addTableForTest() (string, int, error) {
	columns := make(map[string]interface{})
	columns["user_id"] = 1
	columns["name"] = "Some Random Name"
	table_name := "tbl_" + strings.ToLower(helper.RandomString(12))
	columns["table_name"] = table_name
	schema_name := "public"
	columns["schema_name"] = schema_name
	columns["status"] = importer.UPLOADED

	id, err := map_user_dao.Insert(db.MasterTable(), columns, "id")
	return schema_name + "." + table_name, id, err
}

func createRequestForTest(path string, filename string) (*http.Request, error) {
	var b bytes.Buffer
	var msg string

	w := multipart.NewWriter(&b)
	f, err := os.Open(path)
	if err != nil {
		msg = fmt.Sprintf("Could not open file for attaching. Error = %v", err)
		return nil, errors.New(msg)
	}
	defer f.Close()
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		msg = fmt.Sprintf("Could not create form file. Error = %v", err)
		return nil, errors.New(msg)
	}
	if _, err = io.Copy(fw, f); err != nil {
		msg = fmt.Sprintf("Could not copy file to buffer. Error = %v", err)
		return nil, errors.New(msg)
	}
	if fw, err = w.CreateFormField("name"); err != nil {
		msg = fmt.Sprintf("Could not create form field. Error = %v", err)
		return nil, errors.New(msg)
	}
	if _, err = fw.Write([]byte("Some Table")); err != nil {
		msg = fmt.Sprintf("Could not set value to form field 'name'. Error = %v", err)
		return nil, errors.New(msg)
	}
	w.Close()
	req, err := http.NewRequest(http.MethodPost, "/some/url", &b)
	if err != nil {
		msg = fmt.Sprintf("Could not create request. Error = %v", err)
		return nil, errors.New(msg)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	return req, nil
}

func fileImportTest(t *testing.T, filename string, file_type int) {
	// create a temporary directory in test/imports/ and copy the file to
	// that directory. Import process deletes the directory containing the
	// uploaded file.
	tmp_dir := TEST_DIR + "/imports/" + helper.RandomString(8)
	err := os.Mkdir(tmp_dir, os.ModePerm)
	if err != nil {
		t.Errorf("Could not create temporary imports directory. Error = %v", err)
	}
	err = copyFile(TEST_DIR, tmp_dir, filename)
	if err != nil {
		t.Errorf("Could not copy file %s. Error = %s", filename, err)
	}
	if path.Ext(filename) == ".shp" {
		tmp_file := strings.Replace(filename, ".shp", ".shx", 1)
		err = copyFile(TEST_DIR, tmp_dir, tmp_file)
		if err != nil {
			t.Errorf("Could not copy file %s. Error = %s", tmp_file, err)
		}

		tmp_file = strings.Replace(filename, ".shp", ".dbf", 1)
		err = copyFile(TEST_DIR, tmp_dir, tmp_file)
		if err != nil {
			t.Errorf("Could not copy file %s. Error = %s", tmp_file, err)
		}

		tmp_file = strings.Replace(filename, ".shp", ".prj", 1)
		err = copyFile(TEST_DIR, tmp_dir, tmp_file)
		if err != nil {
			t.Errorf("Could not copy file %s. Error = %s", tmp_file, err)
		}
	}

	table_name, id, err := addTableForTest()
	if err != nil {
		t.Errorf("Could not create table. Error = %v", err)
	}
	uploaded_file := importer.UploadedFile{
		IsValid:      true,
		FileType:     file_type,
		OutfilePath:  tmp_dir + "/" + filename,
		ErrorMessage: "",
	}
	my_importer := importer.MakeImporter(map_user_dao)
	my_importer.FileImport(&app_config, uploaded_file, table_name, id)
	row, err := map_user_dao.Find(db.MasterTable(), id)
	if err != nil {
		t.Errorf("Could not fetch table details")
	}
	expected_status := fmt.Sprintf("%d", importer.READY)
	if row["status"] != expected_status {
		t.Errorf("Import failed. Status: expected(%s), found(%s)", expected_status, row["status"])
	}
}

func copyFile(src_dir string, dst_dir string, filename string) error {
	data, err := ioutil.ReadFile(src_dir + "/" + filename)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst_dir+"/"+filename, data, 0644)

	return err
}
