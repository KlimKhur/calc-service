package calc

import (
	"errors"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	var isValid = func(expression string) bool {
		for _, elem := range expression {
			if !(elem >= '0' && elem <= '9' || elem == '+' || elem == '-' || elem == '*' || elem == '/' || elem == '(' || elem == ')') {
				return false
			}
		}
		return true
	}(expression)
	if !isValid {
		return 0, errors.New("Invalid expression 1")
	}

	postfix, err := toPostfix(expression)

	if err != nil {
		return 0, err
	}

	result, err := Solution(postfix)

	if err != nil {
		return 0, err
	}

	return result, nil
}

// Приведение выражения в постфиксному виду
func toPostfix(expression string) (string, error) {

	var stack []rune
	var priority = map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	var newExpression strings.Builder
	for _, elem := range expression {
		if elem >= '0' && elem <= '9' {
			newExpression.WriteRune(elem)
		} else if elem == '(' {
			stack = append(stack, elem)
		} else if elem == ')' {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				newExpression.WriteRune(stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return "", errors.New("Error brackets 1")
			}
			stack = stack[:len(stack)-1]
		} else {
			for len(stack) > 0 && priority[elem] <= priority[stack[len(stack)-1]] {
				newExpression.WriteRune(stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, elem)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == '(' {
			return "", errors.New("Error brackets 2")
		}
		newExpression.WriteRune(stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return newExpression.String(), nil

}

// Решение постфиксного выражения
func Solution(expression string) (float64, error) {
	var stack []float64

	for _, elem := range expression {
		if elem >= '0' && elem <= '9' {
			num, _ := strconv.ParseFloat(string(elem), 64)
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("Invalid expression 2")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch elem {
			case '+':
				stack = append(stack, a+b)
			case '-':
				stack = append(stack, a-b)
			case '*':
				stack = append(stack, a*b)
			case '/':
				if b == 0 {
					return 0, errors.New("Division by zero 1")
				}
				stack = append(stack, a/b)
			}
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("Invalid expression 3")
	}

	return stack[0], nil
}
