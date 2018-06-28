package confexpr

func next_state(key string, state PatternState) (PatternState, error) {
	processors, err := parse_processors(key)
	for _, processor := range processors {
		if len(state.path) == 0 {
			return state, nil
		}
		state, err = processor.next(state)
		if err != nil {
			state.state = Missing
			return state, err
		}
	}
	if len(state.path) == 0 {
		state.state = Complete
	}
	return state, nil
}

//DANGER IMPURE PASS BY REFEENCE
func next_with_collisions(next_pattern, prev_pattern string, prev_state *PatternState, next PatternState) (string, *PatternState) {
	same_class, _ := compare_patterns(prev_pattern, next_pattern)
	if true &&
		(same_class == 0) &&
		prev_state != nil &&
		(next.state == Complete || next.state == Incomplete) &&
		(prev_state.state == Complete || prev_state.state == Incomplete || prev_state.state == Collision) {
		prev_state.state = Collision
		next.state = Collision
	}
	return next_pattern, &next
}
