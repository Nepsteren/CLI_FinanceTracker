package main

import (
	"fmt"
	"strings"
)

func parseInput(input string) ([]string, error) {
	var args []string
	var current strings.Builder

	var quote rune
	escaped := false
	tokenStarted := false

	for _, r := range input {
		if escaped {
			current.WriteRune(r)
			escaped = false
			tokenStarted = true
			continue
		}

		if r == '\\' && quote != '\'' {
			escaped = true
			tokenStarted = true
			continue
		}

		if quote != 0 {
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
			continue
		}

		switch r {
		case '"', '\'':
			quote = r
			tokenStarted = true

		case ' ', '\t':
			if tokenStarted {
				args = append(args, current.String())
				current.Reset()
				tokenStarted = false
			}

		default:
			current.WriteRune(r)
			tokenStarted = true
		}
	}
	if escaped {
		return nil, fmt.Errorf("незавершённая escape-последовательность")
	}

	if quote != 0 {
		return nil, fmt.Errorf("незакрытые кавычки")
	}

	if tokenStarted {
		args = append(args, current.String())
	}

	return args, nil
}
