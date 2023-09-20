package parameters

import "flag"

type Parameters struct {
	Count      bool
	Duplicates bool
	Unique     bool
	IgnoreCase bool
	Fields     int
	Chars      int
}

func ParseFlags() Parameters {
	var flags Parameters
	flag.BoolVar(&flags.Count, "c", false, "count repeated lines")
	flag.BoolVar(&flags.Duplicates, "d", false, "output only duplicate lines")
	flag.BoolVar(&flags.Unique, "u", false, "output only unique lines")
	flag.BoolVar(&flags.IgnoreCase, "i", false, "ignore case when comparing lines")
	flag.IntVar(&flags.Fields, "f", 0, "skip first num fields in each line")
	flag.IntVar(&flags.Chars, "s", 0, "skip first num chars in each line")
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

func CheckFlags(flags Parameters) bool {
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
		return false
	}

	return true
}
