package main

import (
	"flag"
	"fmt"
	"naksha"
	"naksha/db"
	"os"
)

func main() {
	config_file := flag.String("config", "", "'-config=\"path-to-config-file\"'")
	days_before := flag.Int("days", 3, "'-days=\"delete-exports-older-than-x-days\"'")
	flag.Parse()
	if len(*config_file) == 0 {
		fmt.Println("Please enter value for config")
		os.Exit(1)
	}
	if *days_before < 2 {
		fmt.Println("Only the exports older than 2 days can be deleted")
		os.Exit(1)
	}

	app_config, err := naksha.MakeConfig(*config_file)
	if err != nil {
		fmt.Println("Could not parse config file. Error = ", err)
		os.Exit(1)
	}

	map_user_dao, err := db.MapUserDao(app_config)
	if err != nil {
		fmt.Println("Failed to connect to database. Error = ", err)
		os.Exit(1)
	}

	old_exports, err := getExportsOlderThan3Days(map_user_dao, *days_before)
	if err != nil {
		fmt.Println("Could not get exports list. Error = ", err)
		os.Exit(1)
	}

	if len(old_exports) > 0 {
		deleteExports(app_config, map_user_dao, old_exports)
	}
}

func getExportsOlderThan3Days(map_user_dao db.Dao, days_before int) ([]map[string]string, error) {
	query := fmt.Sprintf("SELECT * FROM %v WHERE updated_at <= (CURRENT_DATE - INTERVAL '%v day')", db.MasterExport(), days_before)
	sres := map_user_dao.FetchData(query, []interface{}{})
	if sres.Error != nil {
		return nil, sres.Error
	}

	return sres.Rows, nil
}

func deleteExports(app_config naksha.AppConfig, map_user_dao db.Dao, old_exports []map[string]string) {
	for _, export := range old_exports {
		path := fmt.Sprintf("%v/%v", app_config.ExportsDir(), export["hash"])
		where := map[string]interface{}{"id": export["id"]}
		err := map_user_dao.Delete(db.MasterExport(), where)
		if err == nil {
			err := os.RemoveAll(path)
			if err == nil {
				fmt.Println("Deleted: Path = ", path)
			} else {
				fmt.Printf("Error: path = %v. Error = %v\n", path, err)
			}
		} else {
			fmt.Println("Could not delete export from db. Error = ", err)
		}
	}
}
