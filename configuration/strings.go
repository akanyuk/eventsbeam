package configuration

import (
	"strings"
)

func isLowerCase(rune byte, lowerRune byte) bool {
	return rune == lowerRune
}

func isUpperCase(rune byte, upperRune byte) bool {
	return rune == upperRune
}

func splitByUpperCase(name string) []string {
	var words []string
	firstCharIndex := 0
	upperName := strings.ToUpper(name)
	lowerName := strings.ToLower(name)
	for i := 1; i < len(name); i++ {
		if (isLowerCase(name[i], lowerName[i]) && isUpperCase(name[i-1], upperName[i-1]) ||
			isUpperCase(name[i], upperName[i]) && isLowerCase(name[i-1], lowerName[i-1])) && i-firstCharIndex > 1 {
			words = append(words, name[firstCharIndex:i])
			firstCharIndex = i
		}
	}
	if len(words) == 0 {
		return []string{name}
	}
	words = append(words, name[firstCharIndex:])
	return words
}

func AddDelimiter(sentence string) string {
	sentence = strings.Trim(sentence, " ")
	words := splitByUpperCase(sentence)
	return strings.Join(words, "-")
}
