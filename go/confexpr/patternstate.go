package confexpr

import "fmt"

type PatternState struct {
	path, evaluated_path []string
	state                MatchState
	value                interface{}
}

func (this PatternState) String() string {
	return fmt.Sprintf("State: %v Path: %v Evaluated_Path:%v", this.state, this.path, this.evaluated_path)
}

func NewPatternState(path []string) PatternState {
	return PatternState{
		path:           path,
		evaluated_path: []string{},
		state:          Incomplete,
	}
}
