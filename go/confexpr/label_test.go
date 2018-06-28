package confexpr

import (
	"errors"
	"reflect"
	"testing"
)

type ErrorLabel struct{}

func expect(actual, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v', Actual '%v'", expected, actual)
	}
}

func expectError(actual error, message string, t *testing.T) {
	if actual == nil {
		t.Errorf("Error was expected, none found")
	}
	expect(actual.Error(), message, t)
}

func (this ErrorLabel) both(path []string) ([]string, string, error) {
	return nil, "", errors.New("Dummy error")
}

func TestWrappedLabelWeights(t *testing.T) {
	var l Label = nil
	l = NewDirectHit("a")
	expect(NewWrappedLabel(l, "a", "var").WEIGHT, 0, t)
	l = NewEnumHit()
	expect(NewWrappedLabel(l, "a", "var").WEIGHT, 1, t)
	l = NewRegexHit()
	expect(NewWrappedLabel(l, "a", "var").WEIGHT, 2, t)
	l = NewWildcard()
	expect(NewWrappedLabel(l, "a", "var").WEIGHT, 3, t)
}

func TestWrappedLabelErrorsArePropogatedUp(t *testing.T) {
	wl := NewWrappedLabel(ErrorLabel{}, "", "")
	state := NewPatternState([]string{})
	_, err := wl.next(state)
	expectError(err, "Dummy error", t)
}

func TestWrappedLabelHappyPath(t *testing.T) {
	wl := NewWrappedLabel(NewDirectHit("a"), "", "")
	state := NewPatternState([]string{"a", "b", "c"})
	next_state, err := wl.next(state)
	expect(next_state.path, []string{"b", "c"}, t)
	expect(next_state.evaluated_path, []string{"a"}, t)
	expect(err, nil, t)
}

func TestDirectHitHappyPath(t *testing.T) {
	dh := NewDirectHit("a")
	rest, consumed, err := dh.both([]string{"a", "b", "c"})
	expect(consumed, "a", t)
	expect(rest, []string{"b", "c"}, t)
	expect(err, nil, t)

	rest, consumed, err = dh.both([]string{"d"})
	expect(consumed, "", t)
	expect(rest, []string{}, t)
	expectError(err, "Path not found", t)
}

func TestWildcardHappyPath(t *testing.T) {
	wc := NewWildcard()
	rest, consumed, err := wc.both([]string{"a", "b", "c"})
	expect(consumed, "a", t)
	expect(rest, []string{"b", "c"}, t)
	expect(err, nil, t)
}
