package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kerudev/cuckoo/cuckoo"
)

func main() {
	path := ""
	flag.StringVar(&path, "path", "", "Path where the data is")

	flag.Parse()

	if path == "" {
		fmt.Println("Please provide a path")
		os.Exit(1)
	}

	sample := map[string]string{}
	cuckoo.ReadPath(path, &sample)

	cuckoo.DrawLoop(sample)
}
