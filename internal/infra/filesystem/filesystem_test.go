package filesystem

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadJson(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"name":"Pikachu"}`), 0644))

	var got struct {
		Name string `json:"name"`
	}

	err := ReadJson(path, &got)

	require.NoError(t, err)
	assert.Equal(t, "Pikachu", got.Name)
}

func TestReadJson_ReturnsReadError(t *testing.T) {
	var got map[string]string

	err := ReadJson(filepath.Join(t.TempDir(), "missing.json"), &got)

	assert.Error(t, err)
}

func TestReadJson_ReturnsUnmarshalError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.json")
	require.NoError(t, os.WriteFile(path, []byte(`not-json`), 0644))

	var got map[string]string

	err := ReadJson(path, &got)

	assert.Error(t, err)
}

func TestWriteCSV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "data.csv")

	err := WriteCSV(path, [][]string{{"Number", "Name"}, {"0025", "Pikachu"}})

	require.NoError(t, err)
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"Number", "Name"}, {"0025", "Pikachu"}}, records)
}

func TestWriteCSV_ReturnsCreateError(t *testing.T) {
	err := WriteCSV(t.TempDir(), [][]string{{"Number", "Name"}})

	assert.Error(t, err)
}
