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

func isQuote(char byte) bool {
	if char == '"' || char == '\'' {
		return true
	}

	return false
}

func isSpecialChar(char byte) bool {
	switch char {
	case ' ', '\'', '"', '$', '*', '?', 'n', 't':
		return true
	}

	return false
}

func argsParse(s string) []string {
	s = strings.TrimSpace(s)

	var currentQuote byte

	var result []string
	var b strings.Builder

	for i := 0; i < len(s); i++ {

		currentSymbol := s[i]

		if currentSymbol == '\\' {
			if isSpecialChar(s[i+1]) {
				currentSymbol = s[i+1]
			} else {
				b.WriteByte(s[i+1])
				i++
				continue
			}
		}

		if isQuote(currentSymbol) && currentQuote == 0 {
			currentQuote = currentSymbol
			continue
		}

		if isQuote(currentSymbol) && currentQuote == currentSymbol {
			currentQuote = 0
			continue
		}

		if currentSymbol == ' ' && currentQuote == 0 {
			if b.Len() > 0 {
				result = append(result, b.String())
				b.Reset()
			}
			continue
		}

		b.WriteByte(currentSymbol)
	}

	if b.Len() > 0 {
		result = append(result, b.String())
	}

	return result
}
