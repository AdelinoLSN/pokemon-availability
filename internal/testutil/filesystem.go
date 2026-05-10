package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateTempDir(t *testing.T, dirpath string, pattern string) string {
	t.Helper()

	dir, err := os.MkdirTemp(dirpath, pattern)
	assert.NoError(t, err)
	return dir
}

func CreateTempFile(t *testing.T, dirpath string, pattern string, content string) *os.File {
	t.Helper()

	file, err := os.CreateTemp(dirpath, pattern)
	assert.NoError(t, err)

	_, err = file.Write([]byte(content))
	assert.NoError(t, err)

	return file
}

func CreateFile(t *testing.T, dirpath string, filename string, content string) {
	t.Helper()

	exactPath := filepath.Join(dirpath, filename)
	err := os.WriteFile(exactPath, []byte(content), 0644)
	assert.NoError(t, err)
}
