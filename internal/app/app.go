package app

import (
	"flag"
	"fmt"
	"os"
)

func Run() {
	path := ""
	flag.StringVar(&path, "path", "", "Path where the data is")

	flag.Parse()

	if path == "" {
		fmt.Println("-path: provide a path to loads crons")
		os.Exit(1)
	}

	DrawLoop(path)
}
