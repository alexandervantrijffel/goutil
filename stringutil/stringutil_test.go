package stringutil

import (
	"testing"

	"gotest.tools/assert"
)

func TestNilString(t *testing.T) {
	result := SafeDereferenceString(nil)
	assert.Equal(t, "", result)
}
