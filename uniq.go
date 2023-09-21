package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func getLines(inputParameter string) (lines []string) {
	input := os.Stdin

	if inputParameter != "" {
		input, err := os.Open(inputParameter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error openning input file: %v\n", err)
			os.Exit(1)
		}
		defer input.Close()
	}

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func processLine(line string, count int, flags parameters.Parameters) string {
	if flags.Count {
		return strconv.Itoa(count) + " " + line
	}
	if flags.Duplicates && count > 1 {
		return line
	}
	if flags.Unique && count == 1 {
		return line
	}
	if !flags.Unique && !flags.Duplicates {
		return line
	}

	return ""
}

func uniq(lines []string, flags parameters.Parameters) (output []string) {
	prevLine := lines[0]
	count := 1

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if equal(line, prevLine, flags) {
			count++
		} else {
			output = append(output, processLine(prevLine, count, flags))
			prevLine = line
			count = 1
		}
	}

	output = append(output, processLine(prevLine, count, flags))
	return output
}

func printLines(outputParameter string, lines []string) {
	output := os.Stdout

	if outputParameter != "" {
		output, err := os.Create(outputParameter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating/writing to output file: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
	}

	for _, line := range lines {
		fmt.Fprintln(output, line)
	}
}

func main() {
	flags := parameters.ParseFlags()

	parameters.CheckFlags(flags)

	inputFile, outputFile := parameters.ParseArguments()

	inputLines := getLines(inputFile)

	outputLines := uniq(inputLines, flags)

	printLines(outputFile, outputLines)
}
