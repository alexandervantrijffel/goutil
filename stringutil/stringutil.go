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
func Truncate(s string, maxChars int) string {
	if len(s) > maxChars {
		return s[:maxChars-1] + " *SNIP*"
	}
	return s
}
func EqualsCaseInsensitive(a, b string) (equal bool) {
	return strings.ToLower(strings.TrimSpace(a)) == strings.ToLower(strings.TrimSpace(b))
}
