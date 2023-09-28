package uniq

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"uniq/parameters"
)

type LineResult struct {
	Line         string
	MatchesFlags bool
}

func SkipNFields(line string, fieldsNum int) string {
	fields := strings.Fields(line)
	if len(fields) > fieldsNum && fieldsNum >= 0 {
		return strings.Join(fields[fieldsNum:], " ")
	}

	return ""
}

func SkipNSymbols(line string, symbolsNum int) string {
	if len(line) > symbolsNum && symbolsNum >= 0 {
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

func ProcessLine(line string, count int, flags parameters.CmdFlags) LineResult {
	result := LineResult{Line: line}

	if flags.Count {
		result.Line = fmt.Sprintf("%s %s", strconv.Itoa(count), line)
		result.MatchesFlags = true
	}
	if flags.Duplicates && count > 1 {
		result.MatchesFlags = true
	}
	if flags.Unique && count == 1 {
		result.MatchesFlags = true
	}
	if !flags.Unique && !flags.Duplicates && !flags.Count {
		result.MatchesFlags = true
	}

	return result
}

func GetLines(input io.Reader) []string {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func PrintLines(output io.Writer, lines []string) error {
	for _, line := range lines {
		_, err := fmt.Fprintln(output, line)

		if err != nil {
			return err
		}
	}

	return nil
}

func UtilityUniq(lines []string, flags parameters.CmdFlags) ([]string, error) {
	err := parameters.CheckFlags(flags)

	if err != nil {
		return nil, err
	}

	output := make([]string, 0)

	if len(lines) == 0 {
		return output, nil
	}

	prevLine := lines[0]
	count := 1

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if AreLinesEqual(line, prevLine, flags) {
			count++
			continue
		}

		result := ProcessLine(prevLine, count, flags)
		if result.MatchesFlags {
			output = append(output, result.Line)
		}

		prevLine = line
		count = 1
	}

	result := ProcessLine(prevLine, count, flags)
	if result.MatchesFlags {
		output = append(output, result.Line)
	}

	return output, err
}
