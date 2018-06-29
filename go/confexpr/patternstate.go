package confexpr

import "fmt"

type PatternState struct {
	Path, Evaluated_path []string
	State                MatchState
	Value                interface{}
}

func (this PatternState) String() string {
	return fmt.Sprintf("State: %v Path: %v Evaluated_Path:%v Value:%v", this.State, this.Path, this.Evaluated_path, this.Value)
}

func NewPatternState(path []string) PatternState {
	return PatternState{
		Path:           path,
		Evaluated_path: []string{},
		State:          Incomplete,
	}
}
