package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"naksha/helper"
	"naksha/logger"
	"strconv"
	"strings"
	"time"
)

type Dao interface {
	FetchData(query string, params []interface{}) SelectResult
	SelectAll(table string, columns []string) SelectResult
	SelectWhere(table string, columns []string, where map[string]interface{}) SelectResult
	SelectSelect(sq SelectQuery) SelectResult
	SelectOne(table string, columns []string, where map[string]interface{}) (map[string]string, error)
	Find(table string, id interface{}) (map[string]string, error)
	FindWhere(table string, where map[string]interface{}) (map[string]string, error)
	ModifyData(query string, params []interface{}) (sql.Result, error)
	Insert(table string, values map[string]interface{}, auto_increment_column string) (int, error)
	InsertWithGeometry(table string, values map[string]interface{}, auto_increment_column string) (int, error)
	Update(table string, values map[string]interface{}, where map[string]interface{}) (sql.Result, error)
	UpdateGeometry(table string, value string, where map[string]interface{}) (sql.Result, error)
	Delete(table string, where map[string]interface{}) error
	CountAll(table string) (int, error)
	CountWhere(table string, where map[string]interface{}) (int, error)
	Exec(query string) (sql.Result, error)
	StoredProcNoResult(query string, params []interface{}) error
	TxTransaction(tx_func func(Dao) helper.Result) helper.Result
}

type DaoImpl struct {
	conn  Queryable
	is_db bool
}

func (di *DaoImpl) FetchData(query string, params []interface{}) SelectResult {
	rows, err := di.conn.Query(query, params...)
	if err != nil {
		logger.LogMessage(fmt.Sprintf("query failed: %v\nerrors: %v", query, err))
		return di.makeErrorResult(errors.New("Query failed"))
	}

	defer rows.Close()

	column_names, err1 := rows.Columns()
	if err1 != nil {
		logger.LogMessage(fmt.Sprintf("Could not get columns. query: %v\nerrors: %v", query, err1))
		return di.makeErrorResult(errors.New("Could not get columns"))
	}
	col_len := len(column_names)

	results := make([]map[string]string, 0)
	for rows.Next() {
		columns := make([]interface{}, col_len)
		column_pointers := make([]interface{}, col_len)
		for i := 0; i < col_len; i++ {
			column_pointers[i] = &columns[i]
		}
		err2 := rows.Scan(column_pointers...)
		if err2 != nil {
			logger.LogMessage(fmt.Sprintf("Could not scan. query:%v\nerror: %v", query, err2))
			return di.makeErrorResult(errors.New("Could not fetch results"))
		}
		result := make(map[string]string)
		for i := 0; i < col_len; i++ {
			if columns[i] == nil {
				result[column_names[i]] = ""
			} else {
				switch columns[i].(type) {
				case int, int32, int64:
					result[column_names[i]] = fmt.Sprintf("%v", columns[i])
				case float64:
					result[column_names[i]] = strconv.FormatFloat(columns[i].(float64), 'f', -1, 64)
				case time.Time:
					result[column_names[i]] = helper.FormatTimestamp(columns[i].(time.Time))
				case bool:
					if columns[i].(bool) {
						result[column_names[i]] = "t"
					} else {
						result[column_names[i]] = "f"
					}
				default:
					result[column_names[i]] = fmt.Sprintf("%s", columns[i])
				}
			}
		}
		results = append(results, result)
	}

	return SelectResult{nil, column_names, results}
}

func (di *DaoImpl) SelectAll(table string, columns []string) SelectResult {
	sq := SelectQuery{table, columns, nil, "", 0, 0}
	query, params := di.buildSelectQuery(sq)
	return di.FetchData(query, params)
}

func (di *DaoImpl) SelectWhere(table string, columns []string, where map[string]interface{}) SelectResult {
	sq := SelectQuery{table, columns, where, "", 0, 0}
	query, params := di.buildSelectQuery(sq)
	return di.FetchData(query, params)
}

func (di *DaoImpl) SelectSelect(sq SelectQuery) SelectResult {
	query, params := di.buildSelectQuery(sq)
	return di.FetchData(query, params)
}

// Returns error when record is not found
func (di *DaoImpl) Find(table string, id interface{}) (map[string]string, error) {
	where := map[string]interface{}{"id": id}
	res := di.SelectWhere(table, []string{"*"}, where)
	if res.Error == nil {
		if len(res.Rows) == 0 {
			err := errors.New("No such record")
			return nil, err
		} else {
			return res.Rows[0], nil
		}
	} else {
		return nil, res.Error
	}
}

// Returns error when record is not found
func (di *DaoImpl) FindWhere(table string, where map[string]interface{}) (map[string]string, error) {
	res := di.SelectWhere(table, []string{"*"}, where)
	if res.Error == nil {
		if len(res.Rows) == 0 {
			err := errors.New("No such record")
			return nil, err
		} else {
			return res.Rows[0], nil
		}
	} else {
		return nil, res.Error
	}
}

// returns nil (instead of error) when record is not found
func (di *DaoImpl) SelectOne(table string, columns []string, where map[string]interface{}) (map[string]string, error) {
	res := di.SelectWhere(table, columns, where)
	if res.Error == nil {
		if len(res.Rows) == 0 {
			return nil, nil
		} else {
			return res.Rows[0], nil
		}
	} else {
		return nil, res.Error
	}
}

func (di *DaoImpl) ModifyData(query string, params []interface{}) (sql.Result, error) {
	result, err := di.conn.Exec(query, params...)
	if err != nil {
		logger.LogMessage(fmt.Sprintf("query failed: %v\nerror: %v", query, err))
		return result, errors.New("Query failed")
	} else {
		return result, nil
	}
}

func (di *DaoImpl) Insert(table string, values map[string]interface{}, auto_increment_column string) (int, error) {
	var query string
	var columns = make([]string, 0)
	var placeholders = make([]string, 0)
	var params = make([]interface{}, 0)
	var err error

	count := 0
	for k, v := range values {
		count++
		columns = append(columns, k)
		placeholders = append(placeholders, fmt.Sprintf("$%v", count))
		params = append(params, v)
	}
	last_insert_id := 0
	query = "INSERT INTO " + table + "(" + strings.Join(columns, ", ") + ") " +
		"VALUES (" + strings.Join(placeholders, ", ") + ") "
	if len(auto_increment_column) > 0 {
		query += " RETURNING " + auto_increment_column
		err = di.conn.QueryRow(query, params...).Scan(&last_insert_id)
	} else {
		_, err = di.conn.Exec(query, params...)
	}
	if err != nil {
		logger.LogMessage(fmt.Sprintf("query failed: %v\nerror: %v", query, err))
		return last_insert_id, errors.New("Query failed")
	} else {
		return last_insert_id, nil
	}
}

func (di *DaoImpl) InsertWithGeometry(table string, values map[string]interface{}, auto_increment_column string) (int, error) {
	var query string
	var columns = make([]string, 0)
	var placeholders = make([]string, 0)
	var params = make([]interface{}, 0)
	var err error

	count := 0
	for k, v := range values {
		count++
		columns = append(columns, k)
		if k == "the_geom" {
			placeholders = append(placeholders, fmt.Sprintf("ST_GeomFromEWKT($%v)", count))
		} else {
			placeholders = append(placeholders, fmt.Sprintf("$%v", count))
		}
		params = append(params, v)
	}
	last_insert_id := 0
	query = "INSERT INTO " + table + "(" + strings.Join(columns, ", ") + ") " +
		"VALUES (" + strings.Join(placeholders, ", ") + ") "
	if len(auto_increment_column) > 0 {
		query += " RETURNING " + auto_increment_column
		err = di.conn.QueryRow(query, params...).Scan(&last_insert_id)
	} else {
		_, err = di.conn.Exec(query, params...)
	}
	if err != nil {
		logger.LogMessage(fmt.Sprintf("query failed: %v\nerror: %v", query, err))
		return last_insert_id, errors.New("Query failed")
	} else {
		return last_insert_id, nil
	}
}

func (di *DaoImpl) Update(table string, values map[string]interface{}, where map[string]interface{}) (sql.Result, error) {
	var query string

	set_parts, set_params := di.partsParamsFromMap(values, 1)
	part_len := len(set_parts)
	where_parts, where_params := di.partsParamsFromMap(where, part_len+1)
	all_params := make([]interface{}, 0)
	all_params = append(all_params, set_params...)
	all_params = append(all_params, where_params...)

	query = "UPDATE " + table + " SET " +
		strings.Join(set_parts, ", ")
	if len(where_parts) > 0 {
		query += " WHERE " + strings.Join(where_parts, " AND ")
	}

	result, err := di.conn.Exec(query, all_params...)
	if err != nil {
		logger.LogMessage(fmt.Sprintf("query failed: %v\nerror: %v", query, err))
		return result, errors.New("Query failed")
	} else {
		return result, nil
	}
}

func (di *DaoImpl) UpdateGeometry(table string, value string, where map[string]interface{}) (sql.Result, error) {
	var query string
	params := make([]interface{}, 0)

	if len(value) > 0 {
		query = "UPDATE " + table + " SET the_geom = ST_GeomFromEWKT($1) "
		params = append(params, value)
	} else {
		// geometry columns do not accept empty strings
		query = "UPDATE " + table + " SET the_geom = NULL "
	}

	where_parts, where_params := di.partsParamsFromMap(where, len(params)+1)
	params = append(params, where_params...)

	if len(where_parts) > 0 {
		query += " WHERE " + strings.Join(where_parts, " AND ")
	}

	result, err := di.conn.Exec(query, params...)
	if err != nil {
		logger.LogMessage(fmt.Sprintf("query failed: %v\nerror: %v", query, err))
		return result, errors.New("Query failed")
	} else {
		return result, nil
	}
}

func (di *DaoImpl) Delete(table string, where map[string]interface{}) error {
	var query string

	parts, params := di.partsParamsFromMap(where, 1)
	query = "DELETE FROM " + table
	if len(parts) > 0 {
		query += " WHERE " + strings.Join(parts, " AND ")
	}
	_, err := di.conn.Exec(query, params...)
	if err != nil {
		logger.LogMessage(fmt.Sprintf("query failed: %v\nerror: %v", query, err))
		return errors.New("Query failed")
	} else {
		return nil
	}
}

func (di *DaoImpl) CountAll(table string) (int, error) {
	var count int
	var params = make([]interface{}, 0)

	query := "SELECT COUNT(*) AS cnt FROM " + table
	err := di.conn.QueryRow(query, params...).Scan(&count)
	if err != nil {
		logger.LogMessage(fmt.Sprintf("Query failed: %v\nerror: %v", query, err))
		return count, errors.New("Query failed")
	} else {
		return count, nil
	}
}

func (di *DaoImpl) CountWhere(table string, where map[string]interface{}) (int, error) {
	var count int
	var query string

	parts, params := di.partsParamsFromMap(where, 1)
	query = "SELECT COUNT(*) AS cnt FROM " + table +
		" WHERE " + strings.Join(parts, " AND ")
	err := di.conn.QueryRow(query, params...).Scan(&count)
	if err != nil {
		logger.LogMessage(fmt.Sprintf("Query failed: %v\nerror: %v", query, err))
		return count, errors.New("Query failed")
	} else {
		return count, nil
	}
}

func (di *DaoImpl) Exec(query string) (sql.Result, error) {
	params := make([]interface{}, 0)
	res, err := di.conn.Exec(query, params...)
	if err != nil {
		logger.LogMessage(fmt.Sprintf("Query failed: %v\nerror: %v", query, err))
		return res, errors.New("Query failed")
	} else {
		return res, nil
	}
}

func (di *DaoImpl) StoredProcNoResult(query string, params []interface{}) error {
	var some_val interface{}
	err := di.conn.QueryRow(query, params...).Scan(&some_val)
	if err != nil {
		logger.LogMessage(fmt.Sprintf("Query failed: %v\nerror: %v\nparams=%v", query, err, params))
		return errors.New("Query failed")
	} else {
		return nil
	}
}

func (di *DaoImpl) TxTransaction(tx_func func(Dao) helper.Result) helper.Result {
	var result helper.Result

	if !di.is_db {
		return helper.ErrorResultString("Transactions not allowd on Tx")
	}
	tx, err := di.conn.Begin()
	if err != nil {
		return helper.ErrorResultString("Query failed 1")
	}

	defer func() {
		if p := recover(); p != nil {
			logger.LogMessage("Tx. Panicked. rolling back")
			tx.Rollback()
			panic(p)
		}
	}()

	result = tx_func(DaoFromTx(&NakshaTx{tx}))
	if result.IsSuccess() {
		err = tx.Commit()
		if err != nil {
			logger.LogMessage(fmt.Sprintf("TxTransaction.Commit. Error = %s", err))
			result = helper.ErrorResultString("Query failed: commit")
		}
	} else {
		tx.Rollback()
	}

	return result
}

// start_index: Postgresql users $1, $2 etc. as place holders.
// start_index will be an integer from which the place holders should start.
func (di *DaoImpl) partsParamsFromMap(values map[string]interface{}, start_index int) ([]string, []interface{}) {
	var params []interface{}
	var parts = make([]string, 0)

	for k, v := range values {
		parts = append(parts, fmt.Sprintf("%v = $%v", k, start_index))
		params = append(params, v)
		start_index++
	}

	return parts, params
}

func (di *DaoImpl) buildSelectQuery(sq SelectQuery) (string, []interface{}) {
	var query string
	var columns_str string
	var where_parts []string
	var params []interface{}
	var where_str string
	var order_by_str string
	var limit_str string
	var offset_str string

	columns_str = strings.Join(sq.Columns, ", ")
	if sq.Where == nil {
		where_str = ""
	} else {
		where_parts, params = di.partsParamsFromMap(sq.Where, 1)
		where_str = " WHERE " + strings.Join(where_parts, " AND ")
	}

	if len(sq.OrderBy) == 0 {
		order_by_str = ""
	} else {
		order_by_str = "ORDER BY " + sq.OrderBy
	}

	limit_str = ""
	offset_str = ""
	if sq.Limit > 0 {
		limit_str = fmt.Sprintf("LIMIT %v", sq.Limit)
		if sq.Offset > 0 {
			offset_str = fmt.Sprintf("OFFSET %v", sq.Offset)
		}
	}

	query = fmt.Sprintf("SELECT %v FROM %v %v %v %v %v", columns_str, sq.Table, where_str, order_by_str, limit_str, offset_str)
	query = strings.TrimSpace(query)

	return query, params
}

func (di *DaoImpl) makeErrorResult(err error) SelectResult {
	return SelectResult{err, nil, nil}
}
