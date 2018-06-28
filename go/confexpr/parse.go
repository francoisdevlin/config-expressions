package confexpr

import "regexp"

func parse_processors(key string) ([]WrappedLabel, error) {
	output := []WrappedLabel{}
	if key == "" {
		return output, nil
	}
	entries := regexp.MustCompile("\\.").Split(key, -1)

	for _, entry := range entries {
		var l Label = nil
		if entry == "*" {
			l = NewWildcard()
		} else {
			l = NewDirectHit(entry)
		}
		wl := NewWrappedLabel(l, entry, "var")
		output = append(output, wl)
	}
	return output, nil
}
