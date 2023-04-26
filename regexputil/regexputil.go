package regexputil

import (
	"regexp"
)

// RegexMapExpression executes a regex returning a map of named captures
func RegexMapExpression(expression *regexp.Regexp, input string) (output *map[string]string) {
	match := expression.FindStringSubmatch(input)
	if match != nil {
		result := make(map[string]string)
		for i, name := range expression.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		output = &result
	}
	return
}

// RegexMapString executes a regex returning a map of named captures
func RegexMapString(expression string, input string) *map[string]string {
	return RegexMapExpression(regexp.MustCompile(expression), input)
}
