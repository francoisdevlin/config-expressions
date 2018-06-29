package confexpr

import "errors"

type Label interface {
	both(path []string) ([]string, string, error)
}

type WrappedLabel struct {
	WEIGHT int
	Label
	label, variable string
}

func (this WrappedLabel) next(state PatternState) (PatternState, error) {
	var rest, consumed, err = this.Label.both(state.Path)
	if err != nil {
		return state, err
	}
	var output = NewPatternState(rest)
	output.Evaluated_path = append(state.Evaluated_path, consumed)
	return output, nil
}

type DirectHit struct {
	element string
}

func NewDirectHit(element string) DirectHit {
	return DirectHit{element}
}

func NewWrappedLabel(iLabel Label, label, variable string) WrappedLabel {
	output := WrappedLabel{
		Label:    iLabel,
		label:    label,
		variable: variable,
		WEIGHT:   0,
	}

	if _, ok := iLabel.(DirectHit); ok {
		output.WEIGHT = 0
	} else if _, ok := iLabel.(EnumHit); ok {
		output.WEIGHT = 1
	} else if _, ok := iLabel.(RegexHit); ok {
		output.WEIGHT = 2
	} else if _, ok := iLabel.(Wildcard); ok {
		output.WEIGHT = 3
	}

	return output
}

func (this DirectHit) both(path []string) ([]string, string, error) {
	if path[0] == this.element {
		return path[1:], path[0], nil
	} else {
		return []string{}, "", errors.New("Path not found")
	}
}

type EnumHit struct {
}

func NewEnumHit() EnumHit {
	return EnumHit{}
}

func (this EnumHit) both(path []string) ([]string, string, error) {
	return path[1:], path[0], nil
}

type RegexHit struct {
}

func NewRegexHit() RegexHit {
	return RegexHit{}
}

func (this RegexHit) both(path []string) ([]string, string, error) {
	return path[1:], path[0], nil
}

type Wildcard struct {
}

func NewWildcard() Wildcard {
	return Wildcard{}
}

func (this Wildcard) both(path []string) ([]string, string, error) {
	return path[1:], path[0], nil
}
