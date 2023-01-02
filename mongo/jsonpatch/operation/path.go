package operation

import (
	"regexp"
	"strings"
)

type Path string

// Valid check if given path is in valid format.
func (p Path) Valid() bool {
	regex := regexp.MustCompile(`^(([\w-]+)+\.?)*([\w-]+)+$`)

	return regex.MatchString(string(p))
}

// Equal check if path is equal to given path.
// Single fields can be set to `*` for wildcard.
func (p Path) Equal(comparePath Path) bool {
	if strings.Contains(string(p), "*") {
		stringPath := []rune(p)
		offset := 0

		for index, char := range comparePath {
			switch {
			case index == len(comparePath)-1 && (index+offset) == len(stringPath)-1:
				return stringPath[index+offset] == '*' || stringPath[index+offset] == char
			case len(stringPath) <= index+offset:
				break
			case stringPath[index+offset] == '*':
				if char != '.' {
					offset--
				} else {
					offset++
				}

				continue
			case stringPath[index+offset] != char:
				return false
			}
		}

		return false
	}

	if p == comparePath {
		return true
	}

	return false
}
