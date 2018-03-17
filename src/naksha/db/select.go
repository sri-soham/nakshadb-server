package db

import (
	"fmt"
	"regexp"
	"strings"
)

type Select struct {
	where       []string
	joins       []string
	order_by    []string
	limit       int
	offset      int
	columns     []string
	params      []interface{}
	group_by    string
	having      string
	param_count int
	table       string
	block_start bool
}

func (s *Select) Columns(cols ...string) {
	s.columns = append(s.columns, cols...)
}

func (s *Select) From(table string) {
	s.table = table
}

func (s *Select) Where(frag string, val interface{}) {
	if len(s.where) > 0 {
		s.where = append(s.where, "AND")
	}
	if s.block_start {
		s.where = append(s.where, "(")
	}
	s.whereFrag(frag, val)
	s.block_start = false
}

func (s *Select) OrWhere(frag string, val interface{}) {
	if len(s.where) > 0 {
		s.where = append(s.where, "OR")
	}
	if s.block_start {
		s.where = append(s.where, "(")
	}
	s.whereFrag(frag, val)
	s.block_start = false
}

func (s *Select) BlockBegin() {
	s.block_start = true
}

func (s *Select) BlockEnd() {
	s.where = append(s.where, ")")
	s.block_start = false
}

func (s *Select) GroupBy(gb string) {
	s.group_by = gb
}

func (s *Select) Having(hv string) {
	s.having = hv
}

func (s *Select) Limit(l int) {
	s.limit = l
}

func (s *Select) Offset(o int) {
	s.offset = o
}

func (s *Select) OrderBy(ob string) {
	s.order_by = append(s.order_by, ob)
}

func (s *Select) Join(table string, on string) {
	j := fmt.Sprintf("INNER JOIN %v ON %v", table, on)
	s.joins = append(s.joins, j)
}

func (s *Select) LeftJoin(table string, on string) {
	j := fmt.Sprintf("LEFT OUTER JOIN %v ON %v", table, on)
	s.joins = append(s.joins, j)
}

func (s *Select) RightJoin(table string, on string) {
	j := fmt.Sprintf("RIGHT OUTER JOIN %v ON %v", table, on)
	s.joins = append(s.joins, j)
}

func (s *Select) GetQuery() string {
	query := "SELECT "
	if len(s.columns) == 0 {
		query += " * \n"
	} else {
		query += strings.Join(s.columns, ", ") + " \n"
	}
	query += "FROM " + s.table + " \n"
	if len(s.joins) > 0 {
		query += strings.Join(s.joins, "\n") + " \n"
	}
	if len(s.where) > 0 {
		query += "WHERE " + strings.Join(s.where, " ") + " \n"
	}
	if len(s.group_by) > 0 {
		query += "GROUP BY " + s.group_by + " \n"
	}
	if len(s.having) > 0 {
		query += "HAVING " + s.having + " \n"
	}
	if len(s.order_by) > 0 {
		query += "ORDER BY " + strings.Join(s.order_by, ", ") + " \n"
	}
	if s.limit > 0 {
		query += fmt.Sprintf("LIMIT %v ", s.limit)
		if s.offset > 0 {
			query += fmt.Sprintf("OFFSET %v ", s.offset)
		}
	}

	return query
}

func (s *Select) GetParams() []interface{} {
	return s.params
}

func (s *Select) whereFrag(frag string, val interface{}) {
	frag = strings.TrimSpace(frag)
	parts := regexp.MustCompile("\\s+").Split(frag, -1)
	part_len := len(parts)
	switch part_len {
	case 1:
		s.param_count++
		w := fmt.Sprintf("%v = $%v", parts[0], s.param_count)
		s.where = append(s.where, w)
		s.params = append(s.params, val)
	case 2:
		s.param_count++
		w := fmt.Sprintf("%v %v $%v", parts[0], parts[1], s.param_count)
		s.where = append(s.where, w)
		s.params = append(s.params, val)
	default:
		w := strings.Join(parts, " ")
		s.where = append(s.where, w)
	}
}

func GetSelect() Select {
	sel := Select{
		where:       make([]string, 0),
		joins:       make([]string, 0),
		order_by:    make([]string, 0),
		limit:       0,
		offset:      0,
		columns:     make([]string, 0),
		params:      make([]interface{}, 0),
		group_by:    "",
		having:      "",
		param_count: 0,
		table:       "",
		block_start: false,
	}

	return sel
}

func main() {
	sel := GetSelect()
	sel.Columns("name", "description", "id", "created_date")
	sel.From("tbl_tour")
	sel.BlockBegin()
	sel.Where("id", 10)
	sel.Where("created_date > ", "2017-10-10")
	sel.BlockEnd()
	sel.OrWhere("name IS NOT NULL", nil)
	sel.Join("tbl_user", "tbl_user.id = tbl_tour.user_id")
	sel.GroupBy("name")
	sel.Having("cnt > 10")
	sel.OrderBy("name")
	sel.OrderBy("description")
	sel.Offset(10)
	query := sel.GetQuery()
	fmt.Println(query)
	params := sel.GetParams()
	fmt.Println(params)
}
