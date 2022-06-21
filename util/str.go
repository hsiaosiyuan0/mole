package util

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// remove comments in json, for keep the loc of error is meaningful, the comments
// are replace to their equal-length spaces
func RemoveJsonComments(str string) ([]byte, error) {
	sb := strings.Builder{}

	// 0 - normal
	// 1 - string
	// 2 - singleline comment
	// 3 - multiline comment
	state := 0

	pr := utf8.RuneError

	peek := func(str string) (rune, int, error) {
		r, s := utf8.DecodeRuneInString(str)
		if r == utf8.RuneError && s > 0 {
			return utf8.RuneError, 0, errors.New("deformed json")
		}
		return r, s, nil
	}
	next := func() rune {
		r, s := utf8.DecodeRuneInString(str)
		str = str[s:]
		return r
	}

	for len(str) > 0 {
		r, s, err := peek(str)
		if err != nil {
			return nil, err
		}

		switch r {
		case '"':
			if state == 0 {
				state = 1
			} else if state == 1 && pr != '\\' {
				state = 0
			}
		case '/':
			if state == 0 {
				pp, _, _ := peek(str[s:])
				if pp == '/' {
					state = 2
				} else if pp == '*' {
					state = 3
				}
				if state == 2 || state == 3 {
					next()
					next()
					sb.WriteRune(' ')
					sb.WriteRune(' ')
					continue
				}
			}
		case '\r':
			if state == 2 {
				pp, _, _ := peek(str[s:])
				if pp == '\n' {
					next()
				}
				next()
				sb.WriteRune('\n')
				state = 0
				continue
			}
		case '\n':
			if state == 2 {
				next()
				sb.WriteRune('\n')
				state = 0
				continue
			}
		case '*':
			if state == 3 {
				pp, _, _ := peek(str[s:])
				if pp == '/' {
					next()
					next()
					sb.WriteRune(' ')
					sb.WriteRune(' ')
					state = 0
					continue
				}
			}
		}

		if state == 0 || state == 1 {
			sb.WriteRune(next())
		} else {
			next()
			if state == 3 && (r == '\r' || r == '\n') {
				sb.WriteRune('\n')
			}
			sb.WriteRune(' ')
		}
		pr = r
	}
	return []byte(sb.String()), nil
}
