package service

import (
	"fmt"
	"naksha"
	"naksha/db"
	"naksha/helper"
	"naksha/logger"
	"naksha/validator"
	"strconv"
	"strings"
	"time"
)

const PER_PAGE = 40

type TableRowService interface {
	Data(req *naksha.Request) helper.Result
	Update(req *naksha.Request) helper.Result
	Add(req *naksha.Request) helper.Result
	Delete(table_id string, id string) helper.Result
	Show(table_id string, id string) helper.Result
	TableBelongsToUser(table_id string, user_id int) bool
}

type TableRowServiceImpl struct {
	map_user_dao db.Dao
}

func (trs *TableRowServiceImpl) Data(req *naksha.Request) helper.Result {
	var result helper.Result

	table_name, err := trs.tableNameById(req.UriParams["table_id"])
	if err != nil {
		return helper.ErrorResultString("Query failed")
	}

	page, _ := strconv.Atoi(req.UriParams["page"])
	limit := PER_PAGE
	offset := (page - 1) * PER_PAGE
	order_column := req.Request.FormValue("order_column")
	order_type := req.Request.FormValue("order_type")
	if len(order_column) == 0 {
		order_column = "naksha_id"
	}
	if len(order_type) == 0 {
		order_type = "asc"
	} else {
		if order_type == "asc" || order_type == "desc" {
			// nothing to do
		} else {
			order_type = "asc"
		}
	}
	sq := db.SelectQuery{
		table_name,
		[]string{"*", "ST_AsEWKT(the_geom) AS the_geom"},
		nil,
		order_column + " " + order_type,
		limit,
		offset,
	}
	tmp := trs.map_user_dao.SelectSelect(sq)

	if tmp.Error != nil {
		return helper.ErrorResultString("Query failed 2")
	}
	rows := trs.tableRowsForDisplay(tmp)

	count, err := trs.map_user_dao.CountAll(table_name)
	if err != nil {
		return helper.ErrorResultString("Query failed 3")
	}

	result = helper.MakeSuccessResult()
	result.AddToData("rows", rows)
	result.AddToData("count", count)

	return result
}

func (trs *TableRowServiceImpl) Update(req *naksha.Request) helper.Result {
	var result helper.Result
	var err error
	var update_hash string

	table_name, err := trs.tableNameById(req.UriParams["table_id"])
	if err != nil {
		return helper.ErrorResultString("Could not get table name")
	}

	validator := validator.MakeTableRowValidator(req.Request)
	errs := validator.ValidateUpdate()
	if len(errs) == 0 {
		column := req.Request.PostFormValue("column")
		value := req.Request.PostFormValue("value")
		values := make(map[string]interface{})
		where := make(map[string]interface{})
		where["naksha_id"] = req.UriParams["id"]
		if column == "the_geom" {
			_, err = trs.map_user_dao.UpdateGeometry(table_name, value, where)
			update_hash, err = trs.updateUpdateHash(req.UriParams["table_id"])
			logger.LogMessage("updating geometry-type and style")
			go trs.updateGeometryTypeStyle(req.UriParams["table_id"])
		} else {
			values[column] = value
			_, err = trs.map_user_dao.Update(table_name, values, where)
		}
		if err == nil {
			result = helper.MakeSuccessResult()
			if column == "the_geom" {
				result.AddToData("update_hash", update_hash)
			}
		} else {
			result = helper.ErrorResultString("Query failed")
		}
	} else {
		result = helper.MakeErrorResult()
		result.SetErrors(errs)
	}

	return result
}

func (trs *TableRowServiceImpl) Add(req *naksha.Request) helper.Result {
	var result helper.Result
	var err error
	var naksha_id int
	var update_hash string

	table_name, err := trs.tableNameById(req.UriParams["table_id"])
	if err != nil {
		return helper.ErrorResultString("Could not get table name")
	}

	validator := validator.MakeTableRowValidator(req.Request)
	errs := validator.ValidateAdd()
	if len(errs) == 0 {
		values := make(map[string]interface{})
		with_geometry := req.PostFormValue("with_geometry")
		geometry := req.PostFormValue("geometry")

		st := helper.FormatTimestamp(time.Now())
		values["created_at"] = st
		values["updated_at"] = st
		if with_geometry == "1" {
			values["the_geom"] = geometry
			naksha_id, err = trs.map_user_dao.InsertWithGeometry(table_name, values, "naksha_id")
			update_hash, err = trs.updateUpdateHash(req.UriParams["table_id"])
			go trs.updateGeometryTypeStyle(req.UriParams["table_id"])
		} else {
			naksha_id, err = trs.map_user_dao.Insert(table_name, values, "naksha_id")
		}
		if err == nil {
			result = helper.MakeSuccessResult()
			values["naksha_id"] = naksha_id
			values["update_hash"] = update_hash
			result.AddToData("row", values)
		} else {
			result = helper.ErrorResultString("Query failed")
		}
	} else {
		result = helper.MakeErrorResult()
		result.SetErrors(errs)
	}

	return result
}

func (trs *TableRowServiceImpl) Delete(table_id string, id string) helper.Result {
	var result helper.Result

	table_name, err := trs.tableNameById(table_id)
	if err != nil {
		return helper.ErrorResultString("Could not get table name")
	}

	where := make(map[string]interface{})
	where["naksha_id"] = id
	err = trs.map_user_dao.Delete(table_name, where)
	if err == nil {
		update_hash, err := trs.updateUpdateHash(table_id)
		if err == nil {
			result = helper.MakeSuccessResult()
			result.AddToData("update_hash", update_hash)
		} else {
			result = helper.ErrorResultString("could not update hash")
		}
	} else {
		result = helper.ErrorResultString("Query failed")
	}

	return result
}

func (trs *TableRowServiceImpl) Show(table_id string, id string) helper.Result {
	var result helper.Result

	table_name, err := trs.tableNameById(table_id)
	if err != nil {
		return helper.ErrorResultString("Could not get table name")
	}

	ml_where := map[string]interface{}{"table_id": table_id}
	layer_details, err := trs.map_user_dao.FindWhere(db.MasterLayer(), ml_where)
	if err != nil {
		return helper.ErrorResultError(err)
	}
	iw, err := helper.StringToInfowindow(layer_details["infowindow"])
	if err != nil {
		return helper.ErrorResultString("Error parsing data")
	}

	where := map[string]interface{}{"naksha_id": id}
	res := trs.map_user_dao.SelectWhere(table_name, iw.Fields, where)
	if res.Error != nil {
		return helper.ErrorResultString("Could not get details")
	}
	if len(res.Rows) == 0 {
		return helper.ErrorResultString("No details")
	}

	result = helper.MakeSuccessResult()
	result.AddToData("data", res.Rows[0])

	return result
}

func (trs *TableRowServiceImpl) TableBelongsToUser(table_id string, user_id int) bool {
	var belongs bool

	where := make(map[string]interface{})
	where["id"] = table_id
	where["user_id"] = user_id
	count, err := trs.map_user_dao.CountWhere(db.MasterTable(), where)
	if err == nil {
		belongs = (count == 1)
	} else {
		belongs = false
	}

	return belongs
}

func (trs *TableRowServiceImpl) tableRowsForDisplay(res db.SelectResult) []map[string]string {
	rows := make([]map[string]string, 0)
	for _, c := range res.Rows {
		row := make(map[string]string)
		for k, v := range c {
			if k != "the_geom_webmercator" {
				row[k] = v
			}
		}
		rows = append(rows, row)
	}

	return rows
}

func (trs *TableRowServiceImpl) updateUpdateHash(table_id string) (string, error) {
	update_hash := helper.HashForMapUrl()
	values := map[string]interface{}{"update_hash": update_hash}
	where := map[string]interface{}{"table_id": table_id}

	_, err := trs.map_user_dao.Update(db.MasterLayer(), values, where)
	if err == nil {
		return update_hash, err
	} else {
		return "", err
	}
}

func (trs *TableRowServiceImpl) updateGeometryTypeStyle(table_id string) {
	table_name, err := trs.tableNameById(table_id)
	if err != nil {
		logger.LogMessage("TableRowService:updateGeomtryTypeStyle: could not get table name: " + table_id)
		return
	}
	ml_where := map[string]interface{}{"table_id": table_id}
	layer_details, err := trs.map_user_dao.FindWhere(db.MasterLayer(), ml_where)
	if err != nil {
		logger.LogMessage("TableRowService:updateGeomtryTypeStyle: could not get layer details: " + table_id)
		return
	}
	// Update style and geometry_type only if geometry_type is 'unknown'
	if layer_details["geometry_type"] != "unknown" {
		return
	}

	query := "SELECT ST_GeometryType(the_geom) AS geom_type FROM " + table_name + " WHERE ST_GeometryType(the_geom) IS NOT NULL LIMIT 1"
	tmp := trs.map_user_dao.FetchData(query, make([]interface{}, 0))
	if tmp.Error != nil {
		logger.LogMessage("TableRowService:updateGeometryTypeStyle: could not get geometry type: " + table_name)
		return
	}

	geom_type := tmp.Rows[0]["geom_type"]
	geom_type = strings.ToUpper(geom_type)
	logger.LogMessage(fmt.Sprintf("geom-type = %s", geom_type))
	values := make(map[string]interface{})
	values["updated_at"] = time.Now().Format(time.RFC3339)
	switch geom_type {
	case "ST_MULTIPOLYGON", "ST_POLYGON":
		values["geometry_type"] = "polygon"
		values["style"] = "<Rule>" +
			"<PolygonSymbolizer fill=\"#000000\" fill-opacity=\"0.75\" />" +
			"<LineSymbolizer stroke=\"#ffffff\" stroke-width=\"0.5\" stroke-opacity=\"1.0\" />" +
			"</Rule>"
	case "ST_MULTILINESTRING", "ST_LINESTRING":
		values["geometry_type"] = "linestring"
		values["style"] = "<Rule>" +
			"<LineSymbolizer stroke=\"#ffffff\" stroke-width=\"4\" stroke-opacity=\"1.0\" />" +
			"</Rule>"
	case "ST_MULTIPOINT", "ST_POINT":
		values["geometry_type"] = "point"
		values["style"] = "<Rule>" +
			"<MarkersSymbolizer fill=\"#000000\" stroke=\"#ffffff\" opacity=\"0.75\" stroke-width=\"1\" stroke-opacity=\"1.0\" width=\"10\" height=\"10\" marker-type=\"ellipse\" />" +
			"</Rule>"
	default:
		values["geometry_type"] = "unknown"
		values["style"] = "<Rule>" +
			"<PolygonSymbolizer fill=\"#000000\" fill-opacity=\"0.75\" />" +
			"<LineSymbolizer stroke=\"#ffffff\" stroke-width=\"0.5\" stroke-opacity=\"1.0\" />" +
			"<MarkersSymbolizer fill=\"#000000\" stroke=\"#ffffff\" opacity=\"0.75\" stroke-width=\"1\" stroke-opacity=\"1.0\" width=\"10\" height=\"10\" marker-type=\"ellipse\" />" +
			"</Rule>"
	}

	where := make(map[string]interface{})
	where["table_id"] = table_id
	_, err = trs.map_user_dao.Update(db.MasterLayer(), values, where)
	if err != nil {
		logger.LogMessage("TableRowService:updateGeometryTypeStyle: could update style, geometry_type: " + table_name)
	}
	logger.LogMessage("updated geometry type and style")
}

func (trs *TableRowServiceImpl) tableNameById(id string) (string, error) {
	details, err := trs.map_user_dao.Find(db.MasterTable(), id)
	if err != nil {
		return "", err
	} else {
		return helper.SchemaTableFromDetails(details), nil
	}
}

func MakeTableRowService(map_user_dao db.Dao) TableRowService {
	return &TableRowServiceImpl{map_user_dao}
}
