package app

import "flag"

func Run() {
	var path string
	flag.StringVar(&path, "path", "", "Path where the data is")

	flag.Parse()

	DrawLoop(path)
}
