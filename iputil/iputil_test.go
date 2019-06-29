package iputil

import (
	"net/http"
	"testing"

	"github.com/alexandervantrijffel/goutil/logging"
	"gotest.tools/assert"
)

func TestGetIP(t *testing.T) {
	logging.InitWith("unit test", true)
	ipFromStackpath := "87.208.232.54, 94.46.155.210"
	result := GetIP(&http.Request{RemoteAddr: ipFromStackpath})
	assert.Equal(t, "87.208.232.54", result)
}
