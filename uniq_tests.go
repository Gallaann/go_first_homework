package main

import (
	"testing"

	"uniq/parameters"
	"uniq/uniq"
)

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
		{"Compare similar case characters with ignore case flag", args{"abc", "abc", parameters.CmdFlags{IgnoreCase: true}}, true},
		{"Compare similar case characters without ignore case flag", args{"abc", "abc", parameters.CmdFlags{IgnoreCase: false}}, true},
		{"Compare different case characters with ignore case flag", args{"abc", "ABC", parameters.CmdFlags{IgnoreCase: true}}, true},
		{"Compare different case characters without ignore case flag", args{"abc", "ABC", parameters.CmdFlags{IgnoreCase: false}}, false},
		{"Compare with skip field", args{"one two", "uno two", parameters.CmdFlags{Fields: 1}}, true},
		{"Compare without skip field", args{"one two", "uno two", parameters.CmdFlags{Fields: 0}}, false},
		{"Compare with skip symbols", args{"one two", "ONe two", parameters.CmdFlags{Symbols: 2}}, true},
		{"Compare without skip symbols", args{"one two", "ONe two", parameters.CmdFlags{Symbols: 0}}, false},
		{"Compare with skip fields and symbols", args{"skip one two", "not ONe two", parameters.CmdFlags{Fields: 1, Symbols: 2}}, true},
		{"Compare without skip fields and symbols", args{"skip one two", "not ONe two", parameters.CmdFlags{Fields: 0, Symbols: 0}}, false},
		{"Compare with all compare flags", args{"skip one tWo", "not ONe tWo", parameters.CmdFlags{IgnoreCase: true, Fields: 1, Symbols: 2}}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := uniq.AreLinesEqual(tt.args.lineOne, tt.args.lineTwo, tt.args.flags); got != tt.want {
				t.Errorf("AreLinesEqual() = %v, wantLine %v", got, tt.want)
			}
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
	tests := []struct {
		name     string
		args     args
		wantLine string
	}{
		{name: "TestCountFlag", args: args{line: "apple", count: 3, flags: parameters.CmdFlags{Count: true}}, wantLine: "3 apple"},
		{name: "TestDuplicatesFlag", args: args{line: "banana", count: 2, flags: parameters.CmdFlags{Duplicates: true}}, wantLine: "banana"},
		{name: "TestUniqueFlag", args: args{line: "cherry", count: 1, flags: parameters.CmdFlags{Unique: true}}, wantLine: "cherry"},
		{name: "TestNoFlag", args: args{line: "grape", count: 2, flags: parameters.CmdFlags{}}, wantLine: "grape"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := uniq.ProcessLine(tt.args.line, tt.args.count, tt.args.flags); got != tt.wantLine {
				t.Errorf("ProcessLine() = %v, wantLine %v", got, tt.wantLine)
			}
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
		{"Skip less then number of fields in line", args{"one two three four", 2}, "three four"},
		{"Skip more then number of fields in line", args{"one two three four", 5}, ""},
		{"Skip number is equal to number of fields in line", args{"one two three four", 4}, ""},
		{"Skip 0 fields", args{"abc def ghi", 0}, "abc def ghi"},
		{"Empty line", args{"", 2}, ""},
		{"Skip negative number of fields", args{"one two three", -1}, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := uniq.SkipNFields(tt.args.line, tt.args.fieldsNum); got != tt.want {
				t.Errorf("SkipNFields() = %v, wantLine %v", got, tt.want)
			}
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
		{"Skip less then number of symbols", args{"abcdef", 2}, "cdef"},
		{"Skip more then number of symbols in line", args{"two", 5}, ""},
		{"Skip number is equal to number of symbols in line", args{"four", 4}, ""},
		{"Skip 0 symbols", args{"12345", 0}, "12345"},
		{"Empty line", args{"", 2}, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := uniq.SkipNSymbols(tt.args.line, tt.args.symbolsNum); got != tt.want {
				t.Errorf("SkipNSymbols() = %v, wantLine %v", got, tt.want)
			}
		})
	}
}
