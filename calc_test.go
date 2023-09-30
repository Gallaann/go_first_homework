package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_calc(t *testing.T) {
	t.Parallel()
	type args struct {
		expression string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "Valid Expression Test",
			args:    args{expression: "2 + 3 * 4"},
			want:    14.0,
			wantErr: false,
		},
		{
			name:    "Invalid Expression Test",
			args:    args{expression: "2 + (3 * 4"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Division By Zero Test",
			args:    args{expression: "5 / 0"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid Token Test",
			args:    args{expression: "2 + x"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Big expression Test",
			args:    args{expression: "(2 + 3 * 4) - (5 / 2) + (7 * 8) - (9 + 10 / 2) + (11 * 12) - (13 / 3) + (14 * 15) - (16 + 17 / 4)"},
			want:    370.91666666666663,
			wantErr: false,
		},
		{
			name:    "Whitespaces Test",
			args:    args{expression: "2 + \n4 - \t3 *      \n\r 10\v - 0\f"},
			want:    -24.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := calc(tt.args.expression)
			if tt.wantErr {
				require.Error(t, err, "Expected an error, but got nil")
			} else {
				require.NoError(t, err, "Expected no error, but got an error")
				require.Equal(t, tt.want, got, "UtilityUniq() did not return the expected result")
			}
		})
	}
}

func Test_convertInfixToPostfix(t *testing.T) {
	t.Parallel()
	type args struct {
		tokens []string
	}
	tests := []struct {
		name                  string
		args                  args
		wantPostfixExpression []string
	}{
		{
			name:                  "Simple Infix to Postfix Test",
			args:                  args{tokens: []string{"2", "+", "3"}},
			wantPostfixExpression: []string{"2", "3", "+"},
		},
		{
			name:                  "Complex Infix to Postfix Test",
			args:                  args{tokens: []string{"2", "+", "3", "*", "4"}},
			wantPostfixExpression: []string{"2", "3", "4", "*", "+"},
		},
		{
			name:                  "Complex Infix to Postfix expression with scopes Test",
			args:                  args{tokens: []string{"(", "2", "+", "3", "*", "4", ")", "-", "(", "5", "/", "2", ")", "+", "(", "7", "*", "8", ")", "-", "(", "9", "+", "10", "/", "2", ")", "+", "(", "11", "*", "12", ")", "-", "(", "13", "/", "3", ")", "+", "(", "14", "*", "15", ")", "-", "(", "16", "+", "17", "/", "4", ")"}},
			wantPostfixExpression: []string{"2", "3", "4", "*", "+", "5", "2", "/", "-", "7", "8", "*", "+", "9", "10", "2", "/", "+", "-", "11", "12", "*", "+", "13", "3", "/", "-", "14", "15", "*", "+", "16", "17", "4", "/", "+", "-"},
		},
		{
			name:                  "Infix to Postfix with same precedence Test",
			args:                  args{tokens: []string{"2", "*", "3", "/", "4"}},
			wantPostfixExpression: []string{"2", "3", "*", "4", "/"},
		},
		{
			name:                  "Invalid Infix to Postfix Test",
			args:                  args{tokens: []string{"(", "2", "+", "3", "*", "4"}},
			wantPostfixExpression: nil,
		},
		{
			name:                  "Invalid Infix to Postfix Test",
			args:                  args{tokens: []string{"2", "+", "3", ")", "*", "4"}},
			wantPostfixExpression: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotPostfixExpression := convertInfixToPostfix(tt.args.tokens)
			require.Equal(t, tt.wantPostfixExpression, gotPostfixExpression, "convertInfixToPostfix() did not return the expected result")
		})
	}
}

func Test_evaluatePostfixExpression(t *testing.T) {
	t.Parallel()
	type args struct {
		postfixExpression []string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "Simple Postfix Evaluation Test",
			args:    args{postfixExpression: []string{"2", "3", "+"}},
			want:    5.0,
			wantErr: false,
		},
		{
			name:    "Complex Postfix Evaluation Test",
			args:    args{postfixExpression: []string{"2", "3", "4", "*", "+"}},
			want:    14.0,
			wantErr: false,
		},
		{
			name:    "Complex Expression Test",
			args:    args{postfixExpression: []string{"2", "3", "*", "4", "5", "+", "-"}},
			want:    -3,
			wantErr: false,
		},
		{
			name:    "Complex Expression Test",
			args:    args{postfixExpression: []string{"2", "3", "4", "*", "+", "5", "2", "/", "-", "7", "8", "*", "+", "9", "10", "2", "/", "+", "-", "11", "12", "*", "+", "13", "3", "/", "-", "14", "15", "*", "+", "16", "17", "4", "/", "+", "-"}},
			want:    370.91666666666663,
			wantErr: false,
		},
		{
			name:    "Division By Zero Test",
			args:    args{postfixExpression: []string{"2", "0", "/"}},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Addition and Subtraction Test",
			args:    args{postfixExpression: []string{"5", "2", "+", "3", "-"}},
			want:    4.0,
			wantErr: false,
		},
		{
			name:    "Multiplication and Division Test",
			args:    args{postfixExpression: []string{"6", "3", "*", "2", "/"}},
			want:    9.0,
			wantErr: false,
		},
		{
			name:    "Empty Expression Test",
			args:    args{postfixExpression: []string{}},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid Token Test",
			args:    args{postfixExpression: []string{"2", "x", "+"}},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid Token Test",
			args:    args{postfixExpression: []string{"/", "*", "+", "+"}},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Number only Test",
			args:    args{postfixExpression: []string{"397", "212", "12", "1"}},
			want:    0,
			wantErr: true,
		},
		{
			name:    "One number Test",
			args:    args{postfixExpression: []string{"1937"}},
			want:    1937.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := evaluatePostfixExpression(tt.args.postfixExpression)
			if tt.wantErr {
				require.Error(t, err, "Expected an error, but got nil")
			} else {
				require.NoError(t, err, "Expected no error, but got an error")
				require.Equal(t, tt.want, got, "evaluatePostfixExpression() did not return the expected result")
			}
		})
	}
}

func Test_extractExpressionTokens(t *testing.T) {
	t.Parallel()
	type args struct {
		expression string
	}
	tests := []struct {
		name       string
		args       args
		wantTokens []string
	}{
		{
			name:       "Simple Token Extraction Test",
			args:       args{expression: "2 + 3"},
			wantTokens: []string{"2", "+", "3"},
		},
		{
			name:       "Complex Token Extraction Test",
			args:       args{expression: "2 + (3 * 4)"},
			wantTokens: []string{"2", "+", "(", "3", "*", "4", ")"},
		},
		{
			name:       "Expression with Unary Minus Test",
			args:       args{expression: "2 + (-3 * 4)"},
			wantTokens: []string{"2", "+", "(", "~", "3", "*", "4", ")"},
		},
		{
			name:       "Expression with larger number of digits Test",
			args:       args{expression: "1000 + (-7 * 403)"},
			wantTokens: []string{"1000", "+", "(", "~", "7", "*", "403", ")"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotTokens := extractExpressionTokens(tt.args.expression)
			require.Equal(t, tt.wantTokens, gotTokens, "extractExpressionTokens() did not return the expected result")
		})
	}
}
