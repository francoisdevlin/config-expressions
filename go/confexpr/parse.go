package confexpr

import (
	"regexp"
	"strings"
)

func parse_processors(key string) ([]WrappedLabel, error) {
	output := []WrappedLabel{}
	if key == "" {
		return output, nil
	}
	entries := regexp.MustCompile("\\.").Split(key, -1)

	for _, entry := range entries {
		splitVariables := regexp.MustCompile("\\$").Split(entry, -1)
		entry = splitVariables[0]
		var l Label = nil
		if entry == "*" {
			l = NewWildcard()
		} else if strings.Contains(entry, ",") {
			l = NewEnumHit(regexp.MustCompile(",").Split(entry, -1))
		} else if string(entry[0]) == "/" {
			candidateRegex := string(entry[1 : len(entry)-1])
			l = NewRegexHit(candidateRegex)
		} else {
			l = NewDirectHit(entry)
		}
		wl := NewWrappedLabel(l, entry, "")
		if len(splitVariables) == 2 {
			wl.variable = splitVariables[1]
		}
		output = append(output, wl)
	}
	return output, nil
}
