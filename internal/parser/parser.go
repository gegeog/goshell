package parser

import (
	"strconv"
	"strings"
)

const (
	OutputRedirect = 1
	ErrorRedirect  = 2
)

func Parse(input string) (string, []string, []string, []string) {
	input = strings.TrimSpace(input)

	if len(input) == 0 {
		return "", nil, nil, nil
	}

	var separatedArgs []string
	var outputRedirectPaths []string
	var errorRedirectPaths []string

	if isQuote(input[0]) {
		separatedArgs, outputRedirectPaths, errorRedirectPaths = argsParse(input)
	} else {
		separatedArgs = strings.SplitN(input, " ", 2)
	}

	if len(separatedArgs) == 1 {
		return separatedArgs[0], nil, nil, nil
	}

	if isQuote(input[0]) {
		return separatedArgs[0], separatedArgs[1:], outputRedirectPaths, errorRedirectPaths
	}

	arguments, outputRedirectPaths, errorRedirectPaths := argsParse(separatedArgs[1])
	return separatedArgs[0], arguments, outputRedirectPaths, errorRedirectPaths
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

func argsParse(s string) (args, outRedirectPaths, errRedirectPaths []string) {
	s = strings.TrimSpace(s)

	var lastChar byte
	var currentQuote byte

	var isReadingRedirect bool
	var redirectMode int

	var b strings.Builder

	for i := 0; i < len(s); i++ {
		//QUOTES IN/OUT...
		if isQuote(s[i]) && currentQuote == 0 {
			currentQuote = s[i]
			continue
		}

		if isQuote(s[i]) && currentQuote == s[i] {
			currentQuote = 0
			continue
		}
		//...QUOTES IN/OUT

		//IN QUOTES LIMITS...
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
		//...IN QUOTES LIMITS

		//REDIRECT IN/OUT...
		if currentQuote == 0 && s[i] == '>' {
			redirectMode = OutputRedirect

			if b.Len() > 0 {
				if s[i-1] == '1' || s[i-1] == '2' {
					redirectMode, _ = strconv.Atoi(string(s[i-1]))
					if b.Len() > 1 {
						args = append(args, b.String()[:len(b.String())-1])
					}
				} else {
					args = append(args, b.String())
				}
				b.Reset()
			}

			isReadingRedirect = true
			continue
		}

		if isReadingRedirect && s[i] == ' ' && lastChar != 0 {
			outputFilePath := b.String()
			if redirectMode == OutputRedirect {
				outRedirectPaths = append(outRedirectPaths, outputFilePath)
			} else {
				errRedirectPaths = append(errRedirectPaths, outputFilePath)
			}
			lastChar = 0
			isReadingRedirect = false
			b.Reset()
			continue
		}
		//...REDIRECT IN/OUT

		//REDIRECT PARSING...
		if isReadingRedirect {
			if s[i] != ' ' {
				lastChar = s[i]
				b.WriteByte(s[i])
			}
			continue
		}
		//...REDIRECT PARSING

		//SPECIAL SYMBOLS...
		if currentQuote == 0 && s[i] == '\\' {
			b.WriteByte(s[i+1])
			i++
			continue
		}
		//...SPECIAL SYMBOLS

		//SKIPPING EMPTY SPACE...
		if s[i] == ' ' && currentQuote == 0 {
			if b.Len() > 0 {
				args = append(args, b.String())
				b.Reset()
			}
			continue
		}
		//...SKIPPING EMPTY SPACE

		b.WriteByte(s[i])
	}

	if b.Len() > 0 {
		if isReadingRedirect {
			filePath := b.String()
			if redirectMode == OutputRedirect {
				outRedirectPaths = append(outRedirectPaths, filePath)
			} else {
				errRedirectPaths = append(errRedirectPaths, filePath)
			}
		} else {
			args = append(args, b.String())
		}
	}

	return
}
