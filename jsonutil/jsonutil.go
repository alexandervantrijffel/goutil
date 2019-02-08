package jsonutil

import (
	"encoding/json"

	"github.com/alexandervantrijffel/goutil/errorcheck"
)

func UnmarshalWithLogging(destination interface{}, source []byte) error {
	err := json.Unmarshal(source, destination)
	_ = errorcheck.CheckLogf(err, "Failed to unmarshal data to type %T. (start of) data: %s", destination, firstChars(string(source), 512))
	return err
}

func firstChars(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
