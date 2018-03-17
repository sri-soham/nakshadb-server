package importer

import (
	"encoding/csv"
	"fmt"
	"naksha/logger"
	"os"
	"os/exec"
	"strings"
)

type csvFileImporter struct {
	baseImporter
	db_conn_str string
}

func (cfi *csvFileImporter) doImport() {
	var err error

	outfile_path := cfi.getOutfilePath()
	columns, err := cfi.getColumnNames(outfile_path)
	if err != nil {
		cfi.updateStatus(ERROR)
		return
	}
	cfi.logNotice("got columns")

	err = cfi.createTable(columns)
	if err != nil {
		cfi.logError("Could not create table", err)
		cfi.updateStatus(ERROR)
		return
	}
	cfi.logNotice("created table")

	err = cfi.importData(outfile_path)
	if err != nil {
		cfi.updateStatus(ERROR)
		return
	}
	cfi.updateStatus(IMPORTED)

	err = cfi.prepareCsv()
	if err != nil {
		cfi.logError("could not run stored proc naksha_prepare_csv.", err)
		cfi.updateStatus(ERROR)
		return
	}
	cfi.logNotice("prepared csv table")

	err = cfi.prepareTable()
	if err != nil {
		cfi.logError("could not run stored proc.", err)
		cfi.updateStatus(ERROR)
		return
	}
	cfi.logNotice("table prepared by running stored proc")
	cfi.updateStatus(READY)

	err = cfi.deleteDirectory()
	if err != nil {
		cfi.logError("Could not delete import workspace directory.", err)
	} else {
		cfi.logNotice("Deleted workspace directory containing " + cfi.getOutfilePath())
	}
}

func (cfi *csvFileImporter) getColumnNames(outfile_path string) ([]string, error) {
	columns := make([]string, 0)
	file, err := os.Open(outfile_path)
	if err != nil {
		cfi.logError("Could not open file to read header column", err)
		return columns, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	tmp, err := reader.Read()
	if err != nil {
		cfi.logError("Could not read headers", err)
		return columns, err
	}
	has_the_geom := false
	for i := range tmp {
		if tmp[i] == "the_geom" {
			has_the_geom = true
			break
		}
	}
	columns = append(columns, tmp...)
	if !has_the_geom {
		columns = append(columns, "the_geom")
	}

	return columns, nil
}

func (cfi *csvFileImporter) createTable(columns []string) error {
	query := "CREATE TABLE " + cfi.getTableName() + " (\n"
	parts := make([]string, 0)
	has_the_geom := false
	for i := range columns {
		switch columns[i] {
		case "the_geom":
			parts = append(parts, "the_geom Geometry(Geometry, 4326)")
			has_the_geom = true
		case "naksha_id":
			parts = append(parts, "naksha_id INTEGER NOT NULL")
		default:
			parts = append(parts, columns[i]+" TEXT")
		}
	}
	if !has_the_geom {
		parts = append(parts, "the_geom Geometry(Geometry, 4326)")
	}
	query += strings.Join(parts, ",\n  ")
	query += ");"
	_, err := cfi.map_user_dao.Exec(query)

	return err
}

func (cfi *csvFileImporter) importData(file_path string) error {
	psql_copy := "\\COPY " + cfi.getTableName() + " FROM '" + file_path + "' WITH CSV HEADER DELIMITER AS ','"
	cmd := "psql"
	args := []string{cfi.db_conn_str, "-c", psql_copy}

	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("%v. output = %s. error = %v", file_path, out, err)
		logger.ImporterLog(msg)
	}

	return err
}

func (cfi *csvFileImporter) prepareCsv() error {
	parts := strings.Split(cfi.getTableName(), ".")
	query := "SELECT public.naksha_prepare_csv($1, $2)"
	params := []interface{}{parts[0], parts[1]}
	return cfi.map_user_dao.StoredProcNoResult(query, params)
}
