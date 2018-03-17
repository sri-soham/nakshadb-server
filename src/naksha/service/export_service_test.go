package service_test

import (
	"fmt"
	"naksha/service"
	"net/http"
	"net/url"
	"testing"
)

func TestExportGetStatusInQueue(t *testing.T) {
	map_dao := &ExportStatusInQueueImpl{&DaoImpl{}}
	service := service.MakeExportService(map_dao)
	result := service.GetStatus("132")
	if !result.IsSuccess() {
		t.Errorf("ExportStatusInQueue: could not fetch details")
	}
	_, ok := result.GetDataByKey("export_status")
	if !ok {
		t.Errorf("ExportStatusInQueue: export_status is not set")
	}
	_, ok = result.GetDataByKey("export_name")
	if !ok {
		t.Errorf("ExportStatusInQueue: export_name is not set")
	}
	remove := result.GetBoolData("remove_export_id")
	if remove {
		t.Errorf("ExportStatusInQueue: remove_export_id should be false")
	}
}

func TestExportGetStatusError(t *testing.T) {
	map_dao := &ExportStatusErrorImpl{&DaoImpl{}}
	service := service.MakeExportService(map_dao)
	result := service.GetStatus("132")
	if !result.IsSuccess() {
		t.Errorf("ExportStatusError: could not fetch details")
	}
	_, ok := result.GetDataByKey("export_status")
	if !ok {
		t.Errorf("ExportStatusError: export_status is not set")
	}
	_, ok = result.GetDataByKey("export_name")
	if !ok {
		t.Errorf("ExportStatusError: export_name is not set")
	}
	remove := result.GetBoolData("remove_export_id")
	if !remove {
		t.Errorf("ExportStatusError: remove_export_id should be true")
	}
}

func TestExportGetStatusSuccess(t *testing.T) {
	map_dao := &ExportStatusSuccessImpl{&DaoImpl{}}
	service := service.MakeExportService(map_dao)
	result := service.GetStatus("132")
	if !result.IsSuccess() {
		t.Errorf("ExportStatusSuccess: could not fetch details")
	}
	_, ok := result.GetDataByKey("export_status")
	if !ok {
		t.Errorf("ExportStatusSuccess: export_status is not set")
	}
	_, ok = result.GetDataByKey("export_name")
	if !ok {
		t.Errorf("ExportStatusInSuccess: export_name is not set")
	}
	url := result.GetStringData("download_url")
	expected_url := "/exports/132/download"
	if url != expected_url {
		t.Errorf("ExportStatusSuccess: download-url - expected(%v), found(%v)", expected_url, url)
	}
	remove := result.GetBoolData("remove_export_id")
	if !remove {
		t.Errorf("ExportStatusSuccess: remove_export_id should be true")
	}
}

func TestExportDownloadNoFile(t *testing.T) {
	map_dao := &ExportDownloadNoFileImpl{&DaoImpl{}}
	service := service.MakeExportService(map_dao)
	result := service.Download("32", "../../../test")
	if result.IsSuccess() {
		t.Errorf("ExportDownloadNoFile: success returned even when file does not exist")
	}
}

func TestExportDownloadFile(t *testing.T) {
	map_dao := &ExportDownloadFileImpl{&DaoImpl{}}
	service := service.MakeExportService(map_dao)
	exports_path := "../../../test"
	result := service.Download("32", exports_path)
	if !result.IsSuccess() {
		t.Errorf("ExportDownloadFile: error even when file exists")
	}
	expected_filepath := exports_path + "/abcd1234/somefile.zip"
	filepath := result.GetStringData("file_path")
	if expected_filepath != filepath {
		t.Errorf("ExportDownloadFile: filepath - expected(%v), found(%v)", expected_filepath, filepath)
	}

	expected_content_type := "application/octet-stream"
	content_type := result.GetStringData("content_type")
	if expected_content_type != content_type {
		t.Errorf("ExportDownloadFile: content-type - expected(%v), found(%v)", expected_content_type, content_type)
	}
	expected_filename := "somefile.zip"
	filename := result.GetStringData("filename")
	if expected_filename != filename {
		t.Errorf("ExportDownloadFile: filename - expected(%v), found(%v)", expected_filename, filename)
	}
}

func TestExportOfUser(t *testing.T) {
	map_dao := &ExportOfUserImpl{&DaoImpl{}}
	service := service.MakeExportService(map_dao)
	uri_params := make(map[string]string)
	values := url.Values{}
	values.Set("page", "1")
	request := makeRequest(http.MethodGet, uri_params, values, true)
	result := service.ExportsOfUser(request)
	tmp_exports, ok := result.GetDataByKey("exports")
	if !ok {
		t.Errorf("ExportOfUser: exports not returned")
	}

	_, ok = result.GetDataByKey("pagination_links")
	if !ok {
		t.Errorf("ExportOfUser: pagination-links not returned")
	}

	_, ok = result.GetDataByKey("pagination_text")
	if !ok {
		t.Errorf("ExportOfUser: pagination-text not returned")
	}

	exports := tmp_exports.([]map[string]interface{})
	if len(exports) != 4 {
		t.Errorf("ExportOfUser: exports count - expected(4), found(%d)", len(exports))
	}
	if exports[0]["name"] != "File1" {
		t.Errorf("ExportOfUser: name - expected(File1), found(%v)", exports[0]["name"])
	}
	if exports[0]["type"] != "ESRI Shapefile" {
		t.Errorf("ExportOfUser: type - expected(ESRI Shapefile), found(%v)", exports[0]["type"])
	}
	if exports[0]["status"] != "Success" {
		t.Errorf("ExportOfUser: status - expected(Success), found(%v)", exports[0]["status"])
	}
	download_link := "<a href=\"/exports/1/download\" target=\"_blank\">Download</a>"
	if fmt.Sprintf("%s", exports[0]["download_link"]) != download_link {
		t.Errorf("ExportOfUser: download-link - expected('%v'), found('%v')", download_link, exports[0]["download_link"])
	}

	if exports[1]["type"] != "CSV" {
		t.Errorf("ExportOfUser: type - expected(CSV), found(%v)", exports[1]["type"])
	}
	download_link = "[N/A]"
	if exports[1]["download_link"] != download_link {
		t.Errorf("ExportOfUser: download-link - expected(%v), found(%v)", download_link, exports[1]["download_link"])
	}

	if exports[2]["type"] != "GeoJSON" {
		t.Errorf("ExportOfUser: type - expected(GeoJSON), found(%v)", exports[2]["type"])
	}
	download_link = "[N/A]"
	if exports[2]["download_link"] != download_link {
		t.Errorf("ExportOfUser: download-link - expected(%v), found(%v)", download_link, exports[2]["download_link"])
	}

	if exports[3]["type"] != "KML" {
		t.Errorf("ExportOfUser: type - expected(KML), found(%v)", exports[3]["type"])
	}
}

func TestExportDeleteFail(t *testing.T) {
	map_dao := &ExportDeleteFailImpl{&DaoImpl{}}
	service := service.MakeExportService(map_dao)
	exports_dir := "../../../test"
	result := service.Delete(exports_dir, "1")
	if result.IsSuccess() {
		t.Errorf("ExportDeleteFail: successful even when delete failed")
	}
}

func TestExportDelete(t *testing.T) {
	map_dao := &ExportDeleteImpl{&DaoImpl{}}
	service := service.MakeExportService(map_dao)
	exports_dir := "../../../test"
	result := service.Delete(exports_dir, "1")
	if !result.IsSuccess() {
		t.Errorf("ExportDeleteFail: error when it should have been success")
	}
}

func TestExportBelongsToUserFail(t *testing.T) {
	map_dao := &ExportBelongsToUserImpl{&DaoImpl{}}
	service := service.MakeExportService(map_dao)
	belongs := service.ExportBelongsToUser("20", 1)
	if belongs {
		t.Errorf("ExportBelongsToUserFail: expected(false), found(%v)", belongs)
	}
}

func TestExportBelongsToUser(t *testing.T) {
	map_dao := &ExportBelongsToUserImpl{&DaoImpl{}}
	service := service.MakeExportService(map_dao)
	belongs := service.ExportBelongsToUser("10", 1)
	if !belongs {
		t.Errorf("ExportBelongsToUserFail: expected(true), found(%v)", belongs)
	}
}
