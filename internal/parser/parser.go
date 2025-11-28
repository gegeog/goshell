package parser

import (
	"strings"
)

func Parse(input string) (string, string) {
	input = strings.TrimSpace(input)
	sep := strings.SplitN(input, " ", 2)

	if len(sep) == 0 {
		return "", ""
	}

	if len(sep) == 1 {
		return sep[0], ""
	}

	return sep[0], sep[1]
}

func EchoParse(s string) string {
	var b strings.Builder
	var quoted bool

	for i := 0; i < len(s); i++ {
		sym := s[i]
		if sym == '\'' {
			quoted = !quoted
			continue
		}

		if quoted {
			b.WriteByte(sym)
			continue
		}

		if sym == ' ' {
			if i-1 >= 0 && s[i-1] != ' ' {
				b.WriteByte(sym)
			}

			continue
		}

		b.WriteByte(sym)
	}

	return b.String()
}

func ArgsParse(s string) []string {
	args := strings.Split(s, "'")
	return filterEmpty(args)
}

func filterEmpty(s []string) []string {
	var res []string
	for _, v := range s {
		if len(strings.TrimSpace(v)) > 0 {
			res = append(res, v)
		}
	}

	return res
}
