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

func (this ErrorLabel) both(path []string) ([]string, []string, error) {
	return nil, nil, errors.New("Dummy error")
}

func TestWrappedLabelWeights(t *testing.T) {
	var l Label = nil
	l = NewDirectHit("a")
	expect(NewWrappedLabel(l, "a", "var").WEIGHT, 0, t)
	l = NewEnumHit([]string{})
	expect(NewWrappedLabel(l, "a", "var").WEIGHT, 1, t)
	l = NewRegexHit("app")
	expect(NewWrappedLabel(l, "a", "var").WEIGHT, 2, t)
	l = NewWildcard()
	expect(NewWrappedLabel(l, "a", "var").WEIGHT, 3, t)
	l = NewDeepWildcard(nil)
	expect(NewWrappedLabel(l, "a", "var").WEIGHT, 4, t)
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
	expect(next_state.Path, []string{"b", "c"}, t)
	expect(next_state.Evaluated_path, []string{"a"}, t)
	expect(err, nil, t)
}

func TestDirectHitHappyPath(t *testing.T) {
	dh := NewDirectHit("a")
	rest, consumed, err := dh.both([]string{"a", "b", "c"})
	expect(consumed, []string{"a"}, t)
	expect(rest, []string{"b", "c"}, t)
	expect(err, nil, t)

	rest, consumed, err = dh.both([]string{"d"})
	expect(consumed, []string{}, t)
	expect(rest, []string{}, t)
	expectError(err, "Path not found", t)
}

func TestWildcardHappyPath(t *testing.T) {
	wc := NewWildcard()
	rest, consumed, err := wc.both([]string{"a", "b", "c"})
	expect(consumed, []string{"a"}, t)
	expect(rest, []string{"b", "c"}, t)
	expect(err, nil, t)
}

func TestDeepWildcardHappyPath(t *testing.T) {
	dh := NewDirectHit("b")
	dw := NewDeepWildcard(dh)
	rest, consumed, err := dw.both([]string{"a", "b", "c"})
	expect(consumed, []string{"a"}, t)
	expect(rest, []string{"b", "c"}, t)
	expect(err, nil, t)

	rest, consumed, err = dw.both([]string{"a", "b", "b", "b", "c"})
	expect(consumed, []string{"a", "b", "b"}, t)
	expect(rest, []string{"b", "c"}, t)
	expect(err, nil, t)

	rest, consumed, err = dw.both([]string{"a", "c"})
	expect(consumed, []string{}, t)
	expect(rest, []string{"a", "c"}, t)
	expectError(err, "Sucessor path not found", t)

	dw = NewDeepWildcard(NewEnumHit([]string{"b", "c", "d"}))
	rest, consumed, err = dw.both([]string{"a", "b", "c"})
	expect(consumed, []string{"a", "b"}, t)
	expect(rest, []string{"c"}, t)
	expect(err, nil, t)

	rest, consumed, err = dw.both([]string{"a", "e"})
	expect(consumed, []string{}, t)
	expect(rest, []string{"a", "e"}, t)
	expectError(err, "Sucessor path not found", t)
}

func TestRegexHitHappyPath(t *testing.T) {
	rh := NewRegexHit("a")
	rest, consumed, err := rh.both([]string{"a", "b", "c"})
	expect(consumed, []string{"a"}, t)
	expect(rest, []string{"b", "c"}, t)
	expect(err, nil, t)

	rh = NewRegexHit("a.*")
	rest, consumed, err = rh.both([]string{"a", "b", "c"})
	expect(consumed, []string{"a"}, t)
	expect(rest, []string{"b", "c"}, t)
	expect(err, nil, t)

	rh = NewRegexHit("\\w")
	rest, consumed, err = rh.both([]string{"a", "b", "c"})
	expect(consumed, []string{"a"}, t)
	expect(rest, []string{"b", "c"}, t)
	expect(err, nil, t)

	rh = NewRegexHit("\\w+")
	rest, consumed, err = rh.both([]string{"a", "b", "c"})
	expect(consumed, []string{"a"}, t)
	expect(rest, []string{"b", "c"}, t)
	expect(err, nil, t)

	rh = NewRegexHit("d")
	rest, consumed, err = rh.both([]string{"a", "b", "c"})
	expect(consumed, []string{}, t)
	expect(rest, []string{"a", "b", "c"}, t)
	expectError(err, "Path not found", t)
}
