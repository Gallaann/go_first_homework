package parameters

import (
	"errors"
	"flag"
	"io"
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

func ParseArguments() (io.Reader, io.Writer, error) {
	args := flag.Args()

	switch len(args) {
	case 0:
		return os.Stdin, os.Stdout, nil
	case 1:
		{
			return nil, os.Stdout, nil
		}
	default:
		{
			input, err := os.Open(args[0])
			if err != nil {
				return nil, nil, err
			}
			defer input.Close()

			output, err := os.Create(args[1])
			if err != nil {
				return nil, nil, err
			}
			defer output.Close()

			return input, output, nil
		}
	}
}

func CheckFlags(flags CmdFlags) error {
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
		return errors.New("error: parameters [-c | -d | -u] are interchangeable, in parallel these parameters are meaningless")
	}

	return nil
}
