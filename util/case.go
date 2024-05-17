package util

import (
	"strings"
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
