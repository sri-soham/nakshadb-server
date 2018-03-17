package service

import (
	"fmt"
	"naksha"
	"naksha/db"
	"naksha/helper"
	"naksha/logger"
	"naksha/validator"
	"regexp"
	"strconv"
	"strings"
)

type MapService interface {
	Add(req *naksha.Request) helper.Result
	GetDetails(req *naksha.Request) helper.Result
	Update(req *naksha.Request) helper.Result
	UserMaps(req *naksha.Request) helper.Result
	ShowMap(req *naksha.Request) helper.Result
	Delete(req *naksha.Request) helper.Result
	BaseLayerUpdate(req *naksha.Request) helper.Result
	SearchTables(req *naksha.Request) helper.Result
	QueryData(req *naksha.Request) helper.Result
	AddLayer(req *naksha.Request) helper.Result
	DeleteLayer(req *naksha.Request) helper.Result
	UpdateHash(req *naksha.Request) helper.Result
	LayerOfTable(req *naksha.Request, tiler_url string) helper.Result
	MapBelongsToUser(map_id string, user_id int) bool
}

type MapServiceImpl struct {
	map_user_dao db.Dao
	api_user_dao db.Dao
}

func (ms *MapServiceImpl) Add(req *naksha.Request) helper.Result {
	var result helper.Result

	validator := validator.MakeMapValidator(req.Request)
	errs := validator.Validate()
	if len(errs) > 0 {
		result = helper.MakeErrorResult()
		result.SetErrors(errs)
		return result
	}

	user_id := helper.UserIDFromSession(req.Session)
	layer_id := req.PostFormValue("layer")
	result = ms.layerBelongsToUser(layer_id, user_id)
	if !result.IsSuccess() {
		return result
	}

	result = ms.map_user_dao.TxTransaction(func(dao db.Dao) helper.Result {
		var trans_res helper.Result

		map_values := make(map[string]interface{})
		map_values["user_id"] = user_id
		map_values["name"] = req.PostFormValue("name")
		map_values["hash"] = helper.RandomString(64)
		map_values["base_layer"] = helper.BASE_MAP_OSM
		map_id, err := dao.Insert(db.MasterMap(), map_values, "id")
		if err != nil {
			return helper.ErrorResultString("Query failed 3")
		}

		err_msg := fmt.Sprintf("Add: map_id = %d, layer_id = %v", map_id, layer_id)
		logger.LogMessage(err_msg)
		trans_res = ms.addLayer(dao, fmt.Sprintf("%d", map_id), layer_id)
		if !trans_res.IsSuccess() {
			return trans_res
		}

		trans_res = helper.MakeSuccessResult()
		trans_res.AddToData("redir_url", fmt.Sprintf("/maps/%v/show", map_id))
		return trans_res
	})

	return result
}

func (ms *MapServiceImpl) GetDetails(req *naksha.Request) helper.Result {
	var result helper.Result

	map_details, err := ms.map_user_dao.Find(db.MasterMap(), req.UriParams["id"])
	if err != nil {
		return helper.ErrorResultError(err)
	}

	sel := db.GetSelect()
	sel.Columns("mt.schema_name, mt.name AS table_name", "ml.id AS layer_id")
	sel.From(db.MasterTable() + " AS mt")
	sel.Join(db.MasterLayer()+" AS ml", "mt.id = ml.table_id")
	sel.Join(db.MasterMapLayer()+" AS mml", "ml.id = mml.layer_id")
	sel.Where("mml.map_id", req.UriParams["id"])
	sel.OrderBy("table_name")
	res2 := ms.map_user_dao.FetchData(sel.GetQuery(), sel.GetParams())
	if res2.Error != nil {
		return helper.ErrorResultString("Query failed 2")
	}

	result = helper.MakeSuccessResult()
	result.AddToData("map_details", map_details)
	result.AddToData("base_layers", helper.BaseLayers())
	result.AddToData("tables", res2.Rows)

	return result
}

func (ms *MapServiceImpl) Update(req *naksha.Request) helper.Result {
	var result helper.Result

	validator := validator.MakeMapValidator(req.Request)
	errs := validator.ValidateEdit()

	if len(errs) > 0 {
		result = helper.MakeErrorResult()
		result.SetErrors(errs)
		return result
	}

	values := make(map[string]interface{})
	values["name"] = req.PostFormValue("name")
	where := make(map[string]interface{})
	where["id"] = req.UriParams["id"]
	_, err := ms.map_user_dao.Update(db.MasterMap(), values, where)
	if err != nil {
		result = helper.ErrorResultString("Update failed")
	} else {
		result = helper.MakeSuccessResult()
	}

	return result
}

func (ms *MapServiceImpl) UserMaps(req *naksha.Request) helper.Result {
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
	where := make(map[string]interface{})
	where["user_id"] = user_id
	sq := db.SelectQuery{
		db.MasterMap(),
		[]string{"*"},
		where,
		"name",
		per_page,
		((page - 1) * per_page),
	}
	res1 := ms.map_user_dao.SelectSelect(sq)
	if res1.Error != nil {
		return helper.ErrorResultString("Query failed 1")
	}

	count, err := ms.map_user_dao.CountWhere(db.MasterMap(), where)
	if err != nil {
		return helper.ErrorResultString("Query failed 2")
	}

	paginator := helper.MakePaginator(page, per_page, count, "/maps/index?page={page}")
	result = helper.MakeSuccessResult()
	result.AddToData("maps", res1.Rows)
	result.AddToData("pagination_links", paginator.Links())
	result.AddToData("pagination_text", paginator.Text())

	return result
}

func (ms *MapServiceImpl) ShowMap(req *naksha.Request) helper.Result {
	var result helper.Result

	parts := strings.Split(req.UriParams["map_hash"], "-")
	if len(parts) < 2 {
		return helper.ErrorResultString("Invalid hash")
	}
	where := make(map[string]interface{})
	where["id"] = parts[0]
	where["hash"] = strings.Join(parts[1:], "-")
	map_details, err := ms.map_user_dao.FindWhere(db.MasterMap(), where)
	if err != nil {
		return helper.ErrorResultString("Could not get details")
	}

	user_details, err := ms.map_user_dao.Find(db.MasterUser(), map_details["user_id"])
	if err != nil {
		return helper.ErrorResultError(err)
	}

	sel := db.GetSelect()
	sel.Columns("mt.name", "mt.table_name")
	sel.From(db.MasterTable() + " AS mt")
	sel.Join(db.MasterLayer()+" AS ml", "ml.table_id = mt.id")
	sel.Join(db.MasterMapLayer()+" AS mml", "ml.id = mml.layer_id")
	sel.Where("mml.map_id", map_details["id"])
	ml_res := ms.map_user_dao.FetchData(sel.GetQuery(), sel.GetParams())
	if ml_res.Error != nil {
		return helper.ErrorResultString("Query failed 3")
	}
	if len(ml_res.Rows) == 0 {
		return helper.ErrorResultString("No layers for map")
	}

	queries := make([]string, 0)
	for _, row := range ml_res.Rows {
		query := fmt.Sprintf("SELECT '%v' AS table_name, ST_Extent(the_geom) AS xtnt FROM %v", row["table_name"], user_details["schema_name"]+"."+row["table_name"])
		queries = append(queries, query)
	}
	extents_query := strings.Join(queries, " UNION ALL ")
	ext_res := ms.map_user_dao.FetchData(extents_query, []interface{}{})
	if ext_res.Error != nil {
		return helper.ErrorResultString("Query failed 4")
	}
	extents := make([]string, 0)
	for _, ext := range ext_res.Rows {
		extents = append(extents, ext["table_name"])
		extents = append(extents, ext["xtnt"])
	}

	result = helper.MakeSuccessResult()
	result.AddToData("map_details", map_details)
	result.AddToData("is_google_maps", helper.IsGoogleMapsBaseLayer(map_details["base_layer"]))
	result.AddToData("is_bing_maps", helper.IsBingMapsBaseLayer(map_details["base_layer"]))
	result.AddToData("is_yandex_maps", helper.IsYandexMapsBaseLayer(map_details["base_layer"]))
	result.AddToData("user_details", user_details)
	result.AddToData("layer_data", ml_res.Rows)
	result.AddToData("extents", extents)

	return result
}

func (ms *MapServiceImpl) Delete(req *naksha.Request) helper.Result {
	var result helper.Result

	map_id := req.UriParams["id"]
	ml_where := map[string]interface{}{"map_id": map_id}
	mm_where := map[string]interface{}{"id": map_id}
	result = ms.map_user_dao.TxTransaction(func(txdao db.Dao) helper.Result {
		var trans_result helper.Result

		err := txdao.Delete(db.MasterMapLayer(), ml_where)
		if err != nil {
			return helper.ErrorResultString("Query failed 2")
		}

		err = txdao.Delete(db.MasterMap(), mm_where)
		if err != nil {
			trans_result = helper.ErrorResultString("Query failed 3")
		} else {
			trans_result = helper.MakeSuccessResult()
			trans_result.AddToData("redir_url", "/maps/index")
		}

		return trans_result
	})

	return result
}

func (ms *MapServiceImpl) BaseLayerUpdate(req *naksha.Request) helper.Result {
	var result helper.Result

	user_id := helper.UserIDFromSession(req.Session)
	user, err := ms.map_user_dao.Find(db.MasterUser(), user_id)
	if err != nil {
		return helper.ErrorResultError(err)
	}

	validator := validator.MakeMapValidator(req.Request)
	errs := validator.ValidateBaseLayer(helper.BaseLayers(), user)
	if len(errs) > 0 {
		result = helper.MakeErrorResult()
		result.SetErrors(errs)
		return result
	}

	base_layer := req.PostFormValue("base_layer")
	values := map[string]interface{}{"base_layer": base_layer}
	where := map[string]interface{}{"id": req.UriParams["id"]}
	_, err = ms.map_user_dao.Update(db.MasterMap(), values, where)
	if err != nil {
		result = helper.ErrorResultString("Could not update")
	} else {
		result = helper.MakeSuccessResult()
	}

	return result
}

func (ms *MapServiceImpl) SearchTables(req *naksha.Request) helper.Result {
	var result helper.Result

	table_frag := req.Request.FormValue("table_name")
	user_id := helper.UserIDFromSession(req.Session)
	map_id := req.UriParams["id"]
	query := "SELECT mt.name AS value, ml.id AS layer_id \n" +
		"      FROM " + db.MasterTable() + " AS mt \n" +
		" INNER JOIN " + db.MasterLayer() + " AS ml ON mt.id = ml.table_id \n" +
		"      WHERE mt.name ILIKE $1 \n" +
		"        AND mt.user_id = $2 \n" +
		"        AND ml.id NOT IN \n" +
		"           (SELECT layer_id FROM " + db.MasterMapLayer() + " WHERE map_id = $3) "
	params := []interface{}{table_frag + "%", user_id, map_id}
	sres := ms.map_user_dao.FetchData(query, params)
	if sres.Error != nil {
		result = helper.ErrorResultString("Query failed")
	} else {
		result = helper.MakeSuccessResult()
		result.AddToData("tables", sres.Rows)
	}

	return result
}

func (ms *MapServiceImpl) QueryData(req *naksha.Request) helper.Result {
	var result helper.Result

	query := req.Request.FormValue("query")
	if len(query) == 0 {
		return helper.ErrorResultString("Query is missing")
	}

	query = "SELECT * FROM ( " + query + " ) AS foo"
	sres := ms.api_user_dao.FetchData(query, []interface{}{})
	if sres.Error != nil {
		return helper.ErrorResultString("Query failed")
	}
	result = helper.MakeSuccessResult()
	result.AddToData("data", sres.Rows)

	return result
}

func (ms *MapServiceImpl) AddLayer(req *naksha.Request) helper.Result {
	var result helper.Result

	layer_id := req.PostFormValue("layer_id")
	if len(layer_id) == 0 {
		return helper.ErrorResultString("Layer id is required")
	}
	user_id := helper.UserIDFromSession(req.Session)
	result = ms.layerBelongsToUser(layer_id, user_id)
	if !result.IsSuccess() {
		return result
	}

	map_id := req.UriParams["id"]
	result = ms.addLayer(ms.map_user_dao, map_id, layer_id)

	return result
}

func (ms *MapServiceImpl) DeleteLayer(req *naksha.Request) helper.Result {
	var result helper.Result

	layer_id := req.PostFormValue("layer_id")
	if len(layer_id) == 0 {
		return helper.ErrorResultString("Layer id is required")
	}

	map_id := req.UriParams["id"]
	count_where := map[string]interface{}{"map_id": map_id}
	cnt, err := ms.map_user_dao.CountWhere(db.MasterMapLayer(), count_where)
	if err != nil {
		return helper.ErrorResultString("Query failed")
	}
	if cnt <= 1 {
		return helper.ErrorResultString("Map should have at least one layer")
	}

	where := make(map[string]interface{})
	where["map_id"] = map_id
	where["layer_id"] = layer_id
	err = ms.map_user_dao.Delete(db.MasterMapLayer(), where)
	if err == nil {
		result = helper.MakeSuccessResult()
	} else {
		result = helper.ErrorResultString("Query failed")
	}

	return result
}

func (ms *MapServiceImpl) UpdateHash(req *naksha.Request) helper.Result {
	var result helper.Result

	hash := req.PostFormValue("hash")
	hash_len := len(hash)
	errs := make([]string, 0)
	if hash_len == 0 {
		errs = append(errs, "Hash is required")
	} else if hash_len > 64 {
		errs = append(errs, "Max. allowed length is 64")
	}
	if len(errs) == 0 {
		re := regexp.MustCompile("^[[:alnum:]_-]{1,64}$")
		if !re.MatchString(hash) {
			errs = append(errs, "Only alphabets, digits, underscore and hypen are allowed")
		}
	}
	if len(errs) > 0 {
		return helper.ErrorResultString(strings.Join(errs, "<br />"))
	}

	map_id := req.UriParams["id"]
	result = ms.map_user_dao.TxTransaction(func(dao db.Dao) helper.Result {
		var trans_res helper.Result

		query := "SELECT * FROM " + db.MasterMap() + " WHERE hash = $1 OR id = $2 FOR UPDATE"
		where := []interface{}{hash, map_id}
		sres := dao.FetchData(query, where)
		if sres.Error != nil {
			return helper.ErrorResultString("Could not select rows")
		}
		/**
		 * If two rows are returned, then the asked-for hash is already in use by some other map.
		 * If only one row is returned:
		 *   - id of the row equals the map-id being update, (new) hash is available or the existing hash has not been modified.
		 *   - id of the row does not equal map-id, hash is being used another map and, more importantly, row with map-id has been deleted by some other session/transaction/action.
		 * If no rows are returned, map with map_id has been deleted.
		 */
		switch len(sres.Rows) {
		case 2:
			trans_res = helper.ErrorResultString("Hash is being used by some other map")
		case 1:
			if sres.Rows[0]["id"] == map_id {
				values := map[string]interface{}{"hash": hash}
				where := map[string]interface{}{"id": map_id}
				_, err := dao.Update(db.MasterMap(), values, where)
				if err == nil {
					trans_res = helper.MakeSuccessResult()
				} else {
					trans_res = helper.ErrorResultString("Update failed")
				}
			} else {
				trans_res = helper.ErrorResultString("Hash not available")
			}
		default:
			trans_res = helper.ErrorResultString("Map not available")
		}

		return trans_res
	})

	return result
}

func (ms *MapServiceImpl) LayerOfTable(req *naksha.Request, tiler_url string) helper.Result {
	sel := db.GetSelect()
	sel.Columns("ml.hash, ml.update_hash, ml.infowindow")
	sel.From(db.MasterTable() + " AS mt")
	sel.Join(db.MasterLayer()+" AS ml", "mt.id = ml.table_id")
	sel.Where("mt.schema_name", req.UriParams["schema"])
	sel.Where("mt.table_name", req.UriParams["table"])
	res := ms.map_user_dao.FetchData(sel.GetQuery(), sel.GetParams())
	if res.Error != nil {
		return helper.ErrorResultString("Query failed 2")
	}
	if len(res.Rows) != 1 {
		return helper.ErrorResultString("No such table")
	}
	result := helper.MakeSuccessResult()
	layer_url := tiler_url + "lyr/" + res.Rows[0]["hash"] + "-"
	layer_url += res.Rows[0]["update_hash"] + "/{z}/{x}/{y}.png"
	result.AddToData("layer_url", layer_url)
	result.AddToData("infowindow", res.Rows[0]["infowindow"])

	return result
}

func (ms *MapServiceImpl) MapBelongsToUser(map_id string, user_id int) bool {
	var belongs bool

	where := make(map[string]interface{})
	where["id"] = map_id
	where["user_id"] = user_id
	count, err := ms.map_user_dao.CountWhere(db.MasterMap(), where)
	if err == nil {
		belongs = (count == 1)
	} else {
		belongs = false
	}

	return belongs
}

func (ms *MapServiceImpl) addLayer(dao db.Dao, map_id string, layer_id string) helper.Result {
	var result helper.Result

	map_layer_values := make(map[string]interface{})
	map_layer_values["map_id"] = map_id
	map_layer_values["layer_id"] = layer_id
	map_layer_values["layer_index"] = 1
	_, err := dao.Insert(db.MasterMapLayer(), map_layer_values, "")
	if err != nil {
		result = helper.ErrorResultString("Query failed 4")
	} else {
		result = helper.MakeSuccessResult()
	}

	return result
}

func (ms *MapServiceImpl) layerBelongsToUser(layer_id string, user_id int) helper.Result {
	var result helper.Result

	sel := db.GetSelect()
	sel.Columns("COUNT(*) AS cnt")
	sel.From(db.MasterLayer() + " AS ml")
	sel.Join(db.MasterTable()+" AS mt", "ml.table_id = mt.id")
	sel.Where("ml.id", layer_id)
	sel.Where("mt.user_id", user_id)
	res1 := ms.map_user_dao.FetchData(sel.GetQuery(), sel.GetParams())
	if res1.Error != nil {
		return helper.ErrorResultString("Query failed 1")
	}
	cnt, _ := strconv.Atoi(res1.Rows[0]["cnt"])
	if cnt == 1 {
		result = helper.MakeSuccessResult()
	} else {
		result = helper.ErrorResultString("You cannot create this map")
	}

	return result
}

func MakeMapService(map_user_dao db.Dao, api_user_dao db.Dao) MapService {
	return &MapServiceImpl{map_user_dao, api_user_dao}
}
