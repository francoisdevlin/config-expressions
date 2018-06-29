package confexpr

import "fmt"
import "errors"
import "strings"
import "sort"

type Result struct {
	Key   string
	State PatternState
}

type LookupService struct {
	sortCache map[string][]string
}

func (this *LookupService) determine_match_states(start PatternState, rawConfig map[string]interface{}) ([]Result, error) {
	sorted_matches := this.get_sort_matches(start, rawConfig)

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

func (this *LookupService) get_sort_matches(start PatternState, rawConfig map[string]interface{}) []string {
	cacheState := strings.Join(start.evaluated_path, ".")
	sorted_matches, present := this.sortCache[cacheState]
	if present {
		return sorted_matches
	}
	sorted_matches = []string{}
	for key, _ := range rawConfig {
		sorted_matches = append(sorted_matches, key)
	}

	sort.Slice(sorted_matches, comparatorLambda(sorted_matches))

	this.sortCache[cacheState] = sorted_matches
	return sorted_matches
}

func NewLookupService() LookupService {
	sort_cache := map[string][]string{}
	return LookupService{
		sortCache: sort_cache,
	}
}

func (this *LookupService) Lookup(state PatternState, rawConfig map[string]interface{}) ([]Result, error) {
	results, err := this.determine_match_states(state, rawConfig)
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
				temp, err := this.Lookup(state, map_next_value)
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
