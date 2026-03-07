package service

import (
	"strings"
	"unicode"
	
	"go1fl-sprint6-final-tpl/pkg/morse"
)

func isMorseCode(s string) bool {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return false
	}

	words := strings.Split(trimmed, "   ")
	
	for _, word := range words {
		chars := strings.Split(word, " ")
		
		for _, char := range chars {
			if char == "" {
				continue
			}
			
			for _, r := range char {
				if r != '.' && r != '-' {
					return false
				}
			}
		}
	}
	
	return len(words) > 0 && len(words[0]) > 0
}

func isPlainText(s string) bool {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return false
	}
	
	for _, r := range trimmed {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsPunct(r) || unicode.IsSpace(r) {
			continue
		}
		return false
	}
	return true
}

func AutoDetectAndConvert(input string) (string, error) {
	if input == "" {
		return "", nil
	}
	
	if isMorseCode(input) {
		return morse.ToText(input), nil
	}
	
	if isPlainText(input) {
		return morse.ToMorse(input), nil
	}
	
	for _, r := range input {
		if r == '.' || r == '-' {
			return morse.ToText(input), nil
		}
	}
	return morse.ToMorse(input), nil
}

func IsMorseCode(s string) bool {
	return isMorseCode(s)
}