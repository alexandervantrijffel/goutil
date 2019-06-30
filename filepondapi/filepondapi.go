package filepondapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/alexandervantrijffel/goutil/errorcheck"
	"github.com/alexandervantrijffel/goutil/fileutil"
	"github.com/alexandervantrijffel/goutil/logging"
	"github.com/alexandervantrijffel/goutil/webserviceutil"
	"github.com/julienschmidt/httprouter"
)

type Config struct {
	UploadBasePath string
	FinalBasePath  string
	// const _200MB = (1 << 20) * 200
	MaxPostFormMemoryBytes int64
	// config.THECONFIG.UPLOADBASEURL+
	BaseUrl string
}

var theConfig Config

func SetConfig(c Config) {
	theConfig = c
}
func configIsSet() bool {
	if len(theConfig.UploadBasePath) == 0 {
		logging.Error("filepondapi.Config not initialized. Please call filepondapi.SetConfig() first.")
		return false
	}
	return true
}

// based on https://github.com/AlexandrDobrovolskiy/storage/blob/master/controllers/files.go

func FilePondProcess(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !configIsSet() {
		return
	}
	err := r.ParseMultipartForm(theConfig.MaxPostFormMemoryBytes)
	if err != nil {
		_ = errorcheck.CheckLogf(err, "Failed to parse multipart form")
		webserviceutil.ReturnStatusMessageResponse(w, http.StatusBadRequest, "No files to upload")
		return
	}
	defer func() {
		_ = r.MultipartForm.RemoveAll()
	}()

	var filesList []string
	for _, files := range r.MultipartForm.File {
		for _, file := range files {
			// warning: overwrites existing files!
			err = fileutil.StoreMultipartFile(theConfig.UploadBasePath, file.Filename, file)
			if err != nil {
				webserviceutil.ReturnStatusMessageResponse(w, http.StatusInternalServerError,
					fmt.Sprintf("Failed to store file %s", file.Filename))
				return
			}
			filesList = append(filesList, path.Join(theConfig.BaseUrl, file.Filename))
		}
	}
	resp := map[string]interface{}{"description": "Operation succeeded"}
	resp["files"] = filesList
	webserviceutil.Return(w, http.StatusOK, resp)
}

func FilePondDelete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !configIsSet() {
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_ = errorcheck.CheckLogf(err, "Failed to ready body")
		webserviceutil.ReturnText(w, "Delete failed")
		return
	}
	toDelete := string(body)
	_ = errorcheck.CheckLogf(os.RemoveAll(path.Join(theConfig.UploadBasePath, toDelete)), "RemoveAll failed")
	webserviceutil.ReturnText(w, "Deleted")
}

func FilePondLoad(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	load := r.URL.Query()["load"]
	//use load[0] to get query
	dir := path.Join(theConfig.UploadBasePath, load[0])
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		_ = errorcheck.CheckLogf(err, "Failed to read dir '%s'", dir)
		webserviceutil.ReturnStatusMessageResponse(w, http.StatusNotFound, "Cannot file the specified directory")
	}
	_ = webserviceutil.ReturnFile(w, r, path.Join(dir, files[0].Name()))
}

type ConfirmMessage struct {
	Files []string `json:"files"`
}

func FilePondProcessSubmitForm(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	req := &ConfirmMessage{}
	var filesList []string
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		_ = errorcheck.CheckLogf(err, "Failed to decode request body")
		webserviceutil.ReturnStatusMessageResponse(w, http.StatusBadRequest, "Could not decode request body")
		return
	}
	for _, id := range req.Files {
		src := path.Join(theConfig.UploadBasePath, id)
		dest := path.Join(theConfig.FinalBasePath, id)
		err = os.MkdirAll(dest, os.ModePerm)
		if errorcheck.CheckLogf(err, "Failed to create destination directory %s", dest) != nil {
			return
		}
		err = fileutil.CopyDir(src, dest)
		if err != nil {
			_ = errorcheck.CheckLogf(err, "Error while copying directory")
			webserviceutil.ReturnStatusMessageResponse(w, http.StatusInternalServerError, "Error while copying directory")
			return
		}
		directory, _ := os.Open(dest)
		files, err := directory.Readdir(-1)
		if err != nil {
			_ = errorcheck.CheckLogf(err, "Failed to open destination folder")
			webserviceutil.ReturnStatusMessageResponse(w, http.StatusInternalServerError, "Failed to open destination folder")
			return
		}
		filesList = append(filesList, theConfig.BaseUrl+"/"+id+"/"+files[0].Name())
		_ = errorcheck.CheckLogf(os.RemoveAll(src), "RemoveAll failed")
	}
	resp := map[string]interface{}{"description": "Operation succeeded"}
	resp["files"] = filesList
	webserviceutil.Return(w, http.StatusOK, resp)
}
