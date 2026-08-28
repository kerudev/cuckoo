package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
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

func GetUnix(path string) int64 {
	stat, _ := os.Stat(path)
	return stat.ModTime().UnixNano()
}

func ReadPath(path string, payload *map[string]string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return errors.New("Error while reading file")
	}

	switch filepath.Ext(path) {
	case Csv:
		b, _ := os.ReadFile(path)

		lines, _ := csv.NewReader(bytes.NewReader(b)).ReadAll()
		if len(lines) == 0 {
			return errors.New("File doesn't contain headers nor data")
		}

		lines = lines[1:]
		if len(lines) == 0 {
			return errors.New("File contains just the headers")
		}

		// Two fields are expected on each line: name,cron
		for _, line := range lines {
			(*payload)[line[0]] = line[1]
		}
	case Json:
		err = json.Unmarshal(content, payload)
	}

	if err != nil {
		return errors.New("Error while unmarshaling data (file is empty or JSON has errors)")
	}

	return nil
}
