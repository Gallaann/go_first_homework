package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type parameters struct {
	count      bool
	duplicates bool
	unique     bool
	ignoreCase bool
	fields     int
	chars      int
}

func parseFlags() parameters {
	var flags parameters
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

func checkFlags(flags parameters) bool {
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
		return false
	}

	return true
}

func skipFields(line string, fields int) string {
	fielders := strings.Fields(line)
	if len(fielders) > fields {
		return strings.Join(fielders[fields:], " ")
	}
	return ""
}

func skipChars(line string, chars int) string {
	if len(line) > chars {
		return line[chars:]
	}
	return ""
}

func equal(lineOne, lineTwo string, flags parameters) bool {
	if flags.ignoreCase {
		return strings.EqualFold(skipChars(skipFields(lineOne, flags.fields), flags.chars), skipChars(skipFields(lineTwo, flags.fields), flags.chars))
	}

	return skipChars(skipFields(lineOne, flags.fields), flags.chars) == skipChars(skipFields(lineTwo, flags.fields), flags.chars)
}

func print(line string, flags parameters, count int, output io.Writer) {
	if flags.count {
		fmt.Fprintln(output, count, line)
	} else if flags.duplicates && count > 1 {
		fmt.Fprintln(output, line)
	} else if flags.unique && count == 1 {
		fmt.Fprintln(output, line)
	} else if !flags.unique && !flags.duplicates {
		fmt.Fprintln(output, line)
	}
}

func realisation(input io.Reader, output io.Writer, flags parameters) {
	scanner := bufio.NewScanner(input)

	var prevLine string
	if scanner.Scan() {
		prevLine = scanner.Text()
	}

	count := 1
	for scanner.Scan() {
		line := scanner.Text()

		if equal(line, prevLine, flags) {
			count++
			continue
		} else {
			print(prevLine, flags, count, output)
			prevLine = line
			count = 1
		}
	}

	print(prevLine, flags, count, output)
}

func main() {
	flags := parseFlags()

	var err error

	if !checkFlags(flags) {
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

	realisation(input, output, flags)
}
