package main

import (
	"fmt"
	"os"
	"strconv"

	"calc/stack"
)

func getItems(expression string) (items []string) {
	var currentItem string

	for i, char := range expression {
		switch char {
		case ' ', '\t', '\n', '\r', '\v', '\f': // whitespaces
			continue
		case '(', ')', '+', '-', '*', '/': // operations
			if currentItem != "" {
				items = append(items, currentItem)
				currentItem = ""
			}

			// Handle unary minus
			if char == '-' && (i == 0 || expression[i-1] == '(') {
				currentItem += "-"
			} else {
				items = append(items, string(char))
			}
		default:
			currentItem += string(char)
		}
	}

	if currentItem != "" {
		items = append(items, currentItem)
	}

	return items
}

func infixToPostfix(items []string) (postfixExpression []string) {
	var operationsStack stack.Stack

	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}

	for _, item := range items {
		switch item {
		case "+", "-":
			if operationsStack.IsEmpty() || operationsStack.Peek() == "(" {
				precedence["-"] = 3 // Unary minus has higher precedence
			} else {
				precedence["-"] = 1
			}

			for !operationsStack.IsEmpty() && precedence[operationsStack.Peek()] >= precedence[item] {
				postfixExpression = append(postfixExpression, operationsStack.Pop())
			}

			operationsStack.Push(item)
		case "*", "/":
			for !operationsStack.IsEmpty() && precedence[operationsStack.Peek()] >= precedence[item] {
				postfixExpression = append(postfixExpression, operationsStack.Pop())
			}

			operationsStack.Push(item)
		case "(":
			operationsStack.Push(item)
		case ")":
			for !operationsStack.IsEmpty() && operationsStack.Peek() != "(" {
				postfixExpression = append(postfixExpression, operationsStack.Pop())
			}

			if operationsStack.IsEmpty() {
				return nil
			}

			operationsStack.Pop()
		default:
			postfixExpression = append(postfixExpression, item)
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

func calculatePostfix(postfixExpression []string) (float64, error) {
	var exitStack []float64

	for _, item := range postfixExpression {
		switch item {
		case "+", "-", "*", "/":
			if len(exitStack) < 2 {
				return 0, fmt.Errorf("invalid expression")
			}

			b, a := exitStack[len(exitStack)-1], exitStack[len(exitStack)-2]
			exitStack = exitStack[:len(exitStack)-2]

			switch item {
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
			num, err := strconv.ParseFloat(item, 64)

			if err != nil {
				return 0, fmt.Errorf("invalid item: %s", item)
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
	items := getItems(expression)
	postfix := infixToPostfix(items)
	result, err := calculatePostfix(postfix)

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
