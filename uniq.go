package main

import (
	"bufio"
	"fmt"
	"os"

	"uniq/parameters"
	"uniq/uniq"
)

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

	outputLines, err := uniq.UtilityUniq(inputLines, flags)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	printLines(outputFile, outputLines)
}
