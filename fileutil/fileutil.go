package fileutil

import (
	"io"
	"mime/multipart"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/alexandervantrijffel/goutil/errorcheck"
	"github.com/alexandervantrijffel/goutil/logging"
)

func CopyFile(source string, dest string) (err error) {
	defer func() {
		err = errorcheck.CheckLogf(err, "Failed to copy file from '%s' to '%s'", source, dest)
	}()
	sourcefile, err := os.Open(source)
	if err != nil {
		return errors.Wrapf(err, "Failed to open '%s'", source)
	}
	defer sourcefile.Close()
	destfile, err := os.Create(dest)
	if err != nil {
		return errors.Wrapf(err, "Failed to create '%s'", dest)
	}
	defer destfile.Close()
	_, err = io.Copy(destfile, sourcefile)
	if err != nil {
		return errors.Wrapf(err, "Failed to copy")
	}
	sourceinfo, err := os.Stat(source)
	if err == nil {
		return os.Chmod(dest, sourceinfo.Mode())
	}
	return nil // ignore setting chmod error
}

func CopyDir(source string, dest string) (err error) {
	defer func() {
		err = errorcheck.CheckLogf(err, "Failed to copy dir '%s' to '%s'", source, dest)
	}()
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return errors.Wrap(err, "Failed to get properties of source dir")
	}
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
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
			err = CopyDir(currentSource, currentDestination)
			if err != nil {
				return errors.Wrapf(err, "Failure while copying subfolder '%s' to '%s'", currentSource, currentDestination)
			}
			continue
		}
		err = CopyFile(currentSource, currentDestination)
		if err != nil {
			if err != nil {
				return errors.Wrapf(err, "Failure while copying file '%s' to '%s'", currentSource, currentDestination)
			}
		}
	}
	return
}

func StoreMultipartFile(folder string, name string, fileHeader *multipart.FileHeader) (err error) {
	defer func() {
		err = errorcheck.CheckLogf(err, "Failed to store file %s/%s", folder, name)
	}()

	filePath := path.Join(folder, name)
	f, err := os.Create(filePath)
	if err != nil {
		err = errors.Wrap(err, "failed to create file")
	}
	defer f.Close()
	file, err := fileHeader.Open()
	if err != nil {
		err = errors.Wrap(err, "failed to read file")
	}
	defer file.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		return errors.Wrap(err, "failed to copy file")
	}
	logging.Debugf("Stored file %s", filePath)
	return
}
