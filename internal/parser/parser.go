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

	return sep[0], validate(sep[1])
}

func validate(s string) string {
	var b strings.Builder
	var quoted bool

	for i, v := range s {
		if string(v) == "'" {
			if quoted {
				quoted = false
				continue
			}

			quoted = true
			continue
		}

		if quoted {
			b.WriteRune(v)
			continue
		}

		if string(v) == " " {
			if i-1 >= 0 {
				b.WriteRune(v)
				continue
			}
		}

		b.WriteRune(v)
	}

	return b.String()
}
