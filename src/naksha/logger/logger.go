package logger

import (
	"log"
	"os"
)

var generic_file *os.File
var importer_file *os.File
var generic *log.Logger
var importer *log.Logger

func InitLoggers(logs_dir string) error {
	var err error

	file_path := logs_dir + "/logfile.log"
	generic_file, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return err
	}
	generic = log.New(generic_file, "", log.LstdFlags|log.Llongfile)

	file_path = logs_dir + "/importer.log"
	importer_file, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		generic_file.Close()
	} else {
		importer = log.New(importer_file, "", log.LstdFlags|log.Llongfile)
	}

	return err
}

func LogMessage(msg interface{}) {
	generic.Println(msg)
}

func ImporterLog(msg interface{}) {
	importer.Println(msg)
}

func CloseLoggers() {
	generic_file.Sync()
	generic_file.Close()

	importer_file.Sync()
	importer_file.Close()
}
