package importer

import (
	"naksha/helper"
)

type emptyImporter struct {
	baseImporter
}

func (ei *emptyImporter) doImport() helper.Result {
	var result helper.Result

	err := ei.createTable()
	if err != nil {
		return helper.ErrorResultString("Could not create table")
	}
	ei.logNotice("created table " + ei.getTableName())

	err = ei.prepareTable()
	if err != nil {
		return helper.ErrorResultString("could not prepare table")
	}
	ei.updateStatus(READY)
	result = helper.MakeSuccessResult()

	return result
}

func (ei *emptyImporter) createTable() error {
	query := "CREATE TABLE " + ei.getTableName() + " ( " +
		"  " + PRIMARY_KEY + " SERIAL4 NOT NULL, " +
		"  the_geom Geometry(Geometry, 4326), " +
		"  PRIMARY KEY(" + PRIMARY_KEY + ") " +
		");"
	_, err := ei.map_user_dao.Exec(query)

	return err
}
