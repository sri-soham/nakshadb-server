package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/context"
	"log"
	"naksha"
	"naksha/db"
	"naksha/lib"
	"naksha/logger"
	"naksha/service"
	"net/http"
	"os"
	"time"
)

var router *lib.Router

func setup(config_file_path string) naksha.AppConfig {
	app_config, err := naksha.MakeConfig(config_file_path)
	if err != nil {
		log.Println("Could not parse config file. error = ", err)
		os.Exit(1)
	}

	map_user_dao, db_err := db.MapUserDao(app_config)
	if db_err != nil {
		log.Println("Failed to connect to database. error = ", db_err.Error())
		os.Exit(1)
	}

	api_user_dao, db_err := db.ApiUserDao(app_config)
	if db_err != nil {
		log.Println("Failed to connect to database. error = ", db_err.Error())
		os.Exit(1)
	}

	tmpl, err := lib.ParseTemplates(app_config.TemplatesDir())
	if err != nil {
		log.Println("Template parsing failed. ", err)
		os.Exit(1)
	}

	service_repository := service.GetServicesRepository(map_user_dao, api_user_dao)

	router = lib.MakeRouter(&app_config, tmpl, service_repository)

	t := time.Now()
	asset_timestamp := fmt.Sprintf("%d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
	lib.SetAssetTimestamp(asset_timestamp)

	err = logger.InitLoggers(app_config.LogsDir())
	if err != nil {
		log.Println("Could not open log file: ", err)
		os.Exit(1)
	}

	lib.InitSessionStore(
		app_config.SessionsDir()+"/",
		app_config.AuthKey(),
		app_config.EncKey())

	return app_config
}

type maxBytesHandler struct {
	h http.Handler
	n int64
}

func (m *maxBytesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, m.n)
	m.h.ServeHTTP(w, r)
}

func main() {
	defer logger.CloseLoggers()
	config_file := flag.String("config", "./config/config.xml", "Path to config file")
	flag.Parse()
	app_config := setup(*config_file)
	http.HandleFunc("/", router.ParseAndHandle)
	//http.ListenAndServe("10.0.3.230:8080", context.ClearHandler(http.DefaultServeMux))
	// This is to limit request body size to 8MB.
	// Based on https://stackoverflow.com/questions/28282370/is-it-advisable-to-further-limit-the-size-of-forms-when-using-golang
	http.ListenAndServe(app_config.HostPort(),
		context.ClearHandler(&maxBytesHandler{h: http.DefaultServeMux, n: app_config.MaxRequestBodySize()}))
}
