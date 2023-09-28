package main

import (
	"fmt"
	"os"

	"uniq/parameters"
	"uniq/uniq"
)

func main() {
	flags := parameters.ParseFlags()

	err := parameters.CheckFlags(flags)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	inputStream, outputStream, err := parameters.ParseArguments()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	inputLines := uniq.GetLines(inputStream)

	outputLines, err := uniq.UtilityUniq(inputLines, flags)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = uniq.PrintLines(outputStream, outputLines)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
