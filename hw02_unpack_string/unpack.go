package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	var result strings.Builder
	var prevRune rune
	isEscaped := false

	for _, r := range input {
		if isDigits(r) {
			if err := handleDigits(r, &result, prevRune, &isEscaped); err != nil {
				return "", err
			}
		} else {
			if r == '\\' && !isEscaped {
				isEscaped = true
			} else {
				result.WriteRune(r)
				isEscaped = false
			}
		}
		prevRune = r
	}
	if isEscaped {
		return "", ErrInvalidString
	}

	return result.String(), nil
}

func handleDigits(r rune, result *strings.Builder, prevRune rune, isEscaped *bool) error {
	if isDigits(prevRune) {
		return ErrInvalidString
	}

	if prevRune == 0 || *isEscaped {
		return ErrInvalidString
	}

	count, _ := strconv.Atoi(string(r))
	if count == 0 {
		truncateLastRune(prevRune, result)
	} else {
		for j := 0; j < count-1; j++ {
			result.WriteRune(prevRune)
		}
	}

	*isEscaped = false

	return nil
}

func isDigits(r rune) bool {
	match, _ := regexp.MatchString(`^[0-9]+$`, string(r))
	return match
}

func truncateLastRune(prevRune rune, sb *strings.Builder) {
	s := sb.String()
	_, size := utf8.DecodeLastRuneInString(string(prevRune))
	truncateStr := s[:len(s)-size]
	sb.Reset()
	sb.WriteString(truncateStr)
}
