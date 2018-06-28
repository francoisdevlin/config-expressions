package confexpr

import (
	"testing"
)

func TestNextStateHappyPath(t *testing.T) {
	abc := []string{"a", "b", "c"}
	state := NewPatternState(abc)
	next, err := next_state("a.b.c", state)
	expect(err, nil, t)
	expect(next, PatternState{
		path:           []string{},
		evaluated_path: abc,
		state:          Complete,
	}, t)

	state = NewPatternState(abc)
	next, err = next_state("a.b", state)
	expect(err, nil, t)
	expect(next, PatternState{
		path:           []string{"c"},
		evaluated_path: []string{"a", "b"},
		state:          Incomplete,
	}, t)

	state = NewPatternState(abc)
	next, err = next_state("a.b.c.d", state)
	expect(err, nil, t)
	expect(next, PatternState{
		path:           []string{},
		evaluated_path: []string{"a", "b", "c"},
		state:          Incomplete,
	}, t)

	state = NewPatternState(abc)
	next, err = next_state("d.e.f", state)
	expectError(err, "Path not found", t)
	expect(next, PatternState{
		path:           []string{"a", "b", "c"},
		evaluated_path: []string{},
		state:          Missing,
	}, t)
}
