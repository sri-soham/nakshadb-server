package service

import (
	"fmt"
	"html/template"
	"naksha"
	"naksha/db"
	"naksha/exporter"
	"naksha/helper"
	"naksha/logger"
	"os"
	"path"
	"strconv"
	"strings"
)

type ExportService interface {
	GetStatus(id string) helper.Result
	Download(export_id string, exports_dir string) helper.Result
	ExportsOfUser(req *naksha.Request) helper.Result
	Delete(exports_dir, id string) helper.Result
	ExportBelongsToUser(export_id string, user_id int) bool
}

type ExportServiceImpl struct {
	map_user_dao db.Dao
}

func (es *ExportServiceImpl) GetStatus(id string) helper.Result {
	var result helper.Result

	export, err := es.map_user_dao.Find(db.MasterExport(), id)
	if err != nil {
		return helper.ErrorResultError(err)
	}

	result = helper.MakeSuccessResult()
	export_status, _ := strconv.Atoi(export["status"])
	result.AddToData("export_status", export_status)
	result.AddToData("export_name", export["filename"])
	if export_status == exporter.ST_SUCCESS {
		result.AddToData("download_url", fmt.Sprintf("/exports/%v/download", id))
	}
	if export_status == exporter.ST_IN_QUEUE {
		result.AddToData("remove_export_id", false)
	} else {
		result.AddToData("remove_export_id", true)
	}

	return result
}

func (es *ExportServiceImpl) Download(export_id string, exports_dir string) helper.Result {
	var result helper.Result

	export, err := es.map_user_dao.Find(db.MasterExport(), export_id)
	if err != nil {
		return helper.ErrorResultError(err)
	}
	file_path := exports_dir + "/" + export["hash"] + "/" + export["filename"]
	if export["extension"] == ".shp" {
		file_path += ".zip"
	} else {
		file_path += export["extension"]
	}
	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		return helper.ErrorResultString("File not available")
	}

	result = helper.MakeSuccessResult()
	result.AddToData("file_path", file_path)
	result.AddToData("content_type", "application/octet-stream")
	result.AddToData("filename", path.Base(file_path))

	return result
}

func (es *ExportServiceImpl) ExportsOfUser(req *naksha.Request) helper.Result {
	var result helper.Result
	var page int

	user_id := helper.UserIDFromSession(req.Session)
	page_str := req.Request.FormValue("page")
	if len(page_str) == 0 {
		page_str = "1"
	}
	page, err := strconv.Atoi(page_str)
	if err != nil {
		page = 1
	}
	per_page := 20
	offset := (page - 1) * per_page
	where := map[string]interface{}{"user_id": user_id}
	sq := db.SelectQuery{
		db.MasterExport(),
		[]string{"*"},
		where,
		"filename",
		per_page,
		offset,
	}
	sres := es.map_user_dao.SelectSelect(sq)
	if sres.Error != nil {
		return helper.ErrorResultString("Query failed 1")
	}
	count, err := es.map_user_dao.CountWhere(db.MasterExport(), where)
	if err != nil {
		return helper.ErrorResultString("Query failed")
	}

	exports := make([]map[string]interface{}, 0)
	for _, row := range sres.Rows {
		export := make(map[string]interface{})
		export["name"] = es.nameFromFilename(row["filename"])
		export["type"] = exporter.GetTypeFromExtension(row["extension"])
		status, _ := strconv.Atoi(row["status"])
		switch status {
		case exporter.ST_IN_QUEUE:
			export["status"] = "In Queue"
		case exporter.ST_SUCCESS:
			export["status"] = "Success"
		case exporter.ST_ERROR:
			export["status"] = "Error"
		default:
			export["status"] = "Unknown"
		}
		if status == exporter.ST_SUCCESS {
			link := "<a href=\"/exports/" + row["id"] + "/download\" target=\"_blank\">Download</a>"
			export["download_link"] = template.HTML(link)
		} else {
			export["download_link"] = "[N/A]"
		}
		export["created_at"] = row["created_at"]
		export["updated_at"] = row["updated_at"]
		export["id"] = row["id"]

		exports = append(exports, export)
	}

	paginator := helper.MakePaginator(page, per_page, count, "/exports/index?page={page}")
	result = helper.MakeSuccessResult()
	result.AddToData("exports", exports)
	result.AddToData("pagination_links", paginator.Links())
	result.AddToData("pagination_text", paginator.Text())

	return result
}

func (es *ExportServiceImpl) Delete(exports_dir, id string) helper.Result {
	var result helper.Result

	export, err := es.map_user_dao.Find(db.MasterExport(), id)
	if err != nil {
		return helper.ErrorResultError(err)
	}

	full_path := exporter.ExportedFileDir(exports_dir, export["hash"])
	// if full_path does not exists, os.RemoveAll will return nil
	err = os.RemoveAll(full_path)
	if err != nil {
		logger.LogMessage("Could not delete directory: " + full_path)
		result = helper.ErrorResultString("Error 2")
		return result
	}
	where := map[string]interface{}{"id": id}
	err = es.map_user_dao.Delete(db.MasterExport(), where)
	if err == nil {
		result = helper.MakeSuccessResult()
	} else {
		result = helper.ErrorResultString("Error 3")
	}

	return result
}

func (es *ExportServiceImpl) ExportBelongsToUser(export_id string, user_id int) bool {
	var belongs bool

	where := make(map[string]interface{})
	where["id"] = export_id
	where["user_id"] = user_id
	count, err := es.map_user_dao.CountWhere(db.MasterExport(), where)
	if err == nil {
		belongs = (count == 1)
	} else {
		belongs = false
	}

	return belongs
}

func (es *ExportServiceImpl) nameFromFilename(filename string) string {
	filename = strings.TrimSpace(filename)
	tmp := strings.Split(filename, "_")
	parts := make([]string, 0)
	for _, t := range tmp {
		p := fmt.Sprintf("%s%v", strings.ToUpper(t[0:1]), t[1:])
		parts = append(parts, p)
	}

	return strings.Join(parts, " ")
}

func MakeExportService(map_user_dao db.Dao) ExportService {
	return &ExportServiceImpl{map_user_dao}
}
