package main

import (
	"flag"
	"fmt"
	"naksha"
	"naksha/db"
	"naksha/logger"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func main() {

	sqls_path := flag.String("sqls_path", "./sqls", "Path to directory containing sql files")
	config_file := flag.String("config", "./config/config.xml", "Path to config file")
	flag.Parse()

	if _, err := os.Stat(*sqls_path); os.IsNotExist(err) {
		fmt.Println(*sqls_path + " does not exist")
		os.Exit(1)
	}
	if _, err := os.Stat(*config_file); os.IsNotExist(err) {
		fmt.Println(*config_file + " does not exist")
		os.Exit(1)
	}

	sqls_dir := getSqlsDir(*sqls_path)
	app_config, err := naksha.MakeConfig(*config_file)
	if err != nil {
		fmt.Println("could not parse config file. error = ", err)
		os.Exit(1)
	}

	logger.InitLoggers(app_config.LogsDir())
	main_user_dao, db_err := db.AdminUserDao(app_config)
	if db_err != nil {
		fmt.Println("Failed to connect to database. error = ", db_err.Error())
		os.Exit(1)
	}

	exists, err := migrationsTableExists(main_user_dao)
	if err != nil {
		fmt.Println("Could not check for migrations table. Error = ", err)
		os.Exit(1)
	}

	if !exists {
		err := createMigrationsTable(main_user_dao)
		if err != nil {
			fmt.Println("Could not create migrations table. Error = ", err)
			os.Exit(1)
		}
	}

	max_version, err := getMaxImportedVersion(main_user_dao)
	if err != nil {
		fmt.Println("Could not get max imported version. Error = ", err)
		os.Exit(1)
	}

	sql_files, err := getSqlFilesToBeImported(sqls_dir, max_version)
	if err != nil {
		fmt.Println("Could not get files to be imported. Error = ", err)
		os.Exit(1)
	}

	if len(sql_files) == 0 {
		fmt.Println("Database is up to date")
	} else {
		importFiles(main_user_dao, app_config, sqls_dir, sql_files)
	}
}

func getSqlsDir(path string) string {
	var sqls_dir string

	if strings.HasPrefix(path, "/") {
		sqls_dir = path
	} else {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error: could not get current working directory")
			os.Exit(2)
		}
		sqls_dir, err = filepath.Abs(dir + "/" + path)
		if err != nil {
			fmt.Println("Could not generate absolute path to sqls directory")
			os.Exit(2)
		}
	}

	_, err := os.Stat(sqls_dir)
	if err != nil {
		fmt.Println("Error: sqls_dir = ", sqls_dir, " does not exist")
		os.Exit(2)
	}

	return sqls_dir
}

func migrationsTableExists(main_user_dao db.Dao) (bool, error) {
	where := make(map[string]interface{})
	where["table_schema"] = "public"
	where["table_name"] = db.MasterMigration()
	cnt, err := main_user_dao.CountWhere("information_schema.tables", where)
	exists := (cnt == 1)

	return exists, err
}

func createMigrationsTable(main_user_dao db.Dao) error {
	query := "CREATE TABLE " + db.MasterMigration() + " (version VARCHAR(16), PRIMARY KEY(version));"
	_, err := main_user_dao.Exec(query)

	return err
}

func getMaxImportedVersion(main_user_dao db.Dao) (string, error) {
	query := "SELECT MAX(version) AS mversion FROM " + db.MasterMigration()
	params := make([]interface{}, 0)
	sres := main_user_dao.FetchData(query, params)
	if sres.Error != nil {
		return "0", sres.Error
	}

	return sres.Rows[0]["mversion"], nil
}

func getSqlFilesToBeImported(sqls_dir string, max_version string) ([]string, error) {
	sql_files := make([]string, 0)
	err := filepath.Walk(sqls_dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".sql") {
			parts := strings.Split(info.Name(), ".")
			if parts[0] > max_version {
				sql_files = append(sql_files, parts[0])
			}
		}

		return err
	})

	return sql_files, err
}

func listFiles(sqls_dir string, max_version string) error {
	err := filepath.Walk(sqls_dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".sql") {
			parts := strings.Split(info.Name(), ".")
			if parts[0] > max_version {
				fmt.Println(">  | path = ", path, ", name = ", info.Name())
			} else {
				fmt.Println("<= | path = ", path, ", name = ", info.Name())
			}
		}
		return err
	})

	return err
}

func importFiles(main_user_dao db.Dao, app_config naksha.AppConfig, sqls_dir string, sql_files []string) {
	sort.Strings(sql_files)
	conn_str := db.PsqlConnectionStringMain(app_config)
	cmd := "psql"
	for _, file := range sql_files {
		file_path := sqls_dir + "/" + file + ".sql"
		args := []string{conn_str, "-f", file_path}
		out, err := exec.Command(cmd, args...).CombinedOutput()
		out_str := string(out)
		out_str = strings.ToLower(out_str)
		has_error := false
		if strings.Index(out_str, "error") > -1 {
			has_error = true
		}
		if strings.Index(out_str, "warning") > -1 {
			has_error = true
		}
		if err != nil {
			fmt.Printf("Could not import %s. Output = %s. Error = %v\n", file_path, out_str, err)
			os.Exit(1)
		}
		if has_error {
			fmt.Printf("Could not import %s. Output = %s. Error = %v\n", file_path, out_str, err)
			os.Exit(1)
		}
		values := make(map[string]interface{})
		values["version"] = file
		_, err = main_user_dao.Insert(db.MasterMigration(), values, "")
		if err != nil {
			fmt.Printf("Imported %s but could update database. Error = %v\n", file_path, err)
			os.Exit(1)
		}
		fmt.Println("Imported ", file_path)
	}
}
