package app

import (
	"flag"
	"fmt"
	"os"

	utils "github.com/kerudev/cuckoo/internal/utils"
)

func Run() {
	path := ""
	flag.StringVar(&path, "path", "", "Path where the data is")

	flag.Parse()

	if path == "" {
		fmt.Println("-path: provide a path to loads crons")
		os.Exit(1)
	}

	sample := map[string]string{}
	utils.ReadPath(path, &sample)

	DrawLoop(sample)
}
