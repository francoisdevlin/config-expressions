package confexpr

import (
	"errors"
	"regexp"
	"strings"
)

type Label interface {
	both(path []string) ([]string, []string, error)
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
	for k, v := range state.Variables {
		output.Variables[k] = v
	}
	if this.variable != "" {
		output.Variables[this.variable] = strings.Join(consumed, ".")
	}
	output.Evaluated_path = append(state.Evaluated_path, consumed...)
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
	} else if _, ok := iLabel.(DeepWildcard); ok {
		output.WEIGHT = 4
	}

	return output
}

func (this DirectHit) both(path []string) ([]string, []string, error) {
	if path[0] == this.element {
		return path[1:], path[:1], nil
	} else {
		return []string{}, []string{}, errors.New("Path not found")
	}
}

type EnumHit struct {
	values map[string]bool
}

func NewEnumHit(values []string) EnumHit {
	valueCache := map[string]bool{}
	for _, value := range values {
		valueCache[value] = true
	}
	return EnumHit{
		values: valueCache,
	}
}

func (this EnumHit) both(path []string) ([]string, []string, error) {
	if _, found := this.values[path[0]]; found {
		return path[1:], path[:1], nil
	}
	return path, []string{}, errors.New("Path not found")
}

type RegexHit struct {
	regex string
}

func NewRegexHit(regex string) RegexHit {
	return RegexHit{
		regex: regex,
	}
}

func (this RegexHit) both(path []string) ([]string, []string, error) {
	r := regexp.MustCompile("^" + this.regex + "$")
	if r.MatchString(path[0]) {
		return path[1:], path[:1], nil
	}
	return path, []string{}, errors.New("Path not found")
}

type Wildcard struct {
}

func NewWildcard() Wildcard {
	return Wildcard{}
}

func (this Wildcard) both(path []string) ([]string, []string, error) {
	return path[1:], path[:1], nil
}

type DeepWildcard struct {
	label Label
}

func NewDeepWildcard(label Label) DeepWildcard {
	return DeepWildcard{
		label: label,
	}
}

func (this DeepWildcard) both(path []string) ([]string, []string, error) {
	maxIndex := -1
	for index, _ := range path {
		_, _, err := this.label.both(path[index:])
		if err == nil {
			maxIndex = index
		}
	}
	if maxIndex == -1 {
		return path, []string{}, errors.New("Sucessor path not found")
	}
	return path[maxIndex:], path[:maxIndex], nil
}
