package fileutil

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func StoreFile(path string, name string, fileHeader *multipart.FileHeader) error {
	f, err := os.Create(path + name)
	if err != nil {
		return errors.New("failed to create file")
	}
	defer f.Close()
	file, err := fileHeader.Open()
	if err != nil {
		return errors.New("failed to read file")
	}
	defer file.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		return errors.New("failed to copy file")
	}
	return nil
}

func CopyFile(source string, dest string) error {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourcefile.Close()
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()
	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			return os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return nil
}

func CopyDir(source string, dest string) error {
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	// create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}
	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}
	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()
		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
