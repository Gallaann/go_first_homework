package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"uniq/parameters"
)

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

func equal(lineOne, lineTwo string, flags parameters.Parameters) bool {
	if flags.IgnoreCase {
		return strings.EqualFold(skipChars(skipFields(lineOne, flags.Fields), flags.Chars), skipChars(skipFields(lineTwo, flags.Fields), flags.Chars))
	}

	return skipChars(skipFields(lineOne, flags.Fields), flags.Chars) == skipChars(skipFields(lineTwo, flags.Fields), flags.Chars)
}

func print(line string, flags parameters.Parameters, count int, output io.Writer) {
	if flags.Count {
		fmt.Fprintln(output, count, line)
	} else if flags.Duplicates && count > 1 {
		fmt.Fprintln(output, line)
	} else if flags.Unique && count == 1 {
		fmt.Fprintln(output, line)
	} else if !flags.Unique && !flags.Duplicates {
		fmt.Fprintln(output, line)
	}
}

func uniq(input io.Reader, output io.Writer, flags parameters.Parameters) {
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
		} else {
			print(prevLine, flags, count, output)
			prevLine = line
			count = 1
		}
	}

	print(prevLine, flags, count, output)
}

func main() {
	flags := parameters.ParseFlags()

	var err error

	if !parameters.CheckFlags(flags) {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	input := os.Stdin
	output := os.Stdout

	inputFile, outputFile := parameters.ParseArguments()
	if inputFile != "" {
		input, err = os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error openning input file: %v\n", err)
			os.Exit(1)
		}
		defer input.Close()
	}

	if outputFile != "" {
		output, err = os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating/writing to output file: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
	}

	uniq(input, output, flags)
}
