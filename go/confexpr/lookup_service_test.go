package confexpr

import (
	"testing"
)

func TestLookupServiceGetSortMatches(t *testing.T) {
	lookupService := NewLookupService()
	rawConfig := map[string]interface{}{
		"a": 1,
		"b": 2,
	}
	state := NewPatternState([]string{})
	state.evaluated_path = []string{"1", "2", "3"}

	expect(len(lookupService.sortCache), 0, t)
	result := lookupService.get_sort_matches(state, rawConfig)
	expect(result, []string{"a", "b"}, t)
	expect(len(lookupService.sortCache), 1, t)

	result = lookupService.get_sort_matches(state, rawConfig)
	expect(result, []string{"a", "b"}, t)
	expect(len(lookupService.sortCache), 1, t)

	state.evaluated_path = []string{"4", "5", "6"}
	result = lookupService.get_sort_matches(state, rawConfig)
	expect(result, []string{"a", "b"}, t)
	expect(len(lookupService.sortCache), 2, t)
}
