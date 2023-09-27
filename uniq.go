package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"uniq/parameters"
)

func skipNFields(line string, fieldsNum int) string {
	fields := strings.Fields(line)
	if len(fields) > fieldsNum {
		return strings.Join(fields[fieldsNum:], " ")
	}

	return ""
}

func skipNSymbols(line string, symbolsNum int) string {
	if len(line) > symbolsNum {
		return line[symbolsNum:]
	}

	return ""
}

func areLinesEqual(lineOne, lineTwo string, flags parameters.CmdFlags) bool {
	if flags.IgnoreCase {
		return strings.EqualFold(skipNSymbols(skipNFields(lineOne, flags.Fields), flags.Symbols), skipNSymbols(skipNFields(lineTwo, flags.Fields), flags.Symbols))
	}

	return skipNSymbols(skipNFields(lineOne, flags.Fields), flags.Symbols) == skipNSymbols(skipNFields(lineTwo, flags.Fields), flags.Symbols)
}

func getLines(inputParameter string) (lines []string) {
	input := os.Stdin
	var err error

	if inputParameter != "" {
		input, err = os.Open(inputParameter)
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

func processLine(line string, count int, flags parameters.CmdFlags) string {
	if flags.Count {
		return fmt.Sprintf("%s %s", strconv.Itoa(count), line)
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

func uniq(lines []string, flags parameters.CmdFlags) (output []string) {
	if len(lines) == 0 {
		return
	}

	prevLine := lines[0]
	count := 1

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if areLinesEqual(line, prevLine, flags) {
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
	var err error

	if outputParameter != "" {
		output, err = os.Create(outputParameter)
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
