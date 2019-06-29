package iputil

import (
	"testing"

	"gotest.tools/assert"
)

func TestParseIP(t *testing.T) {
	ipFromStackpath := "87.208.232.54, 94.46.155.210"
	result := removeReverseProxyIP(ipFromStackpath)
	assert.Equal(t, "87.208.232.54", result)
}
