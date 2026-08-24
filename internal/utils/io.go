package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	Csv  = ".csv"
	Json = ".json"
)

var AllowedExt = []string{
	Csv,
	Json,
}

func IsDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}

func ReadPath(path string, payload *map[string]string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error while reading %s", path)
	}

	switch filepath.Ext(path) {
	case Csv:
		b, _ := os.ReadFile(path)

		lines, _ := csv.NewReader(bytes.NewReader(b)).ReadAll()
		if len(lines) == 0 {
			return fmt.Errorf("File doesn't contain headers nor data: %s", path)
		}

		lines = lines[1:]
		if len(lines) == 0 {
			return fmt.Errorf("File contains just the headers: %s", path)
		}

		// Two fields are expected on each line: name,cron
		for _, line := range lines {
			(*payload)[line[0]] = line[1]
		}
	case Json:
		err = json.Unmarshal(content, payload)
	}

	if err != nil {
		return fmt.Errorf("Error while unmarshaling data from %s", path)
	}

	return nil
}
