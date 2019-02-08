package jsonutil

import (
	"testing"

	"github.com/alexandervantrijffel/goutil/logging"
	"github.com/alexandervantrijffel/scrape/config"
	"github.com/stretchr/testify/assert"
)

func TestShouldReportError(t *testing.T) {
	logging.InitWith("goutil_test", config.DEBUG)
	var dummy interface{}
	err := UnmarshalWithLogging(&dummy, []byte(`{ "invalid":json}`))
	assert.NotNil(t, err)
	assert.Equal(t, "invalid character 'j' looking for beginning of value", err.Error())
}
