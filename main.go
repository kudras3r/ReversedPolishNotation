package main

/*
	Calculator based on reversed Polish notation.
	It used Shunting Yard alghorithm (check toPostfix func).
*/

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const operations = "+-*/()"

const (
	Add      = '+'
	Subtract = '-'
	Multiply = '*'
	Divide   = '/'
)

var operationPriority = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
}

func Calc(expression string) (float64, error) {
	var stack []float64
	var tmp float64

	tokens, err := getTokens(expression)
	if err != nil {
		return 0, err
	}

	postfixNotation, err := toPostfix(tokens)
	if err != nil {
		return 0, err
	}

	for _, t := range postfixNotation {
		if n, err := strconv.ParseFloat(t, 64); err == nil {
			stack = append(stack, n)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("Expression is invalid!")
			}

			n1, n2 := stack[len(stack)-2], stack[len(stack)-1]

			switch rune(t[0]) {
			case Add:
				tmp = n1 + n2
			case Subtract:
				tmp = n1 - n2
			case Multiply:
				tmp = n1 * n2
			case Divide:
				if n2 == 0 {
					return 0, errors.New("Division by zero!")
				}
				tmp = n1 / n2
			}

			stack = stack[:len(stack)-2]
			stack = append(stack, tmp)
		}
	}

	return stack[0], nil
}

func getTokens(expression string) ([]string, error) {
	var tokens []string
	var curToken strings.Builder

	for _, c := range expression {
		if c == ' ' {
			continue
		}
		if strings.ContainsRune(operations, c) {
			if curToken.Len() > 0 {
				tokens = append(tokens, curToken.String())
				curToken.Reset()
			}
			tokens = append(tokens, string(c))
		} else if unicode.IsDigit(c) || c == '.' {
			curToken.WriteRune(c)
		} else {
			return nil, errors.New("Invalid input!")
		}
	}

	if curToken.Len() > 0 {
		tokens = append(tokens, curToken.String())
	}

	if len(tokens) == 0 {
		return nil, errors.New("Expression is empty!")
	}

	return tokens, nil
}

func toPostfix(tokens []string) ([]string, error) {
	var res []string
	var stack []rune

	for _, t := range tokens {
		if _, err := strconv.ParseFloat(t, 64); err == nil {
			res = append(res, t)
		} else {
			switch rune(t[0]) {
			case '(':
				stack = append(stack, '(')
			case ')':
				if len(stack) == 0 {
					return nil, errors.New("Incorrect parenthesis!")
				}
				el := stack[len(stack)-1]
				for el != '(' {
					res = append(res, string(el))
					stack = stack[:len(stack)-1]
					if len(stack) == 0 {
						return nil, errors.New("Incorrect parenthesis!")
					}
					el = stack[len(stack)-1]
				}
				stack = stack[:len(stack)-1]
			default:
				for len(stack) > 0 && operationPriority[rune(t[0])] <= operationPriority[stack[len(stack)-1]] {
					res = append(res, string(stack[len(stack)-1]))
					stack = stack[:len(stack)-1]
				}
				stack = append(stack, rune(t[0]))
			}
		}
	}

	for len(stack) > 0 {
		res = append(res, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return res, nil
}

func main() {
	var exp string

	for {
		fmt.Print("Enter the expression: ")
		fmt.Scanln(&exp)

		res, err := Calc(exp)

		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Result:", res)
		}
	}
}
