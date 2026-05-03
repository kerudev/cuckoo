package main

import (
	"flag"
	"fmt"
	"os"

	cuckoo "github.com/kerudev/cuckoo/internal"
)

func main() {
	path := ""
	flag.StringVar(&path, "path", "", "Path where the data is")

	flag.Parse()

	if path == "" {
		fmt.Println("-path: provide a path to loads crons")
		os.Exit(1)
	}

	sample := map[string]string{}
	cuckoo.ReadPath(path, &sample)

	cuckoo.DrawLoop(sample)
}
