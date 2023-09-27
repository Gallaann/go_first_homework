package uniq

import (
	"fmt"
	"strconv"
	"strings"

	"uniq/parameters"
)

func SkipNFields(line string, fieldsNum int) string {
	fields := strings.Fields(line)
	if len(fields) > fieldsNum && fieldsNum >= 0 {
		return strings.Join(fields[fieldsNum:], " ")
	}

	return ""
}

func SkipNSymbols(line string, symbolsNum int) string {
	if len(line) > symbolsNum {
		return line[symbolsNum:]
	}

	return ""
}

func AreLinesEqual(lineOne, lineTwo string, flags parameters.CmdFlags) bool {
	if flags.IgnoreCase {
		return strings.EqualFold(SkipNSymbols(SkipNFields(lineOne, flags.Fields), flags.Symbols), SkipNSymbols(SkipNFields(lineTwo, flags.Fields), flags.Symbols))
	}

	return SkipNSymbols(SkipNFields(lineOne, flags.Fields), flags.Symbols) == SkipNSymbols(SkipNFields(lineTwo, flags.Fields), flags.Symbols)
}

func ProcessLine(line string, count int, flags parameters.CmdFlags) string {
	if flags.Count {
		return fmt.Sprintf("%s %s", strconv.Itoa(count), line)
	}
	if flags.Duplicates && count > 1 {
		return line
	}
	if flags.Unique && count == 1 {
		return line
	}
	if !flags.Unique && !flags.Duplicates && !flags.Count {
		return line
	}

	return ""
}

func UtilityUniq(lines []string, flags parameters.CmdFlags) (output []string) {
	parameters.CheckFlags(flags)

	if len(lines) == 0 {
		return output
	}

	prevLine := lines[0]
	count := 1

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if AreLinesEqual(line, prevLine, flags) {
			count++
		} else {
			outputLine := ProcessLine(prevLine, count, flags)
			if outputLine != "" {
				output = append(output, outputLine)
			}
			prevLine = line
			count = 1
		}
	}

	outputLine := ProcessLine(prevLine, count, flags)
	if outputLine != "" {
		output = append(output, outputLine)
	}

	return output
}
