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

	for i := 0; i < len(s); i++ {
		sym := s[i]
		if sym == '\'' {
			if quoted {
				quoted = false
			} else {
				quoted = true
			}

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
