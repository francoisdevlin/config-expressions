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
		NewWrappedLabel(NewDirectHit("a"), "a", ""),
	}, t)
	expect(err, nil, t)

	output, err = parse_processors("*")
	expect(output, []WrappedLabel{
		NewWrappedLabel(NewWildcard(), "*", ""),
	}, t)
	expect(err, nil, t)

	output, err = parse_processors("/abc/")
	expect(output, []WrappedLabel{
		NewWrappedLabel(NewRegexHit("abc"), "/abc/", ""),
	}, t)
	expect(err, nil, t)

	output, err = parse_processors("a.b")
	expect(output, []WrappedLabel{
		NewWrappedLabel(NewDirectHit("a"), "a", ""),
		NewWrappedLabel(NewDirectHit("b"), "b", ""),
	}, t)
	expect(err, nil, t)

	output, err = parse_processors("a.b.c")
	expect(output, []WrappedLabel{
		NewWrappedLabel(NewDirectHit("a"), "a", ""),
		NewWrappedLabel(NewDirectHit("b"), "b", ""),
		NewWrappedLabel(NewDirectHit("c"), "c", ""),
	}, t)
	expect(err, nil, t)

	//Test Variable Parsing
	output, err = parse_processors("a$var_a.*$var_wild")
	expect(output, []WrappedLabel{
		NewWrappedLabel(NewDirectHit("a"), "a", "var_a"),
		NewWrappedLabel(NewWildcard(), "*", "var_wild"),
	}, t)
	expect(err, nil, t)
}
