package parser

import (
	"strconv"
	"strings"
)

const (
	OutputRedirect = 1
	ErrorRedirect  = 2
)

type ParsedInfo struct {
	Arguments           []string
	OutputRedirectsRest []string
	OutputRedirectsNew  []string
	ErrRedirectRest     []string
	ErrRedirectNew      []string
}

func Parse(input string) (string, ParsedInfo) {
	input = strings.TrimSpace(input)

	var (
		pi            ParsedInfo
		separatedArgs []string
	)

	if len(input) == 0 {
		return "", pi
	}

	if isQuote(input[0]) {
		pi = argsParse(input)
	} else {
		separatedArgs = strings.SplitN(input, " ", 2)
	}

	if len(separatedArgs) == 1 {
		return separatedArgs[0], pi
	}

	if isQuote(input[0]) {
		cmd := pi.Arguments[0]
		pi.Arguments = pi.Arguments[1:]
		return cmd, pi
	}

	pi = argsParse(separatedArgs[1])
	return separatedArgs[0], pi
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

func argsParse(s string) ParsedInfo {
	s = strings.TrimSpace(s)

	var (
		args                                       []string
		outRedirectPaths, errRedirectPaths         []string
		outRedirectPathsRest, errRedirectPathsRest []string

		lastChar          byte
		currentQuote      byte
		isReadingRedirect bool
		isRedirectRest    bool
		redirectMode      int

		b strings.Builder
	)

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

			if i+1 <= len(s)-1 && s[i+1] == '>' {
				isRedirectRest = true
				i++
			}

			continue
		}

		if isReadingRedirect && s[i] == ' ' && lastChar != 0 {
			outputFilePath := b.String()
			if redirectMode == OutputRedirect {
				if !isRedirectRest {
					outRedirectPaths = append(outRedirectPaths, outputFilePath)
				} else {
					outRedirectPathsRest = append(outRedirectPathsRest, outputFilePath)
				}
			} else {
				if !isRedirectRest {
					errRedirectPaths = append(errRedirectPaths, outputFilePath)
				} else {
					errRedirectPathsRest = append(errRedirectPathsRest, outputFilePath)
				}
			}
			lastChar = 0
			isReadingRedirect = false
			isRedirectRest = false
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
				if !isRedirectRest {
					outRedirectPaths = append(outRedirectPaths, filePath)
				} else {
					outRedirectPathsRest = append(outRedirectPathsRest, filePath)
				}
			} else {
				if !isRedirectRest {
					errRedirectPaths = append(errRedirectPaths, filePath)
				} else {
					errRedirectPathsRest = append(errRedirectPathsRest, filePath)
				}
			}
		} else {
			args = append(args, b.String())
		}
	}

	return ParsedInfo{
		Arguments:           args,
		OutputRedirectsNew:  outRedirectPaths,
		OutputRedirectsRest: outRedirectPathsRest,
		ErrRedirectNew:      errRedirectPaths,
		ErrRedirectRest:     errRedirectPathsRest,
	}
}
