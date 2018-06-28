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

func TestNextWithCollisions(t *testing.T) {
	ab := []string{"a", "b"}
	//abc := []string{"a", "b", "c"}
	prev := NewPatternState(ab)
	next := NewPatternState(ab)
	prev.state = Complete
	next.state = Complete
	output_pattern, output_state := next_with_collisions("a.b", "a.b", &prev, next)
	expected_state := NewPatternState(ab)
	expected_state.state = Collision
	expect(output_pattern, "a.b", t)
	expect(prev, expected_state, t)
	expect(*output_state, expected_state, t)
}
