package confexpr

func FilterMissingResults(results []Result) []Result {
	output := []Result{}
	for _, result := range results {
		if result.State.State != Missing {
			output = append(output, result)
		}
	}
	return output
}

func FilterCollidingResults(results []Result) []Result {
	output := []Result{}
	for _, result := range results {
		if result.State.State == Collision {
			output = append(output, result)
		}
	}
	return output
}
