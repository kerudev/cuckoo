package main

import (
	"github.com/kerudev/cuckoo/cuckoo"
)

func main() {
	sample := []string{
		"0 1,16,20 * * *",
		"0 1,16,20 * * *",
		"25 1,16,20 * * *",
		"25 1,16,20 * * *",
		"25 1,16,20 * * *",
		"24 6 * * *",
		"24 7 * * *",
		"24 7 * * *",
		"44 1,12 * * *",
		"46 1,12 * * *",
	}

	cuckoo.DrawLoop(sample)
}
