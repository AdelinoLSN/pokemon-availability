package filesystem

import (
	"encoding/json"
	"os"
)

func ReadJson(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}
