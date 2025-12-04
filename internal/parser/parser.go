package parser

import (
	"strings"
)

func Parse(input string) (string, []string, string) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return "", nil, ""
	}

	var sep []string
	var output string
	if isQuote(input[0]) {
		sep, output = argsParse(input)
	} else {
		sep = strings.SplitN(input, " ", 2)
	}

	if len(sep) == 1 {
		return sep[0], nil, ""
	}

	if isQuote(input[0]) {
		return sep[0], sep[1:], output
	}

	arguments, output := argsParse(sep[1])
	return sep[0], arguments, output
}

func isQuote(char byte) bool {
	if char == '"' || char == '\'' {
		return true
	}

	return false
}

func isSpecialChar(char byte) bool {
	switch char {
	case '"', '\\':
		return true
	}

	return false
}

func argsParse(s string) ([]string, string) {
	s = strings.TrimSpace(s)

	var currentQuote byte

	var isReadingRedirect bool
	var output string
	var lastChar byte

	var result []string
	var b strings.Builder

	for i := 0; i < len(s); i++ {
		if isQuote(s[i]) && currentQuote == 0 {
			currentQuote = s[i]
			continue
		}

		if isQuote(s[i]) && currentQuote == s[i] {
			currentQuote = 0
			continue
		}

		//READING REDIRECT...
		if currentQuote == 0 && s[i] == '>' {
			if b.Len() > 0 {
				if s[i-1] == '1' {
					result = append(result, b.String()[:len(b.String())-1])
				} else {
					result = append(result, b.String())
				}
				b.Reset()
			}

			isReadingRedirect = true
			continue
		}

		if isReadingRedirect && s[i] == ' ' && lastChar != 0 {
			output = b.String()
			b.Reset()
			isReadingRedirect = false
			continue
		}

		if isReadingRedirect {
			if s[i] != ' ' {
				lastChar = s[i]
				b.WriteByte(s[i])
			}
			continue
		}
		//...READING REDIRECT

		if currentQuote == 0 && s[i] == '\\' {
			b.WriteByte(s[i+1])
			i++
			continue
		}

		if currentQuote == '"' && s[i] == '\\' {
			nextIndex := i + 1
			if nextIndex > len(s)-1 {
				break
			}

			if isSpecialChar(s[nextIndex]) {
				b.WriteByte(s[nextIndex])
				i++
				continue
			}
		}

		if s[i] == ' ' && currentQuote == 0 {
			if b.Len() > 0 {
				result = append(result, b.String())
				b.Reset()
			}
			continue
		}

		b.WriteByte(s[i])
	}

	if b.Len() > 0 {
		if isReadingRedirect {
			output = b.String()
		} else {
			result = append(result, b.String())
		}
	}

	return result, output
}
