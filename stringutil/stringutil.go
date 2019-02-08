package stringutil

import (
	"strconv"

	"github.com/alexandervantrijffel/goutil/logging"
)

func AtoiWithLogging(input string) int {
	var iDefault int
	if len(input) == 0 {
		logging.Warning("Converting empty value to int. Setting it to default value 0.")
		return iDefault
	}
	if i, err := strconv.Atoi(input); err == nil {
		return i
	} else {
		logging.Errorf("Failed to convert value '%s' to int", input)
		return iDefault
	}
}
