package utils

import (
	"regexp"
)

func ParseUrlPath(rgx *regexp.Regexp, path string) map[string]string {
	match := rgx.FindStringSubmatch(path)

	params := make(map[string]string)
	for i, name := range rgx.SubexpNames() {
		if i > 0 && i <= len(match) && match[i] != "" {
			params[name] = match[i]
		}
	}
	return params
}
