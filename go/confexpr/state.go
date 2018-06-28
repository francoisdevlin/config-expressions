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
