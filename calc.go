package main

import (
	"fmt"
	"os"
	"strconv"

	"calc/stack"
)

var kPrecedence = map[string]int{"+": 1, "-": 1, "*": 2, "/": 2, "~": 3}

func extractExpressionTokens(expression string) (tokens []string) {
	var currentToken string

	for i, char := range expression {
		switch char {
		case ' ', '\t', '\n', '\r', '\v', '\f': // whitespaces
			continue
		case '(', ')', '+', '-', '*', '/': // operations
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}

			if char == '-' && (i == 0 || expression[i-1] == '(') {
				tokens = append(tokens, "~")
				continue
			}

			tokens = append(tokens, string(char))
		default:
			currentToken += string(char)
		}
	}

	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}

	return tokens
}

func convertInfixToPostfix(tokens []string) (postfixExpression []string) {
	var operationsStack stack.Stack

	for _, token := range tokens {
		switch token {
		case "+", "-", "~":
			for !operationsStack.IsEmpty() && kPrecedence[operationsStack.Peek()] >= kPrecedence[token] {
				postfixExpression = append(postfixExpression, operationsStack.Pop())
			}

			operationsStack.Push(token)
		case "*", "/":
			for !operationsStack.IsEmpty() && kPrecedence[operationsStack.Peek()] >= kPrecedence[token] {
				postfixExpression = append(postfixExpression, operationsStack.Pop())
			}

			operationsStack.Push(token)
		case "(":
			operationsStack.Push(token)
		case ")":
			for !operationsStack.IsEmpty() && operationsStack.Peek() != "(" {
				postfixExpression = append(postfixExpression, operationsStack.Pop())
			}

			if operationsStack.IsEmpty() {
				return nil
			}

			operationsStack.Pop()
		default:
			postfixExpression = append(postfixExpression, token)
		}
	}

	for !operationsStack.IsEmpty() {
		if operationsStack.Peek() == "(" {
			return nil
		}

		postfixExpression = append(postfixExpression, operationsStack.Pop())
	}

	return postfixExpression
}

func evaluatePostfixExpression(postfixExpression []string) (float64, error) {
	var exitStack []float64

	for _, token := range postfixExpression {
		switch token {
		case "+", "-", "*", "/":
			if len(exitStack) < 2 {
				return 0, fmt.Errorf("invalid expression")
			}

			b, a := exitStack[len(exitStack)-1], exitStack[len(exitStack)-2]
			exitStack = exitStack[:len(exitStack)-2]

			switch token {
			case "+":
				exitStack = append(exitStack, a+b)
			case "-":
				exitStack = append(exitStack, a-b)
			case "*":
				exitStack = append(exitStack, a*b)
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}

				exitStack = append(exitStack, a/b)
			}
		default:
			num, err := strconv.ParseFloat(token, 64)

			if err != nil {
				return 0, fmt.Errorf("invalid token: %s", token)
			}

			exitStack = append(exitStack, num)
		}
	}

	if len(exitStack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	return exitStack[0], nil
}

func calc(expression string) (float64, error) {
	tokens := extractExpressionTokens(expression)
	postfixExpression := convertInfixToPostfix(tokens)
	result, err := evaluatePostfixExpression(postfixExpression)

	if err != nil {
		return 0, err
	}

	return result, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run calc.go <expression>")
		return
	}

	expression := os.Args[1]
	result, err := calc(expression)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(result)
}
