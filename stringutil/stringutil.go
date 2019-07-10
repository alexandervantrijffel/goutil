package stringutil

import (
	"encoding/json"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

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
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

// Returns empty string in case the argument is a nil pointer
func SafeDereferenceString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandomAlphanumericString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
func StructToTruncatedString(s interface{}, maxChars int) string {
	marshalled, err := json.Marshal(s)
	if err != nil {
		logging.Errorf("Failed to marshal string %s", err.Error())
		return "StructToTruncatedString: marshal failed"
	}
	return Truncate(string(marshalled), maxChars)
}
func FindStringSubmatchMap(r *regexp.Regexp, s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}
	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}
		captures[name] = match[i]
	}
	return captures
}
