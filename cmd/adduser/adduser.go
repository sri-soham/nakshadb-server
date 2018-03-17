package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"naksha"
	"naksha/db"
	"naksha/helper"
	"naksha/logger"
	"os"
	"strings"
)

func main() {
	config_file := flag.String("config", "./config/config.xml", "'-config=\"path-to-config-file\"'")
	flag.Parse()
	if _, err := os.Stat(*config_file); os.IsNotExist(err) {
		fmt.Println(*config_file + " does not exist")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	name := get_input_max_len(scanner, "Name", 64)
	username := get_input_max_len(scanner, "Desired username", 32)
	password := get_input_min_len(scanner, "Password", 8)

	app_config, err := naksha.MakeConfig(*config_file)
	if err != nil {
		fmt.Println("Could not parse config file. Error = ", err)
		os.Exit(1)
	}

	logger.InitLoggers(app_config.LogsDir())

	main_user_dao, err := db.AdminUserDao(app_config)
	if err != nil {
		fmt.Println("Failed to connect to database. Error = ", err)
		os.Exit(1)
	}

	password_byte := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(password_byte, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Could not generate hash. Error = ", err)
		os.Exit(1)
	}

	where := make(map[string]interface{})
	where["username"] = username
	count, err := main_user_dao.CountWhere(db.MasterUser(), where)
	if err != nil {
		fmt.Println("could not get count")
		os.Exit(1)
	}
	if count >= 1 {
		fmt.Println("Username " + username + " is in use")
		os.Exit(1)
	}

	values := make(map[string]interface{})
	values["name"] = name
	values["username"] = username
	values["password"] = hash
	schema_name := helper.RandomSchemaName(8)
	values["schema_name"] = schema_name
	result := main_user_dao.TxTransaction(func(dao db.Dao) helper.Result {
		stmt := "CREATE SCHEMA " + schema_name
		_, err := dao.Exec(stmt)
		if err != nil {
			return helper.ErrorResultString("could not create schema")
		}

		stmt = "GRANT ALL PRIVILEGES ON SCHEMA " + schema_name + " TO " + app_config.DBUser()
		_, err = dao.Exec(stmt)
		if err != nil {
			return helper.ErrorResultString("could grant privileges to map user")
		}

		stmt = "GRANT USAGE ON SCHEMA " + schema_name + " TO " + app_config.DBApiUser()
		_, err = dao.Exec(stmt)
		if err != nil {
			return helper.ErrorResultString("Could not grant privileges to api user")
		}

		_, err = main_user_dao.Insert(db.MasterUser(), values, "")
		if err != nil {
			return helper.ErrorResultString("Could not add user")
		} else {
			result := helper.MakeSuccessResult()
			result.AddToData("password", password)
			return result
		}
	})
	if result.IsSuccess() {
		fmt.Println("User added.")
	} else {
		errs := result.GetErrors()
		fmt.Println("Error: ", strings.Join(errs, "\n"))
	}
}

func get_input_max_len(scanner *bufio.Scanner, input_name string, max_len int) string {
	input_text := ""

	for input_text == "" {
		fmt.Printf("Enter "+input_name+" (Max. length: %d): ", max_len)
		scanner.Scan()
		input_text = scanner.Text()
		input_text = strings.TrimSpace(input_text)
		input_len := len(input_text)
		if input_len == 0 {
			fmt.Println("Invalid value")
		} else if input_len > max_len {
			fmt.Println("Max. allowed length is ", max_len)
			input_text = ""
		}
	}

	return input_text
}

func get_input_min_len(scanner *bufio.Scanner, input_name string, min_len int) string {
	input_text := ""

	for input_text == "" {
		fmt.Printf("Enter "+input_name+" (Min. length: %d): ", min_len)
		scanner.Scan()
		input_text = scanner.Text()
		input_text = strings.TrimSpace(input_text)
		input_len := len(input_text)
		if input_len == 0 {
			fmt.Println("Invalid value")
		} else if input_len < min_len {
			fmt.Println(input_name, " should have at least ", min_len, " characters")
			input_text = ""
		}
	}

	return input_text
}
