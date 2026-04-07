package cuckoo

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadPath(path string, payload any) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error while reading %s\n", path)
		os.Exit(1)
	}

	err = json.Unmarshal(content, payload)
	if err != nil {
		fmt.Printf("Error while unmarshaling data from %s\n", path)
		os.Exit(1)
	}
}
