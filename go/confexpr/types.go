package confexpr

import "fmt"
import "errors"
import "sort"

type Result struct {
	Key   string
	State PatternState
}

func determine_match_states(start PatternState, rawConfig map[string]interface{}) ([]Result, error) {
	sorted_matches := []string{}
	for key, _ := range rawConfig {
		sorted_matches = append(sorted_matches, key)
	}

	sort.Slice(sorted_matches, comparatorLambda(sorted_matches))

	output := []Result{}
	prev_pattern := ""
	var prev_state *PatternState = nil
	for _, next_pattern := range sorted_matches {
		next_state, _ := next_state(next_pattern, start)
		prev_pattern, prev_state = next_with_collisions(next_pattern, prev_pattern, prev_state, next_state)
		output = append(output, Result{prev_pattern, *prev_state})
	}
	return output, nil
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
