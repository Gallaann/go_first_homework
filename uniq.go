package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
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

func parseArguments() (string, string) {
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

func checkFlags(flags flags) (bool, error) {
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
	if sum > 1 {
		return false, errors.New("parameters [-c | -d | -u] are interchangeable, in parallel these parameters are meaningless")
	}

	return true, nil
}

func skipFields(line string, numFields int) string {
	fields := strings.Fields(line)
	if len(fields) > numFields {
		return strings.Join(fields[numFields:], " ")
	}
	return ""
}

func skipChars(line string, numChars int) string {
	if len(line) > numChars {
		return line[numChars:]
	}
	return ""
}

func main() {
	flags := parseFlags()

	var err error

	if result, err := checkFlags(flags); !result {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var input io.Reader = os.Stdin
	var output io.Writer = os.Stdout

	inputFile, outputFile := parseArguments()
	if inputFile != "" {
		input, err = os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error openning input file: %v\n", err)
			os.Exit(1)
		}
		defer input.(*os.File).Close()
	}

	if outputFile != "" {
		output, err = os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating/writing to output file: %v\n", err)
			os.Exit(1)
		}
		defer output.(*os.File).Close()
	}

	scanner := bufio.NewScanner(input)

	// TODO realisation of uniq

	// TODO output data

	for scanner.Scan() {
		fmt.Fprintln(output, scanner.Text())
	}

	fmt.Println(flags, inputFile, outputFile)
}
