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
		Path:           []string{},
		Evaluated_path: abc,
		State:          Complete,
		Variables:      map[string]string{},
	}, t)

	state = NewPatternState(abc)
	next, err = next_state("a.b", state)
	expect(err, nil, t)
	expect(next, PatternState{
		Path:           []string{"c"},
		Evaluated_path: []string{"a", "b"},
		State:          Incomplete,
		Variables:      map[string]string{},
	}, t)

	state = NewPatternState(abc)
	next, err = next_state("a.b.c.d", state)
	expect(err, nil, t)
	expect(next, PatternState{
		Path:           []string{},
		Evaluated_path: []string{"a", "b", "c"},
		State:          Incomplete,
		Variables:      map[string]string{},
	}, t)

	state = NewPatternState(abc)
	next, err = next_state("d.e.f", state)
	expectError(err, "Path not found", t)
	expect(next, PatternState{
		Path:           []string{"a", "b", "c"},
		Evaluated_path: []string{},
		State:          Missing,
		Variables:      map[string]string{},
	}, t)

	next, err = next_state("a$var_a.b$var_b.c$var_c", state)
	expect(err, nil, t)
	expect(next, PatternState{
		Path:           []string{},
		Evaluated_path: abc,
		State:          Complete,
		Variables: map[string]string{
			"var_a": "a",
			"var_b": "b",
			"var_c": "c",
		},
	}, t)

}

func TestNextWithCollisions(t *testing.T) {
	ab := []string{"a", "b"}
	//abc := []string{"a", "b", "c"}
	prev := NewPatternState(ab)
	next := NewPatternState(ab)
	prev.State = Complete
	next.State = Complete
	output_pattern, output_state := next_with_collisions("a.b", "a.b", &prev, next)
	expected_state := NewPatternState(ab)
	expected_state.State = Collision
	expect(output_pattern, "a.b", t)
	expect(prev, expected_state, t)
	expect(*output_state, expected_state, t)
}
