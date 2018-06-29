package confexpr

import "strings"

func deepWildcardPenalty(labels []WrappedLabel) int {
	output := 0
	for _, label := range labels {
		if _, dw := label.Label.(DeepWildcard); dw {
			output += 10000
		}
	}
	return output
}

func compare_patterns(left, right string) (int, error) {
	left_processors, _ := parse_processors(left)
	right_processors, _ := parse_processors(right)

	left_score := len(left_processors) + deepWildcardPenalty(left_processors)
	right_score := len(right_processors) + deepWildcardPenalty(right_processors)

	if left_score-right_score != 0 {
		return left_score - right_score, nil
	}

	for index := 0; index < len(left_processors); index++ {
		left_proc := left_processors[index]
		right_proc := right_processors[index]
		if left_proc.WEIGHT != right_proc.WEIGHT {
			return left_proc.WEIGHT - right_proc.WEIGHT, nil
		}
	}
	for index := 0; index < len(left_processors); index++ {
		left_proc := left_processors[index]
		right_proc := right_processors[index]
		left_hit, ok := left_proc.Label.(DirectHit)
		right_hit, ok := right_proc.Label.(DirectHit)
		if ok {
			result := strings.Compare(left_hit.element, right_hit.element)
			if result != 0 {
				return result, nil
			}
		}
	}
	//TODO - WILDCARD HANDYCAP
	return 0, nil
}

func comparatorLambda(strings []string) func(int, int) bool {
	return func(i, j int) bool {
		comp_result, _ := compare_patterns(strings[i], strings[j])
		return comp_result <= 0
	}
}
