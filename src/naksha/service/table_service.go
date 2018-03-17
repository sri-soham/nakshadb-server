package service

import (
	"fmt"
	"naksha"
	"naksha/db"
	"naksha/exporter"
	"naksha/helper"
	"naksha/importer"
	"naksha/logger"
	"naksha/validator"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TableService interface {
	UserTables(req *naksha.Request) helper.Result
	CreateTable(req *naksha.Request, app_config *naksha.AppConfig) helper.Result
	CheckStatus(id string) map[string]interface{}
	GetDetails(req *naksha.Request, tiler_url string) helper.Result
	UpdateStyles(req *naksha.Request) helper.Result
	DeleteTable(table_id string) helper.Result
	AddColumn(req *naksha.Request) helper.Result
	DeleteColumn(req *naksha.Request) helper.Result
	Infowindow(r *naksha.Request) helper.Result
	Export(req *naksha.Request, app_config *naksha.AppConfig) helper.Result
	ApiAccess(req *naksha.Request, db_api_user string) helper.Result
	TableBelongsToUser(table_id string, user_id int) bool
}

type TableServiceImpl struct {
	map_user_dao db.Dao
	my_importer  importer.Importer
}

func (ts *TableServiceImpl) UserTables(req *naksha.Request) helper.Result {
	var result helper.Result

	tmp := req.Request.FormValue("page")
	if len(tmp) == 0 {
		tmp = "1"
	}
	page, err := strconv.Atoi(tmp)
	if err != nil {
		page = 1
	}
	per_page := 20
	user_id := helper.UserIDFromSession(req.Session)
	sel := db.GetSelect()
	sel.Columns("*")
	sel.From(db.MasterTable())
	sel.Where("user_id", user_id)
	sel.OrderBy("name")
	sel.Limit(per_page)
	sel.Offset((page - 1) * per_page)
	sres := ts.map_user_dao.FetchData(sel.GetQuery(), sel.GetParams())
	if sres.Error != nil {
		result = helper.MakeFatalResult(sres.Error)
	} else {
		where := map[string]interface{}{"user_id": user_id}
		cnt, err := ts.map_user_dao.CountWhere(db.MasterTable(), where)
		if err != nil {
			result = helper.ErrorResultString("Query failed")
		} else {
			paginator := helper.MakePaginator(page, per_page, cnt, "/dashboard?page={page}")
			result = helper.MakeSuccessResult()
			result.AddToData("tables", sres.Rows)
			result.AddToData("pagination_links", paginator.Links())
			result.AddToData("pagination_text", paginator.Text())
		}
	}

	return result
}

func (ts *TableServiceImpl) CreateTable(req *naksha.Request, app_config *naksha.AppConfig) helper.Result {
	is_empty_table := req.Request.URL.Query().Get("is_empty")
	logger.LogMessage("is-empty-table =  " + is_empty_table)
	if is_empty_table == "1" {
		return ts.handleCreateEmptyTable(req)
	} else {
		return ts.handleUploadFile(req, app_config)
	}
}

func (ts *TableServiceImpl) handleCreateEmptyTable(req *naksha.Request) helper.Result {
	name := req.Request.PostFormValue("name")
	if len(name) == 0 {
		return helper.ErrorResultString("Please enter the name")
	}

	user_id := helper.UserIDFromSession(req.Session)
	table_name, err := ts.tableNameForImport(user_id, name)
	if err != nil {
		return helper.ErrorResultString("Could not generate table name")
	}
	schema_name := helper.SchemaNameFromSession(req.Session)
	id, err := ts.addTable(name, schema_name, table_name, user_id)
	if err != nil {
		return helper.ErrorResultString("Could not add record to table")
	}

	res := ts.my_importer.EmptyImport(schema_name+"."+table_name, id)
	if res.IsSuccess() {
		res.AddToData("url", fmt.Sprintf("/tables/%v/show", id))
	}

	return res
}

func (ts *TableServiceImpl) handleUploadFile(req *naksha.Request, app_config *naksha.AppConfig) helper.Result {
	var result helper.Result

	name := req.Request.PostFormValue("name")
	if len(name) == 0 {
		return helper.ErrorResultString("Name is required")
	}

	uploaded_file := ts.my_importer.HandleUpload(req.Request, app_config.ImportsDir())
	if !uploaded_file.IsValid {
		return helper.ErrorResultString(uploaded_file.ErrorMessage)
	}

	user_id := helper.UserIDFromSession(req.Session)
	table_name, err := ts.tableNameForImport(user_id, name)
	if err != nil {
		return helper.ErrorResultString("Could not generate table name")
	}
	schema_name := helper.SchemaNameFromSession(req.Session)
	id, err := ts.addTable(name, schema_name, table_name, user_id)
	if err != nil {
		return helper.ErrorResultString("Could not add record to table")
	}

	go ts.my_importer.FileImport(app_config, uploaded_file, schema_name+"."+table_name, id)

	result = helper.MakeSuccessResult()
	result.AddToData("id", id)

	return result
}

func (ts *TableServiceImpl) CheckStatus(id string) map[string]interface{} {
	var result = make(map[string]interface{})

	details, err := ts.map_user_dao.Find(db.MasterTable(), id)
	if err != nil {
		result["status"] = "error"
		result["errors"] = err.Error()
	} else {
		result["status"] = "success"
		result["import_name"] = details["name"]
		status, _ := strconv.Atoi(details["status"])
		switch status {
		case importer.READY:
			result["import_status"] = "success"
			result["table_url"] = "/tables/" + id + "/show"
			result["remove_import_id"] = "1"
		case importer.ERROR:
			result["import_status"] = "error"
			result["errors"] = "Import failed"
			result["remove_import_id"] = "1"
		default:
			result["import_status"] = "importing"
		}
	}

	return result
}

func (ts *TableServiceImpl) GetDetails(req *naksha.Request, tiler_url string) helper.Result {
	var result helper.Result

	id := req.UriParams["id"]
	table_details, err := ts.map_user_dao.Find(db.MasterTable(), id)
	if err != nil {
		return helper.ErrorResultError(err)
	}

	user_id := helper.UserIDFromSession(req.Session)
	user_details, err := ts.map_user_dao.Find(db.MasterUser(), user_id)
	if err != nil {
		return helper.ErrorResultString("No such user")
	}

	col_where := make(map[string]interface{})
	col_where["table_name"] = table_details["table_name"]
	col_where["table_schema"] = table_details["schema_name"]
	tmp := ts.map_user_dao.SelectWhere("information_schema.columns", []string{"column_name"}, col_where)
	if tmp.Error != nil {
		return helper.ErrorResultString("Query failed 2")
	}
	columns := make([]string, 0)
	for _, row := range tmp.Rows {
		if !(row["column_name"] == "the_geom_webmercator" || row["column_name"] == "the_geom") {
			columns = append(columns, row["column_name"])
		}
	}

	tmp = ts.map_user_dao.SelectAll(helper.SchemaTableFromDetails(table_details), []string{"ST_Extent(the_geom) AS xtnt"})
	if tmp.Error != nil {
		return helper.ErrorResultString("Query failed 3")
	}
	extent := tmp.Rows[0]["xtnt"]

	lyr_where := make(map[string]interface{})
	lyr_where["table_id"] = req.UriParams["id"]
	lyr_details, err := ts.map_user_dao.FindWhere(db.MasterLayer(), lyr_where)
	if err != nil {
		return helper.ErrorResultError(err)
	}
	layer_id := lyr_details["id"]
	style := lyr_details["style"]
	geometry_type := lyr_details["geometry_type"]
	hash := lyr_details["hash"]
	infowindow := lyr_details["infowindow"]
	update_hash := lyr_details["update_hash"]

	export_formats := exporter.GetAvailableFormats()
	result = helper.MakeSuccessResult()
	result.AddToData("table_details", table_details)
	result.AddToData("columns", strings.Join(columns, ","))
	result.AddToData("url", "/table_rows/"+table_details["id"]+"/")
	result.AddToData("map_url", tiler_url+"lyr/"+hash+"-[ts]/{z}/{x}/{y}.png")
	result.AddToData("extent", extent)
	result.AddToData("layer_id", layer_id)
	result.AddToData("geometry_type", geometry_type)
	result.AddToData("style", style)
	result.AddToData("infowindow", infowindow)
	result.AddToData("update_hash", update_hash)
	result.AddToData("export_formats", export_formats)
	result.AddToData("tables_url", "/tables/"+req.UriParams["id"])
	result.AddToData("base_layers", helper.BaseLayers())
	result.AddToData("base_layer", helper.BASE_MAP_OSM)
	result.AddToData("user_details", user_details)

	return result
}

func (ts *TableServiceImpl) UpdateStyles(req *naksha.Request) helper.Result {
	var result helper.Result
	var style string

	lyr_where := make(map[string]interface{})
	lyr_where["table_id"] = req.UriParams["id"]
	lyr_details, err := ts.map_user_dao.FindWhere(db.MasterLayer(), lyr_where)
	if err != nil {
		return helper.ErrorResultError(err)
	}
	geometry_type := lyr_details["geometry_type"]

	validator := validator.MakeStylesValidator(req.Request, geometry_type)
	errors := validator.Validate()
	style_generator := helper.MakeStyleGenerator(req.Request)
	if len(errors) == 0 {
		switch geometry_type {
		case "polygon":
			style = style_generator.PolygonStyle()
		case "linestring":
			style = style_generator.LineStringStyle()
		case "point":
			style = style_generator.PointStyle()
		}
		values := make(map[string]interface{})
		values["style"] = style
		values["update_hash"] = helper.HashForMapUrl()
		values["updated_at"] = time.Now().Format(time.RFC3339)
		where := make(map[string]interface{})
		where["table_id"] = req.UriParams["id"]
		_, err := ts.map_user_dao.Update(db.MasterLayer(), values, where)
		if err == nil {
			result = helper.MakeSuccessResult()
			result.AddToData("update_hash", values["update_hash"])
		} else {
			result = helper.MakeErrorResult()
			result.SetErrors([]string{"Query failed"})
		}
	} else {
		result = helper.MakeErrorResult()
		result.SetErrors(errors)
	}

	return result
}

func (ts *TableServiceImpl) DeleteTable(table_id string) helper.Result {
	var result helper.Result

	table_details, err := ts.map_user_dao.Find(db.MasterTable(), table_id)
	if err != nil {
		return helper.ErrorResultError(err)
	}

	where := make(map[string]interface{})
	where["id"] = table_id
	table_name := helper.SchemaTableFromDetails(table_details)
	err = ts.map_user_dao.Delete(db.MasterTable(), where)
	if err != nil {
		return helper.ErrorResultString("Could not delete table 1")
	}

	query := "DROP TABLE IF EXISTS " + table_name
	_, err = ts.map_user_dao.Exec(query)
	if err != nil {
		result = helper.ErrorResultString("Could not delete table 2")
	} else {
		result = helper.MakeSuccessResult()
	}

	return result
}

func (ts *TableServiceImpl) AddColumn(req *naksha.Request) helper.Result {
	var result helper.Result

	validator := validator.MakeAddColumnValidator(req.Request)
	errors := validator.Validate()
	if len(errors) > 0 {
		result = helper.MakeErrorResult()
		result.SetErrors(errors)
		return result
	}
	table_details, err := ts.map_user_dao.Find(db.MasterTable(), req.UriParams["id"])
	if err != nil {
		return helper.ErrorResultError(err)
	}

	column_name := req.Request.PostFormValue("name")
	column_where := make(map[string]interface{})
	column_where["table_name"] = table_details["table_name"]
	column_where["column_name"] = column_name
	column_where["schema_name"] = table_details["schema_name"]
	cnt, err := ts.map_user_dao.CountWhere("information_schema.columns", column_where)
	if err != nil {
		return helper.ErrorResultString("Query failed")
	}
	if cnt > 0 {
		return helper.ErrorResultString("Column with given name already exists")
	}
	table_name := helper.SchemaTableFromDetails(table_details)
	data_type_option := req.Request.PostFormValue("data_type")
	data_type := ""
	switch data_type_option {
	case "1":
		data_type = "INTEGER"
	case "2":
		data_type = "DOUBLE PRECISION"
	case "3":
		data_type = "VARCHAR"
	case "4":
		data_type = "TIMESTAMP WITHOUT TIME ZONE"
	default:
		data_type = "VARCHAR"
	}
	query := "ALTER TABLE " + table_name + " ADD COLUMN " + column_name + " " + data_type + " "
	_, err = ts.map_user_dao.Exec(query)
	if err != nil {
		result = helper.ErrorResultString("Could not add column")
	} else {
		result = helper.MakeSuccessResult()
	}

	return result
}

func (ts *TableServiceImpl) DeleteColumn(req *naksha.Request) helper.Result {
	var result helper.Result

	column_name := req.Request.PostFormValue("column_name")
	if len(column_name) == 0 {
		return helper.ErrorResultString("Column name not provided")
	}
	if column_name == "naksha_id" || column_name == "the_geom" ||
		column_name == "created_at" || column_name == "updated_at" ||
		column_name == "the_geom_webmercator" {
		return helper.ErrorResultString("This column cannot be deleted")
	}

	table_details, err := ts.map_user_dao.Find(db.MasterTable(), req.UriParams["id"])
	if err != nil {
		helper.ErrorResultString("Could not fetch details 1")
	}

	where := map[string]interface{}{"table_id": table_details["id"]}
	layer_details, err := ts.map_user_dao.FindWhere(db.MasterLayer(), where)
	if err != nil {
		helper.ErrorResultString("Could not fetch details 2")
	}

	result = ts.map_user_dao.TxTransaction(func(txdao db.Dao) helper.Result {
		var trans_result helper.Result

		infowindow, err := helper.StringToInfowindow(layer_details["infowindow"])
		if err != nil {
			return helper.ErrorResultString("Could not decode infowindow")
		}
		iw_fields := make([]string, 0)
		for _, v := range infowindow.Fields {
			if v != column_name {
				iw_fields = append(iw_fields, column_name)
			}
		}
		infowindow.Fields = iw_fields

		table_name := helper.SchemaTableFromDetails(table_details)
		query := "ALTER TABLE " + table_name + " DROP COLUMN IF EXISTS \"" + column_name + "\" "
		_, err = txdao.Exec(query)
		if err != nil {
			trans_result = helper.ErrorResultString("Could not delete column")
		} else {
			values := map[string]interface{}{"infowindow": helper.InfowindowToString(infowindow)}
			_, err = txdao.Update(db.MasterLayer(), values, where)
			if err != nil {
				trans_result = helper.ErrorResultString("Could not update infowindow fields")
			} else {
				trans_result = helper.MakeSuccessResult()
			}
		}

		return trans_result
	})

	return result
}

func (ts *TableServiceImpl) Infowindow(r *naksha.Request) helper.Result {
	var result helper.Result

	r.Request.ParseForm()
	columns, ok := r.Request.PostForm["columns"]
	if ok {
		valid_count := 0
		column_count := len(columns)
		for _, col := range columns {
			// column name should *not* be qualfied with db or schema name
			parts := strings.Split(col, ".")
			if len(parts) == 1 {
				if col == "the_geom" || col == "the_geom_webmercator" {
					// nothing to do
				} else {
					valid_count++
				}
			}
		}
		if valid_count != column_count {
			return helper.ErrorResultString("Invalid column names")
		}
	} else {
		columns = make([]string, 0)
	}

	values := make(map[string]interface{})
	ifw := helper.Infowindow{columns}
	values["infowindow"] = helper.InfowindowToString(ifw)
	values["updated_at"] = time.Now().Format(time.RFC3339)
	where := make(map[string]interface{})
	where["table_id"] = r.UriParams["id"]

	_, err := ts.map_user_dao.Update(db.MasterLayer(), values, where)
	if err == nil {
		result = helper.MakeSuccessResult()
	} else {
		result = helper.ErrorResultString("Query failed")
	}

	return result
}

func (ts *TableServiceImpl) Export(req *naksha.Request, app_config *naksha.AppConfig) helper.Result {
	var result helper.Result

	frmt_str := req.PostFormValue("format")
	if !exporter.IsValidFormat(frmt_str) {
		return helper.ErrorResultString("Invalid format")
	}

	user_id := helper.UserIDFromSession(req.Session)
	table_id := req.UriParams["id"]
	id, err := exporter.Export(app_config, ts.map_user_dao, user_id, table_id, frmt_str)
	if err != nil {
		result = helper.ErrorResultString("Export failed")
	} else {
		result = helper.MakeSuccessResult()
		result.AddToData("id", id)
	}

	return result
}

func (ts *TableServiceImpl) ApiAccess(req *naksha.Request, db_api_user string) helper.Result {
	var result helper.Result
	var err error
	var query string

	api_access := req.PostFormValue("api_access")
	if len(api_access) == 0 {
		return helper.ErrorResultString("Enter value for api access")
	}
	if !(api_access == "f" || api_access == "t") {
		return helper.ErrorResultString("Invalid value for api access")
	}

	table_details, err := ts.map_user_dao.Find(db.MasterTable(), req.UriParams["id"])
	if err != nil {
		return helper.ErrorResultError(err)
	}

	values := map[string]interface{}{"api_access": api_access}
	where := map[string]interface{}{"id": req.UriParams["id"]}
	_, err = ts.map_user_dao.Update(db.MasterTable(), values, where)
	if err != nil {
		return helper.ErrorResultString("Could not update 2")
	}

	if api_access == "t" {
		query = "GRANT SELECT ON " + helper.SchemaTableFromDetails(table_details) + " TO " + db_api_user
	} else {
		query = "REVOKE SELECT ON " + helper.SchemaTableFromDetails(table_details) + " FROM " + db_api_user
	}
	_, err = ts.map_user_dao.Exec(query)
	if err != nil {
		return helper.ErrorResultString("Could not update 2")
	}

	result = helper.MakeSuccessResult()

	return result
}

func (ts *TableServiceImpl) TableBelongsToUser(table_id string, user_id int) bool {
	var belongs bool

	where := make(map[string]interface{})
	where["id"] = table_id
	where["user_id"] = user_id
	count, err := ts.map_user_dao.CountWhere(db.MasterTable(), where)
	if err == nil {
		belongs = (count == 1)
	} else {
		belongs = false
	}

	return belongs
}

func (ts *TableServiceImpl) tableNameForImport(user_id int, name string) (string, error) {
	name = strings.ToLower(name)
	name = strings.TrimSpace(name)
	re1 := regexp.MustCompile("[^[:alnum:]]")
	re2 := regexp.MustCompile("_+")
	new_name := re1.ReplaceAllString(name, "_")
	new_name = re2.ReplaceAllString(new_name, "_")
	new_name = strings.Trim(new_name, "_")
	// Maximum allowed length of a table name in postgres is 63.
	// We truncate the name at 59 characters and check for presence of tables(s)
	// with similar names. If found, we will append an underscore followed by a
	// digit to make the name unique. Start with 0 and keep incrementing. This
	// allows a maximum of 1000 tables with same/similar name.
	if len(new_name) > 59 {
		new_name = new_name[0:58]
	}
	sql := "SELECT * FROM " + db.MSTR_TABLE + " WHERE user_id = $1 AND table_name LIKE $2 ORDER BY table_name DESC LIMIT 1"
	params := []interface{}{user_id, new_name + "%"}
	sres := ts.map_user_dao.FetchData(sql, params)
	if sres.Error != nil {
		return "", sres.Error
	}
	if len(sres.Rows) == 0 {
		// table with same name as "new_name" does not exist
		return new_name, nil
	} else {
		// table with same name as "new_name" exists
		new_name = sres.Rows[0]["table_name"]
		parts := strings.Split(new_name, "_")
		if len(parts) > 1 {
			// remove last part
			last := parts[len(parts)-1]
			parts = parts[:len(parts)-1]
			tmp, err := strconv.Atoi(last)
			if err == nil {
				// last part is a digit. increment it to make the name unique
				tmp += 1
				new_name = strings.Join(parts, "_") + "_" + strconv.Itoa(tmp)
			} else {
				// last part is *not* a digit, append 0
				new_name = new_name + "_0"
			}
		} else {
			new_name = new_name + "_0"
		}
	}

	return new_name, nil
}

func (ts *TableServiceImpl) addTable(name string, schema_name string, table_name string, user_id int) (int, error) {
	columns := make(map[string]interface{})
	columns["user_id"] = user_id
	columns["name"] = name
	columns["schema_name"] = schema_name
	columns["table_name"] = table_name
	columns["status"] = importer.UPLOADED

	return ts.map_user_dao.Insert(db.MasterTable(), columns, "id")
}

func MakeTableService(map_user_dao db.Dao) TableService {
	my_importer := importer.MakeImporter(map_user_dao)
	return &TableServiceImpl{map_user_dao, &my_importer}
}

func MakeTableServiceWithImporter(map_user_dao db.Dao, my_importer importer.Importer) TableService {
	return &TableServiceImpl{map_user_dao, my_importer}
}
