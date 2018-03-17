package db

import (
	"strings"
	"testing"
)

var dao DaoImpl

func init() {
	dao = DaoImpl{nil, false}
}

func TestColumns(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id")
	sel.Columns("name")
	sel.From("tbl_user")
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT id, name FROM tbl_user"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
}

func TestWhere1(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "name")
	sel.From("tbl_user")
	sel.Where("id", 10)
	query := standardizeString(sel.GetQuery())
	expected := "SELECT id, name FROM tbl_user WHERE id = $1"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 1 {
		t.Errorf("Param count: expected(1), found(%d)", len(params))
	}
}

func TestWhere2(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "name")
	sel.From("tbl_user")
	sel.Where("id > ", 10)
	query := standardizeString(sel.GetQuery())
	expected := "SELECT id, name FROM tbl_user WHERE id > $1"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
}

func TestWhere3(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "name")
	sel.From("tbl_user")
	sel.Where("id >= 10", nil)
	query := standardizeString(sel.GetQuery())
	expected := "SELECT id, name FROM tbl_user WHERE id >= 10"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 0 {
		t.Errorf("Param count: expected(0), found(%d)", len(params))
	}
}

func TestWhereAnd1(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "name")
	sel.From("tbl_user")
	sel.Where("id", 10)
	sel.Where("nid", 30)
	query := standardizeString(sel.GetQuery())
	expected := "SELECT id, name FROM tbl_user WHERE id = $1 AND nid = $2"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 2 {
		t.Errorf("Param count: expected(2), found(%d)", len(params))
	}
}

func TestWhereAnd2(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "name")
	sel.From("tbl_user")
	sel.Where("id", 10)
	sel.Where("nid", 30)
	sel.Where("name LIKE 'who%'", nil)
	query := standardizeString(sel.GetQuery())
	expected := "SELECT id, name FROM tbl_user WHERE id = $1 AND nid = $2 AND name LIKE 'who%'"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 2 {
		t.Errorf("Param count: expected(2), found(%d)", len(params))
	}
}

func TestWhereOr(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "name")
	sel.From("tbl_user")
	sel.Where("id", 10)
	sel.OrWhere("id", 30)
	query := standardizeString(sel.GetQuery())
	expected := "SELECT id, name FROM tbl_user WHERE id = $1 OR id = $2"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 2 {
		t.Errorf("Param count: expected(2), found(%d)", len(params))
	}
}

func TestWhereBlock1(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "first_name", "last_name")
	sel.From("tbl_user")
	sel.Where("country", "10")
	sel.BlockBegin()
	sel.Where("state", "100")
	sel.OrWhere("state", "200")
	sel.BlockEnd()
	query := standardizeString(sel.GetQuery())
	expected := "SELECT id, first_name, last_name FROM tbl_user WHERE country = $1 AND ( state = $2 OR state = $3 )"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 3 {
		t.Errorf("Param count: expected(3), found(%d)", len(params))
	}
}

func TestWhereBlock2(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "first_name", "last_name")
	sel.From("tbl_user")
	sel.Where("country", "10")
	sel.BlockBegin()
	sel.OrWhere("state", "100")
	sel.Where("city", "1000")
	sel.BlockEnd()
	query := standardizeString(sel.GetQuery())
	expected := "SELECT id, first_name, last_name FROM tbl_user WHERE country = $1 OR ( state = $2 AND city = $3 )"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 3 {
		t.Errorf("Param count: expected(3), found(%d)", len(params))
	}
}

func TestWhereBlock3(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "first_name", "last_name")
	sel.From("tbl_user")
	sel.BlockBegin()
	sel.Where("country", "10")
	sel.Where("state", "20")
	sel.BlockEnd()
	sel.OrWhere("first_name LIKE", "Who%")
	query := standardizeString(sel.GetQuery())
	expected := "SELECT id, first_name, last_name FROM tbl_user WHERE ( country = $1 AND state = $2 ) OR first_name LIKE $3"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 3 {
		t.Errorf("Param count: expected(3), found(%d)", len(params))
	}
}

func TestWhereBlock4(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "first_name", "last_name")
	sel.From("tbl_user")
	sel.Where("country", "10")
	sel.BlockBegin()
	sel.Where("state", "100")
	sel.OrWhere("state", "200")
	sel.BlockEnd()
	sel.Where("first_name LIKE", "who%")
	query := standardizeString(sel.GetQuery())
	expected := "SELECT id, first_name, last_name FROM tbl_user WHERE country = $1 AND ( state = $2 OR state = $3 ) AND first_name LIKE $4"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 4 {
		t.Errorf("Param count: expected(4), found(%d)", len(params))
	}
}

func TestLimit1(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id")
	sel.Columns("name")
	sel.From("tbl_user")
	sel.Limit(1)
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT id, name FROM tbl_user LIMIT 1"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
}

func TestLimit2(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id")
	sel.Columns("name")
	sel.From("tbl_user")
	sel.Where("country", "10")
	sel.Limit(1)
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT id, name FROM tbl_user WHERE country = $1 LIMIT 1"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 1 {
		t.Errorf("Param count: expected(1), found(%d)", len(params))
	}
}

func TestOffset1(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id")
	sel.Columns("name")
	sel.From("tbl_user")
	sel.Limit(10)
	sel.Offset(20)
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT id, name FROM tbl_user LIMIT 10 OFFSET 20"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
}

func TestOffset2(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id")
	sel.Columns("name")
	sel.From("tbl_user")
	sel.Where("country", "10")
	sel.Limit(10)
	sel.Offset(20)
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT id, name FROM tbl_user WHERE country = $1 LIMIT 10 OFFSET 20"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 1 {
		t.Errorf("Param count: expected(1), found(%d)", len(params))
	}
}

func TestGroupBy1(t *testing.T) {
	sel := GetSelect()
	sel.Columns("country", "COUNT(*) AS cnt")
	sel.From("tbl_user")
	sel.GroupBy("country")
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT country, COUNT(*) AS cnt FROM tbl_user GROUP BY country"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 0 {
		t.Errorf("Param count: expected(0), found(%d)", len(params))
	}
}

func TestGroupBy2(t *testing.T) {
	sel := GetSelect()
	sel.Columns("country", "COUNT(*) AS cnt")
	sel.From("tbl_user")
	sel.Where("id > 200", nil)
	sel.GroupBy("country")
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT country, COUNT(*) AS cnt FROM tbl_user WHERE id > 200 GROUP BY country"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 0 {
		t.Errorf("Param count: expected(0), found(%d)", len(params))
	}
}

func TestHaving1(t *testing.T) {
	sel := GetSelect()
	sel.Columns("country", "COUNT(*) AS cnt")
	sel.From("tbl_user")
	sel.GroupBy("country")
	sel.Having("cnt > 10")
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT country, COUNT(*) AS cnt FROM tbl_user GROUP BY country HAVING cnt > 10"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 0 {
		t.Errorf("Param count: expected(0), found(%d)", len(params))
	}
}

func TestHaving2(t *testing.T) {
	sel := GetSelect()
	sel.Columns("country", "COUNT(*) AS cnt")
	sel.From("tbl_user")
	sel.Where("id > ", 200)
	sel.GroupBy("country")
	sel.Having("cnt > 10")
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT country, COUNT(*) AS cnt FROM tbl_user WHERE id > $1 GROUP BY country HAVING cnt > 10"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 1 {
		t.Errorf("Param count: expected(1), found(%d)", len(params))
	}
}

func TestOrderBy1(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "first_name", "last_name")
	sel.From("tbl_user")
	sel.OrderBy("first_name")
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT id, first_name, last_name FROM tbl_user ORDER BY first_name"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 0 {
		t.Errorf("Param count: expected(0), found(%d)", len(params))
	}
}

func TestOrderBy2(t *testing.T) {
	sel := GetSelect()
	sel.Columns("id", "first_name", "last_name")
	sel.From("tbl_user")
	sel.Where("id > ", 10)
	sel.OrderBy("first_name")
	sel.OrderBy("last_name")
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT id, first_name, last_name FROM tbl_user WHERE id > $1 ORDER BY first_name, last_name"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 1 {
		t.Errorf("Param count: expected(1), found(%d)", len(params))
	}
}

func TestJoin(t *testing.T) {
	sel := GetSelect()
	sel.Columns("tu.id", "tu.first_name", "tu.last_name", "td.name")
	sel.From("tbl_user AS tu")
	sel.Join("tbl_department AS td", "tu.department_id = td.id")
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT tu.id, tu.first_name, tu.last_name, td.name FROM tbl_user AS tu INNER JOIN tbl_department AS td ON tu.department_id = td.id"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 0 {
		t.Errorf("Param count: expected(0), found(%d)", len(params))
	}
}

func TestLeftJoin(t *testing.T) {
	sel := GetSelect()
	sel.Columns("tu.id", "tu.first_name", "tu.last_name", "td.name")
	sel.From("tbl_user AS tu")
	sel.LeftJoin("tbl_department AS td", "tu.department_id = td.id")
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT tu.id, tu.first_name, tu.last_name, td.name FROM tbl_user AS tu LEFT OUTER JOIN tbl_department AS td ON tu.department_id = td.id"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 0 {
		t.Errorf("Param count: expected(0), found(%d)", len(params))
	}
}

func TestRightJoin(t *testing.T) {
	sel := GetSelect()
	sel.Columns("tu.id", "tu.first_name", "tu.last_name", "td.name")
	sel.From("tbl_user AS tu")
	sel.RightJoin("tbl_department AS td", "tu.department_id = td.id")
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT tu.id, tu.first_name, tu.last_name, td.name FROM tbl_user AS tu RIGHT OUTER JOIN tbl_department AS td ON tu.department_id = td.id"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 0 {
		t.Errorf("Param count: expected(0), found(%d)", len(params))
	}
}

func TestAll(t *testing.T) {
	sel := GetSelect()
	sel.Columns("tu.id", "tu.first_name", "tu.last_name", "td.name")
	sel.From("tbl_user AS tu")
	sel.Join("tbl_department AS td", "td.id = tu.department_id")
	sel.Where("tu.id >", 10)
	sel.Where("td.id = 3", nil)
	sel.OrderBy("tu.first_name")
	sel.Limit(10)
	sel.Offset(20)
	query := sel.GetQuery()
	query = standardizeString(query)
	expected := "SELECT tu.id, tu.first_name, tu.last_name, td.name " +
		"FROM tbl_user AS tu " +
		"INNER JOIN tbl_department AS td ON td.id = tu.department_id " +
		"WHERE tu.id > $1 AND td.id = 3 ORDER BY tu.first_name LIMIT 10 OFFSET 20"
	if query != expected {
		t.Errorf("Expected: %s\nFound: %s\n", expected, query)
	}
	params := sel.GetParams()
	if len(params) != 1 {
		t.Errorf("Param count: expected(1), found(%d)", len(params))
	}
}

func TestSQSelect(t *testing.T) {
	sq := SelectQuery{
		"tbl_user",
		[]string{"*"},
		nil,
		"",
		0,
		0,
	}
	query, params := dao.buildSelectQuery(sq)
	query = standardizeString(query)
	expected := "SELECT * FROM tbl_user"
	if query != expected {
		t.Errorf("SQSelect. Expected (%v), found(%v)", expected, query)
	}
	if len(params) != 0 {
		t.Errorf("SQSelect. Param count: expected(0), found(%d)", len(params))
	}
}

func TestSQSelectWhere1(t *testing.T) {
	sq := SelectQuery{
		"tbl_user",
		[]string{"*"},
		map[string]interface{}{"id": 10},
		"",
		0,
		0,
	}
	query, params := dao.buildSelectQuery(sq)
	query = standardizeString(query)
	expected := "SELECT * FROM tbl_user WHERE id = $1"
	if query != expected {
		t.Errorf("SQSelectWhere1. Expected (%v), found(%v)", expected, query)
	}
	if len(params) != 1 {
		t.Errorf("SQSelectWhere1. Param count: expected(1), found(%d)", len(params))
	}
}

func TestSQSelectWhere2(t *testing.T) {
	sq := SelectQuery{
		"tbl_user",
		[]string{"id", "first_name", "last_name"},
		map[string]interface{}{"country": 10, "state": 20},
		"",
		0,
		0,
	}
	query, params := dao.buildSelectQuery(sq)
	query = standardizeString(query)
	expected1 := "SELECT id, first_name, last_name FROM tbl_user WHERE country = $1 AND state = $2"
	expected2 := "SELECT id, first_name, last_name FROM tbl_user WHERE state = $1 AND country = $2"
	// map is not ordered, so country or state might appear first in the WHERE clause
	if query != expected1 && query != expected2 {
		t.Errorf("SQSelectWhere2. Expected ('%v' OR '%v'), found('%v')", expected1, expected2, query)
	}
	if len(params) != 2 {
		t.Errorf("SQSelectWhere2. Param count: expected(2), found(%d)", len(params))
	}
}

func TestSQSelectWhereOrder(t *testing.T) {
	sq := SelectQuery{
		"tbl_user",
		[]string{"id", "first_name", "last_name"},
		map[string]interface{}{"country": 10},
		"first_name",
		0,
		0,
	}
	query, params := dao.buildSelectQuery(sq)
	query = standardizeString(query)
	expected := "SELECT id, first_name, last_name FROM tbl_user WHERE country = $1 ORDER BY first_name"
	if query != expected {
		t.Errorf("SQSelectWhereOrder. Expected (%v), found(%v)", expected, query)
	}
	if len(params) != 1 {
		t.Errorf("SQSelectWhereOrder. Param count: expected(1), found(%d)", len(params))
	}
}

func TestSQSelectWhereOrderLimit(t *testing.T) {
	sq := SelectQuery{
		"tbl_user",
		[]string{"*"},
		map[string]interface{}{"country": 10, "state": 12},
		"first_name",
		10,
		0,
	}
	query, params := dao.buildSelectQuery(sq)
	query = standardizeString(query)
	expected1 := "SELECT * FROM tbl_user WHERE country = $1 AND state = $2 ORDER BY first_name LIMIT 10"
	expected2 := "SELECT * FROM tbl_user WHERE state = $1 AND country = $2 ORDER BY first_name LIMIT 10"
	// map is not ordered, so country or state might appear first in the WHERE clause
	if query != expected1 && query != expected2 {
		t.Errorf("SQSelectWhereOrderLimit. Expected ('%v' or '%v'), found(%v)", expected1, expected2, query)
	}
	if len(params) != 2 {
		t.Errorf("SQSelectWhereOrderLimit. Param count: expected(2), found(%d)", len(params))
	}
}

func TestSQSelectWhereOrderLimitOffset(t *testing.T) {
	sq := SelectQuery{
		"tbl_user",
		[]string{"*"},
		map[string]interface{}{"state": 10, "country": 12},
		"first_name",
		10,
		20,
	}
	query, params := dao.buildSelectQuery(sq)
	query = standardizeString(query)
	expected1 := "SELECT * FROM tbl_user WHERE country = $1 AND state = $2 ORDER BY first_name LIMIT 10 OFFSET 20"
	expected2 := "SELECT * FROM tbl_user WHERE state = $1 AND country = $2 ORDER BY first_name LIMIT 10 OFFSET 20"
	if query != expected1 && query != expected2 {
		t.Errorf("SQSelectWhereOrderLimitOffset. Expected ('%v' or '%v'), found(%v)", expected1, expected2, query)
	}
	if len(params) != 2 {
		t.Errorf("SQSelectWhereOrderLimitOffset. Param count: expected(2), found(%d)", len(params))
	}
}

func TestSQSelectOrder(t *testing.T) {
	sq := SelectQuery{
		"tbl_user",
		[]string{"*"},
		nil,
		"first_name",
		0,
		0,
	}
	query, params := dao.buildSelectQuery(sq)
	query = standardizeString(query)
	expected := "SELECT * FROM tbl_user ORDER BY first_name"
	if query != expected {
		t.Errorf("SQSelectOrder. Expected (%v), found(%v)", expected, query)
	}
	if len(params) != 0 {
		t.Errorf("SQSelectOrder. Param count: expected(0), found(%d)", len(params))
	}
}

func TestSQSelectOrderLimit(t *testing.T) {
	sq := SelectQuery{
		"tbl_user",
		[]string{"id", "name", "city", "state", "country"},
		nil,
		"name",
		20,
		0,
	}
	query, params := dao.buildSelectQuery(sq)
	query = standardizeString(query)
	expected := "SELECT id, name, city, state, country FROM tbl_user ORDER BY name LIMIT 20"
	if query != expected {
		t.Errorf("SQSelectOrderLimit. Expected (%v), found(%v)", expected, query)
	}
	if len(params) != 0 {
		t.Errorf("SQSelectOrderLimit. Param count: expected(0), found(%d)", len(params))
	}
}

func TestSQSelectOrderLimitOffset(t *testing.T) {
	sq := SelectQuery{
		"tbl_user",
		[]string{"id", "name", "city", "state", "country"},
		nil,
		"name",
		20,
		100,
	}
	query, params := dao.buildSelectQuery(sq)
	query = standardizeString(query)
	expected := "SELECT id, name, city, state, country FROM tbl_user ORDER BY name LIMIT 20 OFFSET 100"
	if query != expected {
		t.Errorf("SQSelectOrderLimitOffset. Expected (%v), found(%v)", expected, query)
	}
	if len(params) != 0 {
		t.Errorf("SQSelectOrderLimitOffset. Param count: expected(0), found(%d)", len(params))
	}
}

func standardizeString(str string) string {
	str = strings.TrimSpace(str)
	parts := strings.Fields(str)
	return strings.Join(parts, " ")
}
