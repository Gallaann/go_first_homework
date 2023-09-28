package uniq

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"testing"

	"uniq/parameters"
)

func TestUtilityUniq(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		input       []string
		flags       parameters.CmdFlags
		output      []string
		shouldError bool
	}{
		{
			name:   "Empty input",
			input:  []string{},
			flags:  parameters.CmdFlags{},
			output: []string{},
		},
		{
			name: "Without flags",
			input: []string{
				"I love music.",
				"I love music.",
				"I love music.",
				"",
				"I love music of Kartik.",
				"I love music of Kartik.",
			},
			flags: parameters.CmdFlags{},
			output: []string{
				"I love music.",
				"",
				"I love music of Kartik.",
			},
		},
		{
			name: "Test count lines, no duplicates",
			input: []string{
				"I love music.",
				"",
				"I love music.",
				"I love music of Kartik.",
				"Kartik.",
				"I love music of Kartik.",
			},
			flags: parameters.CmdFlags{
				Count: true,
			},
			output: []string{
				"1 I love music.",
				"1 ",
				"1 I love music.",
				"1 I love music of Kartik.",
				"1 Kartik.",
				"1 I love music of Kartik.",
			},
		},
		{
			name: "Test count lines, duplicates",
			input: []string{
				"I love music.",
				"I love music.",
				"I love music.",
				"I love music.",
				"I love music of Kartik.",
				"I love music of Kartik.",
				"I love music of Kartik.",
			},
			flags: parameters.CmdFlags{
				Count: true,
			},
			output: []string{
				"4 I love music.",
				"3 I love music of Kartik.",
			},
		},
		{
			name: "Test count lines, empty lines in the start",
			input: []string{
				"",
				"",
				"I love music.",
				"I love music.",
				"I love music of Kartik.",
				"Kartik.",
			},
			flags: parameters.CmdFlags{
				Count: true,
			},
			output: []string{
				"2 ",
				"2 I love music.",
				"1 I love music of Kartik.",
				"1 Kartik.",
			},
		},
		{
			name: "Test count lines, empty lines in the middle",
			input: []string{
				"I love music.",
				"I love music.",
				"",
				"",
				"I love music of Kartik.",
				"Kartik.",
			},
			flags: parameters.CmdFlags{
				Count: true,
			},
			output: []string{
				"2 I love music.",
				"2 ",
				"1 I love music of Kartik.",
				"1 Kartik.",
			},
		},
		{
			name: "Test count lines, empty lines in the end",
			input: []string{
				"I love music.",
				"I love music.",
				"I love music of Kartik.",
				"Kartik.",
				"",
				"",
			},
			flags: parameters.CmdFlags{
				Count: true,
			},
			output: []string{
				"2 I love music.",
				"1 I love music of Kartik.",
				"1 Kartik.",
				"2 ",
			},
		},
		{
			name: "Test duplicates flag, no duplicates",
			input: []string{
				"I love music.",
				"",
				"I love music of Kartik.",
				"Kartik.",
			},
			flags: parameters.CmdFlags{
				Duplicates: true,
			},
			output: []string{},
		},
		{
			name: "Test duplicates flag, duplicates",
			input: []string{
				"I love music.",
				"I love music.",
				"I love music.",
				"Kartik.",
				"Kartik.",
			},
			flags: parameters.CmdFlags{
				Duplicates: true,
			},
			output: []string{
				"I love music.",
				"Kartik.",
			},
		},
		{
			name: "Test duplicates flag, empty lines",
			input: []string{
				"",
				"",
				"",
			},
			flags: parameters.CmdFlags{
				Duplicates: true,
			},
			output: []string{
				"",
			},
		},
		{
			name: "Test duplicates flag, empty lines in the start",
			input: []string{
				"",
				"",
				"I love music.",
				"I love music.",
				"I love music of Kartik.",
				"Kartik.",
			},
			flags: parameters.CmdFlags{
				Duplicates: true,
			},
			output: []string{
				"",
				"I love music.",
			},
		},
		{
			name: "Test duplicates flag, empty lines in the middle",
			input: []string{
				"I love music.",
				"I love music.",
				"",
				"",
				"I love music of Kartik.",
				"Kartik.",
			},
			flags: parameters.CmdFlags{
				Duplicates: true,
			},
			output: []string{
				"I love music.",
				"",
			},
		},
		{
			name: "Test duplicates flag, empty lines in the end",
			input: []string{
				"I love music.",
				"I love music.",
				"I love music of Kartik.",
				"Kartik.",
				"",
				"",
			},
			flags: parameters.CmdFlags{
				Duplicates: true,
			},
			output: []string{
				"I love music.",
				"",
			},
		},
		{
			name: "Test unique flag, one empty line",
			input: []string{
				"",
			},
			flags: parameters.CmdFlags{
				Unique: true,
			},
			output: []string{
				"",
			},
		},
		{
			name: "Test unique flag, empty lines",
			input: []string{
				"",
				"",
				"",
			},
			flags: parameters.CmdFlags{
				Unique: true,
			},
			output: []string{},
		},
		{
			name: "Test unique flag, ordinary",
			input: []string{
				"I love music.",
				"I love music.",
				"I love music.",
				"",
				"I love music of Kartik.",
				"I love music of Kartik.",
			},
			flags: parameters.CmdFlags{
				Unique: true,
			},
			output: []string{
				"",
			},
		},
		{
			name: "Test ignore case flag, ordinary",
			input: []string{
				"I love music.",
				"I lOvE muSic.",
				"i love MUSIC.",
				"",
				"I love music of Kartik.",
				"I love music of Kartik.",
			},
			flags: parameters.CmdFlags{
				IgnoreCase: true,
			},
			output: []string{
				"I love music.",
				"",
				"I love music of Kartik.",
			},
		},
		{
			name: "Test skip fields flag, ordinary",
			input: []string{
				"I love music.",
				"He love music.",
				"Me love music.",
				"",
				"I love music of Kartik.",
				"We love music of Kartik.",
			},
			flags: parameters.CmdFlags{
				Fields: 1,
			},
			output: []string{
				"I love music.",
				"",
				"I love music of Kartik.",
			},
		},
		{
			name: "Test skip symbols flag, ordinary",
			input: []string{
				"I love music.",
				"A love music.",
				"B love music.",
				"",
				"I love music of Kartik.",
				"H love music of Kartik.",
			},
			flags: parameters.CmdFlags{
				Symbols: 1,
			},
			output: []string{
				"I love music.",
				"",
				"I love music of Kartik.",
			},
		},
		{
			name: "Test count && duplicates flags",
			input: []string{
				"testing",
			},
			flags: parameters.CmdFlags{
				Count:      true,
				Duplicates: true,
				Unique:     true,
				IgnoreCase: true,
				Fields:     1,
				Symbols:    1,
			},
			output:      nil,
			shouldError: true,
		},
		{
			name: "Test count && unique flags",
			input: []string{
				"testing",
			},
			flags: parameters.CmdFlags{
				Count:      true,
				Duplicates: true,
				Unique:     true,
				IgnoreCase: true,
				Fields:     1,
				Symbols:    1,
			},
			output:      nil,
			shouldError: true,
		},
		{
			name: "Test unique && duplicates flags",
			input: []string{
				"testing",
			},
			flags: parameters.CmdFlags{
				Count:      true,
				Duplicates: true,
				Unique:     true,
				IgnoreCase: true,
				Fields:     1,
				Symbols:    1,
			},
			output:      nil,
			shouldError: true,
		},
		{
			name: "Test count && duplicates && unique flags",
			input: []string{
				"testing",
			},
			flags: parameters.CmdFlags{
				Count:      true,
				Duplicates: true,
				Unique:     true,
				IgnoreCase: true,
				Fields:     1,
				Symbols:    1,
			},
			output:      nil,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := UtilityUniq(tt.input, tt.flags)
			if tt.shouldError {
				require.Error(t, err, "Expected an error, but got nil")
			} else {
				require.NoError(t, err, "Expected no error, but got an error")
				require.Equal(t, tt.output, got, "UtilityUniq() did not return the expected result")
			}
		})
	}
}

func Test_AreLinesEqual(t *testing.T) {
	t.Parallel()
	type args struct {
		lineOne string
		lineTwo string
		flags   parameters.CmdFlags
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Compare similar case characters with ignore case flag",
			args{"abc", "abc", parameters.CmdFlags{IgnoreCase: true}},
			true,
		},
		{
			"Compare similar case characters without ignore case flag",
			args{"abc", "abc", parameters.CmdFlags{}},
			true,
		},
		{
			"Compare different case characters with ignore case flag",
			args{"abc", "ABC", parameters.CmdFlags{IgnoreCase: true}},
			true,
		},
		{
			"Compare different case characters without ignore case flag",
			args{"abc", "ABC", parameters.CmdFlags{}},
			false,
		},
		{
			"Compare with skip field",
			args{"one two", "uno two", parameters.CmdFlags{Fields: 1}},
			true,
		},
		{
			"Compare without skip field",
			args{"one two", "uno two", parameters.CmdFlags{}},
			false,
		},
		{
			"Compare with skip symbols",
			args{"one two", "ONe two", parameters.CmdFlags{Symbols: 2}},
			true,
		},
		{
			"Compare without skip symbols",
			args{"one two", "ONe two", parameters.CmdFlags{}},
			false,
		},
		{
			"Compare with skip fields and symbols",
			args{"skip one two", "not ONe two", parameters.CmdFlags{Fields: 1, Symbols: 2}},
			true,
		},
		{
			"Compare without skip fields and symbols",
			args{"skip one two", "not ONe two", parameters.CmdFlags{}},
			false,
		},
		{
			"Compare with all compare flags",
			args{"skip one tWo", "not ONe tWo", parameters.CmdFlags{IgnoreCase: true, Fields: 1, Symbols: 2}},
			true,
		},
		{
			"Compare without all compare flags",
			args{"skip one tWo", "not ONe tWo", parameters.CmdFlags{}},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := AreLinesEqual(tt.args.lineOne, tt.args.lineTwo, tt.args.flags)
			require.Equal(t, tt.want, got, "AreLinesEqual() did not return the expected result")
		})
	}
}

func Test_ProcessLine(t *testing.T) {
	t.Parallel()
	type args struct {
		line  string
		count int
		flags parameters.CmdFlags
	}
	type LineResult struct {
		Line         string
		MatchesFlags bool
	}
	tests := []struct {
		name string
		args args
		want LineResult
	}{
		{
			name: "Test count flag",
			args: args{
				line:  "Music",
				count: 3,
				flags: parameters.CmdFlags{
					Count: true,
				},
			},
			want: LineResult{
				"3 Music",
				true,
			},
		},
		{
			name: "Test duplicates flag",
			args: args{
				line:  "Music",
				count: 2,
				flags: parameters.CmdFlags{
					Duplicates: true,
				},
			},
			want: LineResult{
				"Music",
				true,
			},
		},
		{
			name: "Test duplicates flag, 1 matching line",
			args: args{
				line:  "Music",
				count: 1,
				flags: parameters.CmdFlags{
					Duplicates: true,
				},
			},
			want: LineResult{
				"Music",
				false,
			},
		},
		{
			name: "Test unique flag, 1 matching line",
			args: args{
				line:  "Music",
				count: 1,
				flags: parameters.CmdFlags{
					Unique: true,
				},
			},
			want: LineResult{
				"Music",
				true,
			},
		},
		{
			name: "Test unique flag, more then 1 matching line",
			args: args{
				line:  "Music",
				count: 3,
				flags: parameters.CmdFlags{
					Unique: true,
				},
			},
			want: LineResult{
				"Music",
				false,
			},
		},
		{
			name: "Test no flag",
			args: args{
				line:  "Music",
				count: 2,
				flags: parameters.CmdFlags{},
			},
			want: LineResult{
				"Music",
				true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ProcessLine(tt.args.line, tt.args.count, tt.args.flags)
			require.Equal(t, tt.want.Line, got.Line, "ProcessLine() did not return the expected result")
		})
	}
}

func Test_SkipNFields(t *testing.T) {
	t.Parallel()
	type args struct {
		line      string
		fieldsNum int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Skip less then number of fields in line",
			args{"one two three four", 2},
			"three four",
		},
		{
			"Skip more then number of fields in line",
			args{"one two three four", 5},
			"",
		},
		{
			"Skip number is equal to number of fields in line",
			args{"one two three four", 4},
			"",
		},
		{
			"Skip 0 fields",
			args{"abc def ghi", 0},
			"abc def ghi",
		},
		{
			"Empty line, skip some fields",
			args{"", 2},
			"",
		},
		{
			"Empty line, skip zero fields",
			args{"", 0},
			"",
		},
		{
			"Empty line, skip negative number of fields",
			args{"", -2},
			"",
		},
		{
			"Skip negative number of fields",
			args{"one two three", -1},
			"",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := SkipNFields(tt.args.line, tt.args.fieldsNum)
			require.Equal(t, tt.want, got, "SkipNFields() did not return the expected result")
		})
	}
}

func Test_SkipNSymbols(t *testing.T) {
	t.Parallel()
	type args struct {
		line       string
		symbolsNum int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Skip less then number of symbols",
			args{"Precompiled", 3},
			"compiled",
		},
		{
			"Skip more then number of symbols in line",
			args{"two", 5},
			"",
		},
		{
			"Skip number is equal to number of symbols in line",
			args{"four", 4},
			"",
		},
		{
			"Skip 0 symbols",
			args{"12345", 0},
			"12345",
		},
		{
			"Not empty line, skip negative number of symbols",
			args{"word", -3},
			"",
		},
		{
			"Empty line on input, skip some symbols",
			args{"", 2},
			"",
		},
		{
			"Empty line on input, skip 0 symbols",
			args{"", 0},
			"",
		},
		{
			"Empty line on input, skip negative number of symbols",
			args{"", -3},
			"",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := SkipNSymbols(tt.args.line, tt.args.symbolsNum)
			require.Equal(t, tt.want, got, "SkipNSymbols() did not return the expected result")
		})
	}
}

func TestGetLines(t *testing.T) {
	t.Parallel()
	type args struct {
		input io.Reader
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty input",
			args: args{input: bytes.NewBufferString("")},
			want: []string{},
		},
		{
			name: "Empty lines input",
			args: args{input: bytes.NewBufferString("\n\n\n\n")},
			want: []string{
				"",
				"",
				"",
				"",
			},
		},
		{
			name: "Typical input",
			args: args{input: bytes.NewBufferString(
				"I love music.\n\nI love music of Kartik.",
			)},
			want: []string{
				"I love music.",
				"",
				"I love music of Kartik.",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := GetLines(tt.args.input)
			require.Equal(t, tt.want, got, "GetLines() did not return the expected result")
		})
	}
}

func TestPrintLines(t *testing.T) {
	t.Parallel()
	type args struct {
		lines []string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
		wantErr    bool
	}{
		{
			name: "Ordinary output",
			args: args{lines: []string{
				"I love music.",
				"",
				"I love music of Kartik.",
			}},
			wantOutput: "I love music.\n\nI love music of Kartik.\n",
		},
		{
			name:       "Empty lines",
			args:       args{lines: []string{}},
			wantOutput: "",
		},
		{
			name: "New lines",
			args: args{lines: []string{
				"",
				"",
				"",
			}},
			wantOutput: "\n\n\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := &bytes.Buffer{}
			err := PrintLines(output, tt.args.lines)

			if tt.wantErr {
				require.Error(t, err, "Expected an error, but got nil")
			} else {
				gotOutput := output.String()
				require.Equal(t, tt.wantOutput, gotOutput, "PrintLines() did not return the expected result")
			}
		})
	}
}
