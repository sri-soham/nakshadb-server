package exporter_test

import (
	"flag"
	"fmt"
	"naksha"
	"naksha/db"
	"naksha/exporter"
	"naksha/importer"
	"naksha/logger"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

var app_config naksha.AppConfig
var map_user_dao db.Dao
var sqlfile *string
var table_name = "tbl_sample_export"
var table_id string
var user_id = 1

func init() {
	var err error

	config_path := flag.String("config_path", "", "Full path to config file")
	sqlfile = flag.String("sqlfile", "", "Path to sql file to be imported")
	flag.Parse()
	err_count := 0
	if len(*config_path) == 0 {
		fmt.Println("Please pass -config_path command line parameter")
		err_count++
	}
	if len(*sqlfile) == 0 {
		fmt.Println("Please pass -sqlfile command line parameter")
		err_count++
	}
	if _, err := os.Stat(*config_path); err != nil {
		fmt.Println("Config file ", *config_path, " does not exist")
		err_count++
	}
	if _, err := os.Stat(*sqlfile); err != nil {
		fmt.Println("Sql file ", *sqlfile, " does not exist")
		err_count++
	}
	if err_count > 0 {
		os.Exit(1)
	}

	app_config, err = naksha.MakeConfig(*config_path)
	if err != nil {
		fmt.Println("Could not parse config file. error = ", err)
		os.Exit(1)
	}

	map_user_dao, err = db.MapUserDao(app_config)

	err = logger.InitLoggers("../../../test")
	if err != nil {
		fmt.Println("Could not open log file: ", err)
		os.Exit(1)
	}

	table_id, err = testImportSqlFile()
	if err != nil {
		fmt.Println("Error = ", err)
		os.Exit(1)
	}
}

func testImportSqlFile() (string, error) {
	var table_id string

	cmd := "psql"
	conn_str := db.PsqlConnectionString(app_config)
	args := []string{conn_str, "-f", *sqlfile}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		fmt.Printf("%s. Output = %s. Error = %s\n", *sqlfile, out, err)
		return "", err
	}

	where := map[string]interface{}{"table_name": table_name}
	sres := map_user_dao.SelectWhere(db.MSTR_TABLE, []string{"*"}, where)
	if sres.Error != nil {
		fmt.Println("could not select from ", db.MSTR_TABLE, ". Error = ", sres.Error)
		return "", sres.Error
	}
	if len(sres.Rows) == 0 {
		values := make(map[string]interface{})
		values["user_id"] = user_id
		values["name"] = "Sample Export"
		values["table_name"] = table_name
		values["schema_name"] = "public"
		values["status"] = importer.READY
		tmp, err := map_user_dao.Insert(db.MSTR_TABLE, values, "id")
		table_id = fmt.Sprintf("%d", tmp)
		if err != nil {
			fmt.Println("Could not add record to ", db.MSTR_TABLE, ". Error = ", db.MSTR_TABLE, err)
			return "", err
		}
	} else {
		table_id = sres.Rows[0]["id"]
	}

	return table_id, nil
}

func TestIsValidFormat(t *testing.T) {
	var is_valid bool

	is_valid = exporter.IsValidFormat(fmt.Sprintf("%d", exporter.SHAPE_FILE))
	if !is_valid {
		t.Errorf("Valid shape file format not recognized")
	}

	is_valid = exporter.IsValidFormat(fmt.Sprintf("%d", exporter.CSV_FILE))
	if !is_valid {
		t.Errorf("Valid CSV file format not recognized")
	}

	is_valid = exporter.IsValidFormat(fmt.Sprintf("%d", exporter.GEOJSON_FILE))
	if !is_valid {
		t.Errorf("Valid GeoJSON file format not recognized")
	}

	is_valid = exporter.IsValidFormat(fmt.Sprintf("%d", exporter.KML_FILE))
	if !is_valid {
		t.Errorf("Valid KML file format not recognized")
	}

	is_valid = exporter.IsValidFormat("0")
	if is_valid {
		t.Errorf("Invalid format (0) is being considered valid")
	}
}

func TestExportAsCsv(t *testing.T) {
	testExportAsFormat(t, fmt.Sprintf("%d", exporter.SHAPE_FILE), "Shape File", 4)
	testExportAsFormat(t, fmt.Sprintf("%d", exporter.CSV_FILE), "CSV File", 2)
	testExportAsFormat(t, fmt.Sprintf("%d", exporter.GEOJSON_FILE), "GeoJSON File", 2)
	testExportAsFormat(t, fmt.Sprintf("%d", exporter.KML_FILE), "KML File", 2)
}

// Shape files are taking longer to export.
func testExportAsFormat(t *testing.T, format_str string, format_desc string, wait_time time.Duration) {
	export_id, err := exporter.Export(&app_config, map_user_dao, user_id, table_id, format_str)
	if err != nil {
		t.Errorf("%s: Export failed. Error = %s", format_desc, err)
	} else {
		// wait for a second for the go routine to execute
		time.Sleep(wait_time * time.Second)
		export, err := map_user_dao.Find(db.MasterExport(), export_id)
		if err != nil {
			t.Errorf("%s: Could not fetch export details. Err = %v", format_desc, err)
		} else {
			status, ok := strconv.Atoi(export["status"])
			if ok != nil {
				t.Errorf("%s: Invalid value for status: %d", format_desc, status)
			} else {
				switch status {
				case exporter.ST_ERROR:
					t.Errorf("%s: Export failed", format_desc)
				case exporter.ST_IN_QUEUE:
					t.Errorf("%s: Export is still in queue", format_desc)
				}
			}
		}
	}
}
