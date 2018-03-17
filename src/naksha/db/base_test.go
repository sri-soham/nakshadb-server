package db_test

import (
	"flag"
	"fmt"
	"naksha/db"
	"naksha/helper"
	"naksha/logger"
	"os"
	"os/exec"
	"testing"
)

var should_run_db_test bool
var dao db.Dao
var db_conn_str string
var path_to_sql string
var table string

func init() {
	dbname := flag.String("dbname", "", "Database name")
	dbuser := flag.String("dbuser", "", "Database user")
	dbpass := flag.String("dbpass", "", "Database password")
	// have to pass full path.
	sqlfile := flag.String("sqlfile", "", "Sql file path")
	flag.Parse()

	err_count := 0
	if len(*dbname) == 0 {
		err_count++
	}
	if len(*dbuser) == 0 {
		err_count++
	}
	if len(*dbpass) == 0 {
		err_count++
	}
	if len(*sqlfile) == 0 {
		err_count++
	}

	if err_count == 0 {
		should_run_db_test = true
		db_conn, err := db.InitDBConnection(*dbuser, *dbpass, *dbname, "disable")
		if err != nil {
			fmt.Println("Could not connect to database. Error = %s", err)
			os.Exit(1)
		}
		dao = db.DaoFromConn(db_conn)
		path_to_sql = *sqlfile
		db_conn_str = fmt.Sprintf("postgresql://%v:%v@localhost/%v?sslmode=disable", *dbuser, *dbpass, *dbname)
	} else if err_count == 4 {
		should_run_db_test = false
	} else {
		fmt.Println("All 4 flags should be passed: dbname, dbuser, dbpass and sqlfile")
		os.Exit(1)
	}
	table = "tbl_city"
}

func populateDb() {
	cmd := "psql"
	args := []string{db_conn_str, "-f", path_to_sql}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		fmt.Printf("output = %s. Error = %v\n", out, err)
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	if !should_run_db_test {
		return
	}

	code := m.Run()
	os.Exit(code)
}

func TestFetchData(t *testing.T) {
	populateDb()

	query := "SELECT * FROM " + table + " WHERE id > 90 ORDER BY id"
	params := make([]interface{}, 0)
	sr := dao.FetchData(query, params)
	if sr.Error != nil {
		t.Errorf("FetchData failed. Error = %v", sr.Error)
	}
	if len(sr.Rows) != 10 {
		t.Errorf("FetchData. Row count: expected(10), found(%d)", len(sr.Rows))
	}
	if len(sr.Rows[0]) != 5 {
		t.Errorf("FetchDAta. Column count: expected(5), found(%d)", len(sr.Rows[0]))
	}
	if sr.Rows[0]["id"] != "91" {
		t.Errorf("FetchData. Invalid id: expected(91), found(%d)", sr.Rows[0]["id"])
	}
}

func TestSelectWhere(t *testing.T) {
	populateDb()

	columns := []string{"id", "name"}
	where := make(map[string]interface{})
	where["id"] = 10
	sr := dao.SelectWhere(table, columns, where)
	if sr.Error != nil {
		t.Errorf("SelectWhere failed. Error = %v", sr.Error)
	}
	if len(sr.Rows) != 1 {
		t.Errorf("Rows: expected(1), found(%d)", len(sr.Rows))
	}
	if sr.Rows[0]["id"] != "10" {
		t.Errorf("Wrong id: expected = 10, found = %s", sr.Rows[0]["id"])
	}
	if sr.Rows[0]["name"] != "City 9" {
		t.Errorf("Wrond name: expected = 'City 9', found = '%v'", sr.Rows[0]["name"])
	}
}

func TestSelectAll(t *testing.T) {
	populateDb()

	columns := []string{"id", "name"}
	sr := dao.SelectAll(table, columns)
	if sr.Error != nil {
		t.Errorf("SelectAll failed. Error = %v", sr.Error)
	}
	if len(sr.Rows) != 100 {
		t.Errorf("SelectAll. Rows - expected(100), found(%d)", len(sr.Rows))
	}
}

func TestSelectSelect(t *testing.T) {
	populateDb()
	sq := db.SelectQuery{
		table,
		[]string{"*"},
		map[string]interface{}{"id": "53"},
		"",
		0,
		0,
	}
	sr := dao.SelectSelect(sq)
	if sr.Error != nil {
		t.Errorf("SelectSelect failed. Error = %v", sr.Error)
	}
	if len(sr.Rows) != 1 {
		t.Errorf("SelectSelect. Row count - expected(1), found(%d)", len(sr.Rows))
	}
}

func TestFindNotFound(t *testing.T) {
	populateDb()
	city, err := dao.Find(table, 134)
	if city != nil {
		t.Errorf("FindNotFound: city should be nil")
	}
	if err == nil {
		t.Errorf("FindNotFound: err should not be nil")
	}
}

func TestFind(t *testing.T) {
	populateDb()
	city, err := dao.Find(table, 34)
	if city == nil {
		t.Errorf("Find: city should not be nil")
	}
	if err != nil {
		t.Errorf("Find: err should be nil")
	}
	if city["name"] != "City 33" {
		t.Errorf("Find: name - expected('City 33'), found('%v')", city["name"])
	}
}

func TestFindWhereNotFound(t *testing.T) {
	populateDb()
	city, err := dao.FindWhere(table, map[string]interface{}{"id": 134})
	if city != nil {
		t.Errorf("FindWhereNotFound: city should be nil")
	}
	if err == nil {
		t.Errorf("FindWhereNotFound: err should not be nil")
	}
}

func TestFindWhere(t *testing.T) {
	populateDb()
	city, err := dao.FindWhere(table, map[string]interface{}{"id": 34})
	if city == nil {
		t.Errorf("FindWhere: city should not be nil")
	}
	if err != nil {
		t.Errorf("FindWhere: err should be nil")
	}
	if city["name"] != "City 33" {
		t.Errorf("Find: name - expected('City 33'), found('%v')", city["name"])
	}
}

func TestSelectOneNotFound(t *testing.T) {
	populateDb()
	city, err := dao.SelectOne(table, []string{"*"}, map[string]interface{}{"id": 134})
	if city != nil {
		t.Errorf("SelectOneNotFound: city should be nil")
	}
	if err != nil {
		t.Errorf("SelectOneNotFound: err should be nil")
	}
}

func TestSelectOne(t *testing.T) {
	populateDb()
	city, err := dao.SelectOne(table, []string{"*"}, map[string]interface{}{"id": 34})
	if city == nil {
		t.Errorf("SelectOne: city should not be nil")
	}
	if err != nil {
		t.Errorf("SelectOne: err should be nil")
	}
	if city["name"] != "City 33" {
		t.Errorf("SelectOne: expected('City 33'), found('%s')", city["name"])
	}
}

func TestModifyData(t *testing.T) {
	populateDb()

	name := "City 100"
	id := 1
	query := "UPDATE tbl_city SET name = $1 WHERE id = $2"
	params := []interface{}{name, id}
	_, err := dao.ModifyData(query, params)
	if err != nil {
		t.Errorf("ModifyData. Error: %v", err)
	}

	where := map[string]interface{}{
		"id": id,
	}
	sres := dao.SelectWhere(table, []string{"*"}, where)
	if sres.Error != nil {
		t.Errorf("ModifyData. Fetch query failed. Error: %v", sres.Error)
	}
	if len(sres.Rows) != 1 {
		t.Errorf("ModifyData. Row count: expected(1), found(%d)", len(sres.Rows))
	}
	if sres.Rows[0]["name"] != name {
		t.Errorf("ModifyData. Name not updated. Expected(%v), found(%v)", name, sres.Rows[0]["name"])
	}
}

func TestInsert(t *testing.T) {
	populateDb()

	values := make(map[string]interface{})
	values["name"] = "City 200"
	values["the_geom"] = "0101000020E6100000FD6A0E10CCD1D93FD881734694F6EC3F"
	id, err := dao.Insert(table, values, "id")
	if err != nil {
		t.Errorf("Insert. Error: %v", err)
	}
	if id <= 100 {
		t.Errorf("Insert. Invalid id. Should be greater than 100 but is %d", id)
	}
}

func TestInsertWithGeometry(t *testing.T) {
	populateDb()

	count1, err := dao.CountAll(table)
	if err != nil {
		t.Errorf("InsertWithGeometry. 1st count query failed. Error = %v", err)
	}
	values := make(map[string]interface{})
	values["name"] = "City 101"
	values["the_geom"] = "SRID=4326;POINT(17 78)"
	_, err = dao.InsertWithGeometry(table, values, "id")
	if err != nil {
		t.Errorf("InsertWithGeometry. Error: %v", err)
	}

	count2, err := dao.CountAll(table)
	if err != nil {
		t.Errorf("InsertWithGeometry. 2nd count query failed. Error = %v", err)
	}

	if (count1 + 1) != count2 {
		t.Errorf("InsertWithGeometry. Row not added. Row count: expected(%d), found(%d)", count2, (count1 + 1))
	}
}

func TestUpdate(t *testing.T) {
	populateDb()

	values := make(map[string]interface{})
	name := "City 1000"
	values["name"] = name
	where := make(map[string]interface{})
	id := 74
	where["id"] = id
	_, err := dao.Update(table, values, where)
	if err != nil {
		t.Errorf("Update. Error: %v", err)
	}

	sres := dao.SelectWhere(table, []string{"*"}, where)
	if sres.Rows[0]["name"] != name {
		t.Errorf("Update. Name not updated. Expected (%v), found(%v)", name, sres.Rows[0]["name"])
	}
}

func TestUpdateGeometry(t *testing.T) {
	populateDb()

	value := "SRID=4326;POINT(100.3454 200.9234)"
	where := map[string]interface{}{"id": 10}
	_, err := dao.UpdateGeometry(table, value, where)
	if err != nil {
		t.Errorf("UpdateGeometry. Error = %v", err)
	}
	expected := "0101000020E610000075029A081B1659400B24287E8C1D6940"
	sres := dao.SelectWhere(table, []string{"*"}, where)
	if sres.Rows[0]["the_geom"] != expected {
		t.Errorf("UpdateGeometry. Update failed: expected(%v), found(%v)", expected, sres.Rows[0]["the_geom"])
	}
}

func TestDelete(t *testing.T) {
	populateDb()

	where := map[string]interface{}{"id": 10}
	err := dao.Delete(table, where)
	if err != nil {
		t.Errorf("Delete. Error: %v", err)
	}

	sres := dao.SelectWhere(table, []string{"*"}, where)
	if len(sres.Rows) != 0 {
		t.Errorf("Delete. Row count: expected(0), found(%d)", len(sres.Rows))
	}
}

func TestCountAll(t *testing.T) {
	populateDb()

	count, err := dao.CountAll(table)
	if err != nil {
		t.Errorf("CountAll. Error: %v", err)
	}
	if count != 100 {
		t.Errorf("CountAll. Count: expected(100), found(%d)", count)
	}
}

func TestCountWhere(t *testing.T) {
	populateDb()

	where := map[string]interface{}{"id": 28}
	count, err := dao.CountWhere(table, where)
	if err != nil {
		t.Errorf("CountWhere. Error: %v", err)
	}
	if count != 1 {
		t.Errorf("CountWhere. Count: expected(1), found(%d)", count)
	}
}

func TestExec(t *testing.T) {
	populateDb()

	query := "ALTER TABLE " + table + " ADD COLUMN test VARCHAR(32) DEFAULT NULL"
	_, err := dao.Exec(query)
	if err != nil {
		t.Errorf("Exec. Error: %v", err)
	}

	where := map[string]interface{}{"id": 10}
	sres := dao.SelectWhere(table, []string{"*"}, where)
	if err != nil {
		t.Errorf("Exec. Select Error: %v", err)
	}
	if len(sres.Columns) != 6 {
		t.Errorf("Exec. Column count: expected(6), found(%d)", len(sres.Columns))
	}
}

func TestStoredProcNoResult(t *testing.T) {
	populateDb()

	proc_query := "SELECT naksha_test_1($1)"
	err := dao.StoredProcNoResult(proc_query, []interface{}{table})
	if err != nil {
		t.Errorf("StoredProcNoResult. Error = %v", err)
	}

	query := "SELECT * FROM " + table + " WHERE the_geom_webmercator IS NULL"
	sres := dao.FetchData(query, []interface{}{})
	if len(sres.Rows) != 0 {
		t.Errorf("StoredProcNoResult. the_geom_webmercator is empty")
	}
}

func TestTxTransaction1(t *testing.T) {
	populateDb()

	name1 := "City 1001"
	name2 := "City 1002"
	result := dao.TxTransaction(func(dao db.Dao) helper.Result {
		var trans_res helper.Result

		values1 := make(map[string]interface{})
		values1["name"] = name1
		where1 := make(map[string]interface{})
		where1["id"] = 1
		_, err := dao.Update(table, values1, where1)
		if err != nil {
			msg := fmt.Sprintf("Query failed 1. Error = %v", err)
			return helper.ErrorResultString(msg)
		}

		values2 := make(map[string]interface{})
		values2["name"] = name2
		where2 := make(map[string]interface{})
		where2["id"] = 2
		_, err = dao.Update(table, values2, where2)
		if err != nil {
			msg := fmt.Sprintf("Query failed 2. Error = %v", err)
			return helper.ErrorResultString(msg)
		}

		trans_res = helper.MakeSuccessResult()
		return trans_res
	})
	if !result.IsSuccess() {
		t.Errorf("TxTransaction1: Error: %v", result.GetErrors())
	}

	query := "SELECT id, name FROM " + table + " WHERE id <= 2 ORDER BY id ASC"
	sres := dao.FetchData(query, []interface{}{})
	if sres.Rows[0]["name"] != name1 {
		t.Errorf("TxTransaction1: 1- Name mistmatch: expected(%v), found(%v)", name1, sres.Rows[0]["name"])
	}
	if sres.Rows[1]["name"] != name2 {
		t.Errorf("TxTransaction1: 2- Name mistmatch: expected(%v), found(%v)", name2, sres.Rows[1]["name"])
	}
}

// First query of transaction fails
func TestTxTransaction2(t *testing.T) {
	populateDb()

	// If a query fails, error message will be logged to the error log.
	// If logger is not initiated, program will panic. To avoid that, initiating
	// the logger to the /tmp directory. /tmp directory is a linux only thing.
	logger.InitLoggers("/tmp")
	defer logger.CloseLoggers()

	name1 := "City 0"
	name2 := "City 1"
	result := dao.TxTransaction(func(dao db.Dao) helper.Result {
		var trans_res helper.Result

		values1 := make(map[string]interface{})
		values1["name"] = "City 1001"
		where1 := make(map[string]interface{})
		where1["uid"] = 1
		_, err := dao.Update(table, values1, where1)
		if err != nil {
			msg := fmt.Sprintf("Query failed 1. Error = %v", err)
			return helper.ErrorResultString(msg)
		}

		values2 := make(map[string]interface{})
		values2["name"] = "City 1002"
		where2 := make(map[string]interface{})
		where2["id"] = 2
		_, err = dao.Update(table, values2, where2)
		if err != nil {
			msg := fmt.Sprintf("Query failed 2. Error = %v", err)
			return helper.ErrorResultString(msg)
		}

		trans_res = helper.MakeSuccessResult()
		return trans_res
	})
	if result.IsSuccess() {
		t.Errorf("TxTransaction1: result should be error")
	}

	query := "SELECT id, name FROM " + table + " WHERE id <= 2 ORDER BY id ASC"
	sres := dao.FetchData(query, []interface{}{})
	if sres.Rows[0]["name"] != name1 {
		t.Errorf("TxTransaction1: 1- Name mistmatch: expected(%v), found(%v)", name1, sres.Rows[0]["name"])
	}
	if sres.Rows[1]["name"] != name2 {
		t.Errorf("TxTransaction1: 2- Name mistmatch: expected(%v), found(%v)", name2, sres.Rows[1]["name"])
	}
}

// Second query of transaction fails
func TestTxTransaction3(t *testing.T) {
	populateDb()

	// If a query fails, error message will be logged to the error log.
	// If logger is not initiated, program will panic. To avoid that, initiating
	// the logger to the /tmp directory. /tmp directory is a linux only thing.
	logger.InitLoggers("/tmp")
	defer logger.CloseLoggers()

	name1 := "City 0"
	name2 := "City 1"
	result := dao.TxTransaction(func(dao db.Dao) helper.Result {
		var trans_res helper.Result

		values1 := make(map[string]interface{})
		values1["name"] = "City 1001"
		where1 := make(map[string]interface{})
		where1["id"] = 1
		_, err := dao.Update(table, values1, where1)
		if err != nil {
			msg := fmt.Sprintf("Query failed 1. Error = %v", err)
			return helper.ErrorResultString(msg)
		}

		values2 := make(map[string]interface{})
		values2["name"] = "City 1002"
		where2 := make(map[string]interface{})
		where2["uid"] = 2
		_, err = dao.Update(table, values2, where2)
		if err != nil {
			msg := fmt.Sprintf("Query failed 2. Error = %v", err)
			return helper.ErrorResultString(msg)
		}

		trans_res = helper.MakeSuccessResult()
		return trans_res
	})
	if result.IsSuccess() {
		t.Errorf("TxTransaction1: Result should be error")
	}

	query := "SELECT id, name FROM " + table + " WHERE id <= 2 ORDER BY id ASC"
	sres := dao.FetchData(query, []interface{}{})
	if sres.Rows[0]["name"] != name1 {
		t.Errorf("TxTransaction1: 1- Name mistmatch: expected(%v), found(%v)", name1, sres.Rows[0]["name"])
	}
	if sres.Rows[1]["name"] != name2 {
		t.Errorf("TxTransaction1: 2- Name mistmatch: expected(%v), found(%v)", name2, sres.Rows[1]["name"])
	}
}
