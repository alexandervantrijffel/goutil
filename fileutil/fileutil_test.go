package fileutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomizeFilePath(t *testing.T) {
	path := "/test/a/folder/file1.jpg"
	result := RandomizeFileName(path)
	assert.Regexp(t, "/test/a/folder/file1-[a-zA-Z0-9]{8}.jpg", result)
}
func TestRandomizeFileName(t *testing.T) {
	path := "file1.jpg"
	result := RandomizeFileName(path)
	assert.Regexp(t, "file1-[a-zA-Z0-9]{8}.jpg", result)
}
func TestRandomizeFileNoExtension(t *testing.T) {
	path := "file1"
	result := RandomizeFileName(path)
	assert.Regexp(t, "file1-[a-zA-Z0-9]{8}", result)
}
