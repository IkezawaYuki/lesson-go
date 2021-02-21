package path_mirror

import (
	"errors"
	"strings"
)

var ErrBadPattern = errors.New("syntax error in pattern")

func Match(pattern, name string) (matched bool, err error) {
Pattern:
	for len(pattern) > 0 {
		var star bool
		var chunk string
		star, chunk, pattern = scanChunk(pattern)

		if star && chunk == "" {
			return !strings.Contains(name, "/"), nil
		}

		t, ok, err := matchChunk(chunk, name)
	}
}

func matchChunk(chunk, s string) (rest string, ok bool, err error) {
	for len(chunk) > 0 {

	}
}

func scanChunk(pattern string) (star bool, chunk, rest string) {
	for len(pattern) > 0 && pattern[0] == '*' {
		pattern = pattern[1:]
		star = true
	}
	inrange := false
	var i int

Scan:
	for i = 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '\\':
			if i+1 < len(pattern) {
				i++
			}
		case '[':
			inrange = true
		case ']':
			inrange = false
		case '*':
			if !inrange {
				break Scan
			}
		}
	}
	return star, pattern[0:i], pattern[i:]
}
