package errorcheck

import (
	"fmt"

	"github.com/alexandervantrijffel/goutil/logging"
	"github.com/pkg/errors"
)

func CheckPanic(e error) {
	if e != nil {
		logging.Error(e)
		panic(e)
	}
}

func CheckPanicWrap(e error, action string, v ...interface{}) {
	if e != nil {
		newE := errors.Wrapf(e, action, v...)
		logging.Error(newE)
		panic(newE)
	}
}

// CheckLogf if e != nil enriches the error message with the action text and additional context from v... and returns the extended error, otherwise nil
func CheckLogf(e error, action string, v ...interface{}) error {
	if e != nil {
		err := errors.Wrapf(e, action, v...)
		logging.Error(err)
		return err
	}
	return e
}

func CheckLogFatal(e error, action string) {
	if e != nil {
		logging.Fatal(action)
	}
}

func CheckLogFatalf(e error, action string, v ...interface{}) {
	if e != nil {
		newE := errors.Wrapf(e, action, v...)
		logging.Fatal(newE)
	}
}

func LogAndWrapAsErrorWarning(action string, v ...interface{}) error {
	errorMessage := fmt.Sprintf(action, v...)
	logging.Warning(errorMessage)
	return errors.New(errorMessage)
}

func LogAndWrapAsError(action string, v ...interface{}) error {
	errorMessage := fmt.Sprintf(action, v...)
	logging.Error(errorMessage)
	return errors.New(errorMessage)
}
