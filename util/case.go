package util

import (
	"strings"
	"unicode"
)

func ConvertToSnakeCase(s string) string {
	var result string
	for i, c := range s {
		if i > 0 && c >= 'A' && c <= 'Z' && s[i-1] != '/' && s[i-1] != '_' {
			result += "_"
		}

		str := string(c)
		if str != " " {
			result += str
		}
	}

	return strings.ToLower(result)
}

func ConvertToUpperCamelCase(s string) string {
	var result string
	for i, c := range s {
		if i == 0 {
			result += strings.ToUpper(string(c))
		} else if i > 0 && s[i-1] == '_' {
			result += strings.ToUpper(string(c))
		} else {
			result += string(c)
		}
	}
	result = strings.ReplaceAll(result, "_", "")

	return result
}

func CamelToSnake(s string) string {
	var result []rune

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && unicode.IsLower(rune(s[i-1])) {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}

func SnakeToCamel(s string) string {
	var result []rune

	words := strings.Split(s, "_")
	for i, word := range words {
		if i > 0 {
			result = append(result, ' ')
		}
		for j, r := range word {
			if j == 0 {
				result = append(result, unicode.ToUpper(r))
			} else {
				result = append(result, r)
			}
		}
	}

	return string(result)
}
