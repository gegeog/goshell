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
	s = strings.TrimSpace(s)

	var quoted bool
	var prev string

	var res []string
	var b strings.Builder

	for _, v := range s {
		sym := string(v)

		if quoted && sym != "'" {
			b.WriteString(sym)
			continue
		}

		if quoted {
			quoted = false
			res = append(res, b.String())
			b.Reset()
			continue
		}

		if sym == "'" {
			quoted = true
			res = append(res, filterSpaces(b.String())...)
			b.Reset()
			continue
		}

		if sym == " " && prev != " " {
			res = append(res, filterSpaces(b.String())...)
			b.Reset()
			prev = sym
			continue
		}

		b.WriteString(sym)
		prev = sym
	}

	res = append(res, b.String())

	return res
}

func filterSpaces(s string) []string {
	var res []string
	var sl = strings.Split(s, " ")

	for _, v := range sl {
		if trimmed := strings.TrimSpace(v); len(trimmed) > 0 {
			res = append(res, trimmed)
		}
	}

	return res
}
