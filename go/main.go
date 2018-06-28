package main

import "fmt"
import "confexpr"

func main() {
	rawConfig := map[string]interface{}{
		"a": 1,
		"b": map[string]interface{}{
			"c": 1,
			"d": 2,
		},
		"b.e": 5,
	}

	start := confexpr.NewPatternState([]string{"b", "c"})
	fmt.Println(confexpr.Missing)
	fmt.Println(confexpr.Lookup(start, rawConfig))
}
