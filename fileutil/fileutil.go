package fileutil

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/alexandervantrijffel/goutil/errorcheck"
	"github.com/alexandervantrijffel/goutil/logging"
	"github.com/alexandervantrijffel/goutil/stringutil"
)

func CopyFile(source string, dest string) error {
	var err error
	defer func() {
		err = errorcheck.CheckLogf(err, "Failed to copy file from '%s' to '%s'", source, dest)
	}()
	sourcefile, err := os.Open(source)
	if err != nil {
		return errors.Wrapf(err, "Failed to open source")
	}
	defer sourcefile.Close()
	destfile, err := os.Create(dest)
	if err != nil {
		return errors.Wrapf(err, "Failed to create destination")
	}
	defer destfile.Close()
	if _, err = io.Copy(destfile, sourcefile); err != nil {
		err = errors.Wrapf(err, "Failed to copy")
	}
	if sourceinfo, err := os.Stat(source); err == nil {
		return os.Chmod(dest, sourceinfo.Mode())
	}
	return nil // ignore setting chmod error
}

func CopyDir(source string, dest string) error {
	return errorcheck.CheckLogf(copyDir(source, dest), "Failed to copy dir '%s' to '%s'", source, dest)
}
func copyDir(source string, dest string) error {
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return errors.Wrap(err, "Failed to get properties of source dir")
	}
	if err = os.MkdirAll(dest, sourceinfo.Mode()); err != nil {
		return errors.Wrap(err, "Failed to create destination file")
	}
	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)
	if err != nil {
		return errors.Wrap(err, "Failed to read from source directory")
	}
	for _, obj := range objects {
		currentSource := path.Join(source, obj.Name())
		currentDestination := path.Join(dest, obj.Name())
		if obj.IsDir() {
			if err = CopyDir(currentSource, currentDestination); err != nil {
				return errors.Wrapf(err, "Failure while copying subfolder '%s' to '%s'", currentSource, currentDestination)
			}
			continue
		}
		if err = CopyFile(currentSource, currentDestination); err != nil {
			return errors.Wrapf(err, "Failure while copying file '%s' to '%s'", currentSource, currentDestination)
		}
	}
	return nil
}

func StoreMultipartFile(folder string, name string, fileHeader *multipart.FileHeader) (err error) {
	defer func() {
		err = errorcheck.CheckLogf(err, "Failed to store file %s/%s", folder, name)
	}()
	if err = os.MkdirAll(folder, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create folder")
	}
	filePath := path.Join(folder, name)
	f, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer f.Close()
	file, err := fileHeader.Open()
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}
	defer file.Close()
	if _, err = io.Copy(f, file); err != nil {
		return errors.Wrap(err, "failed to copy file")
	}
	logging.Debugf("Stored file %s", filePath)
	return
}
func RandomizeFileName(org string) string {
	fileName := path.Base(org)
	folder := path.Dir(org)
	rand := stringutil.RandomAlphanumericString(8)
	dot := strings.Index(fileName, ".")
	if dot == -1 {
		return fileName + "-" + rand
	}
	return path.Join(folder, fileName[0:dot]+"-"+rand+fileName[dot:])
}
