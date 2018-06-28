package confexpr

import "fmt"
import "errors"
import "sort"

func next_state(key string, state PatternState) (PatternState, error) {
	processors, _ := parse_processors(key)
	fmt.Println(processors)
	fmt.Println(state)
	for _, processor := range processors {
		if len(state.path) == 0 {
			return state, nil
		}
		state, err := processor.next(state)
		fmt.Println(state)
		if err != nil {
			state.state = Missing
			return state, nil
		}
	}
	if len(state.path) == 0 {
		state.state = Complete
	}
	return state, nil
}

type Result struct {
	Key   string
	State PatternState
}

func determine_match_states(start PatternState, rawConfig map[string]interface{}) ([]Result, error) {
	previous_iteration_pattern := ""
	//_ = previous_interation_pattern
	var prev_state *PatternState = nil
	//sorted_matches \c := rawConfig
	sorted_matches := []string{}
	for key, _ := range rawConfig {
		sorted_matches = append(sorted_matches, key)
	}

	sort.Slice(sorted_matches, comparatorLambda(sorted_matches))
	//fmt.Println(sorted_matches)

	output := []Result{}
	for _, match := range sorted_matches {
		same_class, _ := compare_patterns(previous_iteration_pattern, match)
		next_state, _ := next_state(match, start)
		if true &&
			(same_class == 0) &&
			prev_state != nil &&
			(next_state.state == Complete || next_state.state == Incomplete) &&
			(prev_state.state == Complete || prev_state.state == Incomplete || prev_state.state == Collision) {
			prev_state.state = Collision
			next_state.state = Collision
		}
		previous_iteration_pattern = match
		prev_state = &next_state
		output = append(output, Result{match, next_state})
	}
	return output, nil
}

func Lookup(state PatternState, rawConfig map[string]interface{}) ([]Result, error) {
	results, err := determine_match_states(state, rawConfig)
	if err != nil {
		return results, err
	}
	output := []Result{}
	for _, result := range results {
		state := result.State
		match := result.Key
		if state.state == Complete {
			state.value = rawConfig[match]
			output = append(output, Result{match, state})
		} else if state.state == Incomplete && len(state.path) > 0 {
			next_value := rawConfig[match]
			if map_next_value, ok := next_value.(map[string]interface{}); ok {
				temp, err := Lookup(state, map_next_value)
				output = append(output, temp...)
				if err != nil {
					return output, err
				}
			} else {
				return output, errors.New(fmt.Sprintf("Error, state: %v", state))
			}
		} else {
			output = append(output, Result{match, state})
		}
	}
	return output, nil
}
