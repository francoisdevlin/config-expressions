package confexpr

func next_state(key string, state PatternState) (PatternState, error) {
	processors, err := parse_processors(key)
	for _, processor := range processors {
		if len(state.Path) == 0 {
			return state, nil
		}
		state, err = processor.next(state)
		if err != nil {
			state.State = Missing
			return state, err
		}
	}
	if len(state.Path) == 0 {
		state.State = Complete
	}
	return state, nil
}

//DANGER IMPURE PASS BY REFEENCE
func next_with_collisions(next_pattern, prev_pattern string, prev_state *PatternState, next PatternState) (string, *PatternState) {
	same_class, _ := compare_patterns(prev_pattern, next_pattern)
	if true &&
		(same_class == 0) &&
		prev_state != nil &&
		(next.State == Complete || next.State == Incomplete) &&
		(prev_state.State == Complete || prev_state.State == Incomplete || prev_state.State == Collision) {
		prev_state.State = Collision
		next.State = Collision
	}
	return next_pattern, &next
}
