package main

import (
	"flag"
	"fmt"
)

type flags struct {
	count      bool
	duplicates bool
	unique     bool
	ignoreCase bool
	fields     int
	chars      int
}

func parseFlags() flags {
	var flags flags
	flag.BoolVar(&flags.count, "c", false, "count repeated lines")
	flag.BoolVar(&flags.duplicates, "d", false, "output only duplicate lines")
	flag.BoolVar(&flags.unique, "u", false, "output only unique lines")
	flag.BoolVar(&flags.ignoreCase, "i", false, "ignore case when comparing lines")
	flag.IntVar(&flags.fields, "f", 0, "skip first num fields in each line")
	flag.IntVar(&flags.chars, "s", 0, "skip first num chars in each line")
	flag.Parse()
	return flags
}
func checkFlags(flags flags) bool {
	sum := 0
	if flags.count {
		sum++
	}
	if flags.duplicates {
		sum++
	}
	if flags.unique {
		sum++
	}
	if sum != 1 {
		return false
	}

	return true
}

func main() {
	flags := parseFlags()

	fmt.Println(checkFlags(flags))
}
