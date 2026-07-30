package app

import "flag"

func Run() {
	path := ""
	flag.StringVar(&path, "path", "", "Path where the data is")

	flag.Parse()

	DrawLoop(path)
}
