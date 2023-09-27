package parameters

import (
	"flag"
	"fmt"
	"os"
)

type CmdFlags struct {
	Count      bool
	Duplicates bool
	Unique     bool
	IgnoreCase bool
	Fields     int
	Symbols    int
}

func ParseFlags() CmdFlags {
	var flags CmdFlags
	flag.BoolVar(&flags.Count, "c", false, "count repeated lines")
	flag.BoolVar(&flags.Duplicates, "d", false, "output only duplicate lines")
	flag.BoolVar(&flags.Unique, "u", false, "output only unique lines")
	flag.BoolVar(&flags.IgnoreCase, "i", false, "ignore case when comparing lines")
	flag.IntVar(&flags.Fields, "f", 0, "skip first num fields in each line")
	flag.IntVar(&flags.Symbols, "s", 0, "skip first num symbols in each line")
	flag.Parse()
	return flags
}

func ParseArguments() (string, string) {
	args := flag.Args()

	switch len(args) {
	case 0:
		return "", ""
	case 1:
		return args[0], ""
	default:
		return args[0], args[1]
	}
}

func CheckFlags(flags CmdFlags) {
	sum := 0
	if flags.Count {
		sum++
	}
	if flags.Duplicates {
		sum++
	}
	if flags.Unique {
		sum++
	}

	if sum > 1 {
		fmt.Fprintln(os.Stderr, "error: parameters [-c | -d | -u] are interchangeable, in parallel these parameters are meaningless")
		os.Exit(1)
	}
}
