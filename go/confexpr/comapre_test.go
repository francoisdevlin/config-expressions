package confexpr

import (
	"sort"
	"testing"
)

func TestSortStuff(t *testing.T) {
	evaluateSort([]string{"a", "b", "c"}, t)
	evaluateSort([]string{"a", "a.a", "a.b"}, t)
	evaluateSort([]string{"z", "a.a", "a.b"}, t)
	evaluateSort([]string{"z", "*", "a.a", "a.b"}, t)
}

func evaluateSort(expected []string, t *testing.T) {
	actual := make([]string, len(expected), len(expected))
	copy(actual, expected)
	sort.Slice(actual, comparatorLambda(actual))
	expect(actual, expected, t)
}
