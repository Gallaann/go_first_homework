package main

import (
	"fmt"
	"uniq/parameters"
	"uniq/uniq"
)

func main() {
	flags := parameters.ParseFlags()

	err := parameters.CheckFlags(flags)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	inputStream, outputStream, err := parameters.ParseArguments()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	inputLines := uniq.GetLines(inputStream)

	outputLines, err := uniq.UtilityUniq(inputLines, flags)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = uniq.PrintLines(outputStream, outputLines)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}
