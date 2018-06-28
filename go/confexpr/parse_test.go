package confexpr

import (
	"testing"
)

func TestParseProcessorsHappyPaths(t *testing.T) {
	output, err := parse_processors("")
	expect(output, []WrappedLabel{}, t)
	expect(err, nil, t)

	output, err = parse_processors("a")
	expect(output, []WrappedLabel{
		NewWrappedLabel(NewDirectHit("a"), "a", "var"),
	}, t)
	expect(err, nil, t)

	output, err = parse_processors("*")
	expect(output, []WrappedLabel{
		NewWrappedLabel(NewWildcard(), "*", "var"),
	}, t)
	expect(err, nil, t)

	output, err = parse_processors("a.b")
	expect(output, []WrappedLabel{
		NewWrappedLabel(NewDirectHit("a"), "a", "var"),
		NewWrappedLabel(NewDirectHit("b"), "b", "var"),
	}, t)
	expect(err, nil, t)

	output, err = parse_processors("a.b.c")
	expect(output, []WrappedLabel{
		NewWrappedLabel(NewDirectHit("a"), "a", "var"),
		NewWrappedLabel(NewDirectHit("b"), "b", "var"),
		NewWrappedLabel(NewDirectHit("c"), "c", "var"),
	}, t)
	expect(err, nil, t)
}
