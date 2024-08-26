package parse

import (
	"fmt"
	"strings"
	"unicode"
)

type tokenType int

const (
	Start tokenType = iota
	End
	Plus
	Separator
	Colon
	Minus
	Num
	Nil
)

type token struct {
	Type  tokenType
	Value int
}

func lexer(input string) ([]token, error) {
	var tokens []token

	patterns := []struct {
		str   string
		token tokenType
	}{
		{"<span class=\"table-data file-modified\" title=\"", Start},
		{"</span>", End},
		{"\">", Separator},
		{":", Colon},
		{"+", Plus},
		{"-", Minus},
	}

	i := 0
	for i < len(input) {
		matched := false

		for _, pattern := range patterns {
			if strings.HasPrefix(input[i:], pattern.str) {
				tokens = append(tokens, token{Type: pattern.token})
				i += len(pattern.str)
				matched = true
				break
			}
		}

		if matched {
			continue
		}

		if num, len := parseInt(input[i:]); len > 0 {
			tokens = append(tokens, token{Type: Num, Value: num})
			i += len
			continue
		}

		if unicode.IsSpace(rune(input[i])) {
			i++
			continue
		}

		return nil, fmt.Errorf("unexpected character")
	}

	return tokens, nil
}
