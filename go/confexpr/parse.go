package confexpr

import (
	"regexp"
	"strings"
)

func MarshallFromEntry(entry string, prevLabel Label) Label {
	if entry == "**" {
		return NewDeepWildcard(prevLabel)
	} else if entry == "*" {
		return NewWildcard()
	} else if strings.Contains(entry, ",") {
		return NewEnumHit(regexp.MustCompile(",").Split(entry, -1))
	} else if string(entry[0]) == "/" {
		candidateRegex := string(entry[1 : len(entry)-1])
		return NewRegexHit(candidateRegex)
	} else {
		return NewDirectHit(entry)
	}
}

func parse_processors(key string) ([]WrappedLabel, error) {
	labelList := []WrappedLabel{}
	if key == "" {
		return labelList, nil
	}
	entries := regexp.MustCompile("\\.").Split(key, -1)

	var prevLabel Label = nil
	for index := len(entries) - 1; index >= 0; index-- {
		entry := entries[index]
		splitVariables := regexp.MustCompile("\\$").Split(entry, -1)
		entry = splitVariables[0]
		prevLabel = MarshallFromEntry(entry, prevLabel)
		wl := NewWrappedLabel(prevLabel, entry, "")
		if len(splitVariables) == 2 {
			wl.variable = splitVariables[1]
		}
		labelList = append(labelList, wl)
	}

	output := []WrappedLabel{}
	for index := len(labelList) - 1; index >= 0; index-- {
		output = append(output, labelList[index])
	}
	return output, nil
}
