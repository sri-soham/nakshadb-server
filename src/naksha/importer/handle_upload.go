package importer

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"naksha/helper"
	"naksha/logger"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type UploadedFile struct {
	IsValid      bool
	FileType     int
	OutfilePath  string
	ErrorMessage string
}

type importFileHandler struct {
	req            *http.Request
	workspace_dir  string
	directory_name string
	file_name      string
}

func (ifh *importFileHandler) handle() UploadedFile {
	var uf UploadedFile
	var err error

	err = ifh.copyUploadedFile()
	if err != nil {
		return UploadedFile{false, INVALID_FILE, "", err.Error()}
	}

	ext := strings.ToLower(path.Ext(ifh.file_name))
	if ext == ".zip" {
		err := ifh.processZipFile()
		if err != nil {
			return UploadedFile{false, INVALID_FILE, "", err.Error()}
		}

		ext = strings.ToLower(path.Ext(ifh.file_name))
	}

	switch ext {
	case ".csv":
		uf = UploadedFile{true, CSV_FILE, ifh.fullPath(), ""}
	case ".json", ".geojson":
		uf = UploadedFile{true, GEOJSON_FILE, ifh.fullPath(), ""}
	case ".shp":
		uf = UploadedFile{true, SHAPE_FILE, ifh.fullPath(), ""}
	case ".kml":
		uf = UploadedFile{true, KML_FILE, ifh.fullPath(), ""}
	default:
		uf = UploadedFile{false, -10, "", "Invalid file format"}
	}

	return uf
}

func (ifh *importFileHandler) createDirectory() error {
	ifh.directory_name = fmt.Sprintf("%s/%s_%v", ifh.workspace_dir, helper.RandomString(16), time.Now().Unix())
	err := os.Mkdir(ifh.directory_name, 0775)
	if err != nil {
		msg := fmt.Sprintf("Could not create directory - %s . %s", ifh.directory_name, err)
		logger.ImporterLog(msg)
		return errors.New("Could not create directory")
	}

	return nil
}

func (ifh *importFileHandler) copyUploadedFile() error {
	ifh.req.ParseMultipartForm(4194304) // 4MB
	infile, handler, err := ifh.req.FormFile("file")
	if err != nil {
		return err
	}
	defer infile.Close()

	err = ifh.createDirectory()
	if err != nil {
		return err
	}

	// 2017-11-10
	// If request is submitted with "mime/multipart"
	//   var b bytes.Buffer
	//   w := multipart.NewWriter(&b)
	//   fw, err := w.CreateFormFile("file", "../../../bus_stop.csv")
	// then ifh.file_name = <directory> + "../../../bus_stop.csv"
	// Came across this during testing.
	// Just pick up the file name from handler.Filename to be on safe side.
	ifh.file_name = path.Base(handler.Filename)
	outfile_path := ifh.fullPath()
	outfile, err := os.OpenFile(outfile_path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		msg := fmt.Sprintf("Could not open destination file - %s . %s", outfile_path, err)
		logger.ImporterLog(msg)
		return errors.New("Could not copy file")
	}
	defer outfile.Close()
	io.Copy(outfile, infile)
	logger.ImporterLog("Uploaded file: " + handler.Filename)

	return nil
}

func (ifh *importFileHandler) processZipFile() error {
	err := ifh.unzipFile()
	if err != nil {
		return errors.New("Could not unzip file")
	} else {
		return ifh.decipherUploadedFileType()
	}
}

func (ifh *importFileHandler) unzipFile() error {
	cmd := "unzip"
	outfile_path := ifh.fullPath()
	args := []string{outfile_path, "-d", ifh.directory_name}
	err := exec.Command(cmd, args...).Run()
	if err != nil {
		logger.ImporterLog("Unzipping failed: " + outfile_path + ". " + err.Error())
	}

	return err
}

func (ifh *importFileHandler) decipherUploadedFileType() error {
	var err error

	files_list, err := ioutil.ReadDir(ifh.directory_name)
	if err != nil {
		return errors.New("Could not get files list")
	}

	// Even for shape files, zip file should contain at most 5 files.
	no_of_files := len(files_list)
	if no_of_files > 6 {
		return errors.New("Invalid file")
	}

	if no_of_files == 2 {
		// zip file and the unzipped file
		for _, f := range files_list {
			filename := f.Name()
			ext := strings.ToLower(path.Ext(filename))
			if ext != ".zip" {
				ifh.file_name = filename
			}
		}
		err = nil
	} else {
		// more than one file, must be shape file
		valid_count := 0
		for _, f := range files_list {
			filename := f.Name()
			ext := strings.ToLower(path.Ext(filename))
			switch ext {
			case ".shp":
				valid_count++
				ifh.file_name = filename
			case ".shx":
				valid_count++
			case ".dbf":
				valid_count++
			case ".prj":
				valid_count++
			}
		}
		if valid_count == 4 {
			err = nil
		} else {
			err = errors.New("One of .shp, .shx, .prj, .dbf files is missing")
		}
	}

	return err
}

func (ifh *importFileHandler) fullPath() string {
	return ifh.directory_name + "/" + ifh.file_name
}
