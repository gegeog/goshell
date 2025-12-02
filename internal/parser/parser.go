package parser

import (
	"strings"
)

func Parse(input string) (string, []string) {
	input = strings.TrimSpace(input)
	sep := strings.SplitN(input, " ", 2)

	if len(sep) == 0 {
		return "", nil
	}

	if len(sep) == 1 {
		return sep[0], nil
	}

	return sep[0], argsParse(sep[1])
}

func isQuote(char rune) bool {
	if char == '"' || char == '\'' {
		return true
	}

	return false
}

func argsParse(s string) []string {
	s = strings.TrimSpace(s)

	var currentQuote rune
	var isQuoted bool

	var result []string
	var b strings.Builder

	// 'asdfadsfa""'
	for _, v := range s {
		if isQuote(v) && currentQuote == 0 {
			currentQuote = v
			isQuoted = true
			continue
		}

		if isQuote(v) && v == currentQuote {
			currentQuote = 0
			isQuoted = false
			continue
		}

		if v == ' ' && !isQuoted {
			if b.Len() > 0 {
				result = append(result, b.String())
				b.Reset()
			}
			continue
		}

		b.WriteRune(v)
	}

	if b.Len() > 0 {
		result = append(result, b.String())
	}

	return result
}
