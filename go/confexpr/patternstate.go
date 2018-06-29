package confexpr

import "fmt"

type PatternState struct {
	Path, Evaluated_path []string
	State                MatchState
	Value                interface{}
	Variables            map[string]string
}

func (this PatternState) String() string {
	return fmt.Sprintf("State: %v Path: %v Evaluated_Path:%v Value:%v Variables:%v", this.State, this.Path, this.Evaluated_path, this.Value, this.Variables)
}

func NewPatternState(path []string) PatternState {
	return PatternState{
		Path:           path,
		Evaluated_path: []string{},
		State:          Incomplete,
		Variables:      map[string]string{},
	}
}
