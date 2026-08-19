package utilities

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveJSON(dir, id string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, id+".json")
	return os.WriteFile(path, data, 0o644)
}
