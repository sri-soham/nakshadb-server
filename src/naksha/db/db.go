package db

import (
	"database/sql"
	"errors"
	"fmt"
	"naksha"
)

const (
	MSTR_USER      = "mstr_user"
	MSTR_TABLE     = "mstr_table"
	MSTR_LAYER     = "mstr_layer"
	MSTR_MAP       = "mstr_map"
	MSTR_MAP_LAYER = "mstr_map_layer"
	MSTR_EXPORT    = "mstr_export"
	MSTR_MIGRATION = "mstr_migration"
)

type SelectResult struct {
	Error   error
	Columns []string
	Rows    []map[string]string
}

type SelectQuery struct {
	Table   string
	Columns []string
	Where   map[string]interface{}
	OrderBy string
	Limit   int
	Offset  int
}

type NakshaTx struct {
	*sql.Tx
}

func (nt *NakshaTx) Begin() (*sql.Tx, error) {
	return nil, errors.New("Begin not allowed on Tx")
}

type Queryable interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Begin() (*sql.Tx, error)
}

func InitDBConnection(db_user string, db_pass string, db_name string, ssl_mode string) (*sql.DB, error) {
	conn_string := "user=" + db_user + " dbname=" + db_name + " password=" + db_pass + " sslmode=" + ssl_mode
	conn, err := sql.Open("postgres", conn_string)

	return conn, err
}

func AdminUserDao(app_config naksha.AppConfig) (Dao, error) {
	conn, err := InitDBConnection(app_config.DBAdminUser(), app_config.DBAdminPass(), app_config.DBName(), app_config.DBSslMode())
	if err == nil {
		return &DaoImpl{conn, true}, err
	} else {
		return nil, err
	}
}

func MapUserDao(app_config naksha.AppConfig) (Dao, error) {
	conn, err := InitDBConnection(app_config.DBUser(), app_config.DBPass(), app_config.DBName(), app_config.DBSslMode())
	if err == nil {
		return &DaoImpl{conn, true}, err
	} else {
		return nil, err
	}
}

func ApiUserDao(app_config naksha.AppConfig) (Dao, error) {
	conn, err := InitDBConnection(app_config.DBApiUser(), app_config.DBApiPass(), app_config.DBName(), app_config.DBSslMode())
	if err == nil {
		return &DaoImpl{conn, true}, err
	} else {
		return nil, err
	}
}

func DaoFromConn(conn *sql.DB) Dao {
	return &DaoImpl{conn, true}
}

func DaoFromTx(tx *NakshaTx) Dao {
	return &DaoImpl{tx, false}
}

func MasterTable() string {
	return "public." + MSTR_TABLE
}

func MasterUser() string {
	return "public." + MSTR_USER
}

func MasterLayer() string {
	return "public." + MSTR_LAYER
}

func MasterMap() string {
	return "public." + MSTR_MAP
}

func MasterMapLayer() string {
	return "public." + MSTR_MAP_LAYER
}

func MasterExport() string {
	return "public." + MSTR_EXPORT
}

func MasterMigration() string {
	return "public." + MSTR_MIGRATION
}

func PsqlConnectionString(app_config naksha.AppConfig) string {
	return fmt.Sprintf("postgresql://%v:%v@%v/%v?sslmode=%v",
		app_config.DBUser(),
		app_config.DBPass(),
		app_config.DBHost(),
		app_config.DBName(),
		app_config.DBSslMode(),
	)
}

func PsqlConnectionStringMain(app_config naksha.AppConfig) string {
	return fmt.Sprintf("postgresql://%v:%v@%v/%v?sslmode=%v",
		app_config.DBAdminUser(),
		app_config.DBAdminPass(),
		app_config.DBHost(),
		app_config.DBName(),
		app_config.DBSslMode(),
	)
}

func Ogr2ogrConnectionString(app_config naksha.AppConfig) string {
	return fmt.Sprintf("PG:host=%v user=%v dbname=%v password=%v",
		app_config.DBHost(), app_config.DBUser(), app_config.DBName(), app_config.DBPass())
}
