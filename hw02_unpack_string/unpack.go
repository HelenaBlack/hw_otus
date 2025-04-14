package main

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}
	var result []rune
	var preRune rune
	escaped := false

	for _, r := range input {
		if unicode.IsDigit(preRune) {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(r) {
			if preRune == 0 || escaped {
				return "", ErrInvalidString
			}

			count, _ := strconv.Atoi(string(r))
			if count == 0 {
				result = result[:len(result)-1]
			} else {
				for j := 0; j < count-1; j++ {
					result = append(result, preRune)
				}
			}

			escaped = false
		} else {
			if r == '\\' && !escaped {
				escaped = true
			} else {
				result = append(result, r)
				escaped = false
			}
		}
		preRune = r
	}
	if escaped {
		return "", ErrInvalidString
	}
	return string(result), nil
}
