package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Calc(exp string) (float64, error) {
	if exp == "" {
		return 0, fmt.Errorf("empty")
	}
	var stack []string
	pattern, err := convert_to_reverse_pollish_notation(exp)
	if err != nil {
		return 0, err
	}
	stack = []string{pattern[0]}
	for _, v := range pattern[1:] {
		if unicode.IsDigit([]rune(v)[0]) {
			stack = append(stack, v)
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("Error")
			}
			if v == "+" {
				stack = append(stack[:len(stack)-2], sum(stack[len(stack)-2:]))
			} else if v == "-" {
				stack = append(stack[:len(stack)-2], differance(stack[len(stack)-2:]))
			} else if v == "*" {
				stack = append(stack[:len(stack)-2], multiply(stack[len(stack)-2:]))
			} else if v == "/" {
				t, err := division(stack[len(stack)-2:])
				if err != nil {
					return -1, err
				}
				stack = append(stack[:len(stack)-2], t)
			}
		}
	}
	t, _ := strconv.ParseFloat(stack[0], 64)
	return t, nil
}

func sum(numbers []string) string {
	var sm float64
	for _, v := range numbers {
		t, _ := strconv.ParseFloat(v, 64)
		sm += t
	}
	return strconv.FormatFloat(sm, 'g', -1, 64)
}

func differance(numbers []string) string {
	sm, _ := strconv.ParseFloat(numbers[0], 64)
	for _, v := range numbers[1:] {
		t, _ := strconv.ParseFloat(v, 64)
		sm -= t
	}
	return strconv.FormatFloat(sm, 'g', -1, 64)
}

func multiply(numbers []string) string {
	var sm float64 = 1
	for _, v := range numbers {
		t, _ := strconv.ParseFloat(v, 64)
		sm *= t
	}
	return strconv.FormatFloat(sm, 'g', -1, 64)
}

func division(numbers []string) (string, error) {
	sm, _ := strconv.ParseFloat(numbers[0], 64)
	for _, v := range numbers[1:] {
		t, _ := strconv.ParseFloat(v, 64)
		if t == 0 {
			return "", fmt.Errorf("division by zero")
		}
		sm /= t
	}
	return strconv.FormatFloat(sm, 'g', -1, 64), nil
}

func convert_to_reverse_pollish_notation(pattern string) ([]string, error) {
	var lst []string
	var stack, ouput_stack []string
	operators := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}
	pattern = strings.ReplaceAll(pattern, "+", " + ")
	pattern = strings.ReplaceAll(pattern, "-", " - ")
	pattern = strings.ReplaceAll(pattern, "*", " * ")
	pattern = strings.ReplaceAll(pattern, "/", " / ")
	pattern = strings.ReplaceAll(pattern, "(", " ( ")
	pattern = strings.ReplaceAll(pattern, ")", " ) ")
	lst = strings.Fields(pattern)
	for _, v := range lst {
		if unicode.IsDigit([]rune(v)[0]) {
			ouput_stack = append(ouput_stack, v)
		} else if precedence, exists := operators[v]; exists {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if topPrecedence, topExists := operators[top]; topExists && topPrecedence >= precedence {
					ouput_stack = append(ouput_stack, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, v)
		} else if v == "(" {
			stack = append(stack, v)
		} else if v == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				ouput_stack = append(ouput_stack, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("too many closing parentheses")
			}
			stack = stack[:len(stack)-1]
		}
	}
	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, fmt.Errorf("too many opening parentheses")
		}
		ouput_stack = append(ouput_stack, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return ouput_stack, nil
}

func main() {
	fmt.Println(Calc("((7+1)/(2+2)*4)/8*(32-((4+12)*2))-1"))
}
