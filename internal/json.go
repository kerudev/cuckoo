package cuckoo

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadPath(path string, payload any) (error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error while reading %s", path)
	}

	err = json.Unmarshal(content, payload)
	if err != nil {
		return fmt.Errorf("Error while unmarshaling data from %s", path)
	}

	return nil
}
