package importer

import (
	"fmt"
	"naksha/logger"
	"os/exec"
)

type ogr2OgrImporter struct {
	baseImporter
	db_conn_str string
}

func (oi *ogr2OgrImporter) doImport() {
	var err error

	err = oi.toPostgreSQL()
	if err != nil {
		oi.updateStatus(ERROR)
		return
	}
	oi.updateStatus(IMPORTED)
	oi.logNotice("imported data with ogr2ogr")

	err = oi.prepareTable()
	if err != nil {
		oi.logError("could not run stored proc.", err)
		oi.updateStatus(ERROR)
		return
	}
	oi.logNotice("table prepared by running stored proc")
	oi.updateStatus(READY)

	err = oi.deleteDirectory()
	if err != nil {
		oi.logError("Could not delete import workspace directory.", err)
	} else {
		oi.logNotice("Deleted workspace directory containing " + oi.getOutfilePath())
	}
}

func (oi *ogr2OgrImporter) toPostgreSQL() error {
	file_path := oi.getOutfilePath()
	cmd := "ogr2ogr"
	args := []string{"-f", "PostgreSQL", oi.db_conn_str, file_path, "-nln", oi.getTableName(),
		"-nlt", "PROMOTE_TO_MULTI", "-lco", "GEOMETRY_NAME=the_geom",
		"-lco", "FID=" + PRIMARY_KEY, "-t_srs", "EPSG:4326"}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("%v. output = %s. error = %v", file_path, out, err)
		logger.ImporterLog(msg)
	}

	return err
}
