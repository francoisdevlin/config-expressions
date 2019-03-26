package confexpr

import (
	"encoding/json"
	"testing"
)

func TestLookupServiceGetSortMatches(t *testing.T) {
	lookupService := NewLookupService()
	rawConfig := map[string]interface{}{
		"a": 1,
		"b": 2,
	}
	state := NewPatternState([]string{})
	state.Evaluated_path = []string{"1", "2", "3"}

	expect(len(lookupService.sortCache), 0, t)
	result := lookupService.get_sort_matches(state, rawConfig)
	expect(result, []string{"a", "b"}, t)
	expect(len(lookupService.sortCache), 1, t)

	result = lookupService.get_sort_matches(state, rawConfig)
	expect(result, []string{"a", "b"}, t)
	expect(len(lookupService.sortCache), 1, t)

	state.Evaluated_path = []string{"4", "5", "6"}
	result = lookupService.get_sort_matches(state, rawConfig)
	expect(result, []string{"a", "b"}, t)
	expect(len(lookupService.sortCache), 2, t)
}

func TestLookupServiceLookupHappyPath(t *testing.T) {
	data := `
{
	"development": {
		"app1.db.url" : "jdbc:h2:/sample/path",
		"app2.db.url" : "jdbc:h2:/sample/path",
		"app3.db.url" : "jdbc:h2:/sample/path"
	}
}
	`
	rawConfig := map[string]interface{}{}
	json.Unmarshal([]byte(data), &rawConfig)
	lookupService := NewLookupService()
	start := NewPatternState([]string{"development", "app1", "db", "url"})
	results, err := lookupService.Lookup(start, rawConfig)

	expect(results, []Result{
		Result{Key: "app1.db.url", State: PatternState{
			Value:          "jdbc:h2:/sample/path",
			Evaluated_path: []string{"development", "app1", "db", "url"},
			Path:           []string{},
			State:          Complete,
			Variables:      map[string]string{},
		}},
		Result{Key: "app2.db.url", State: PatternState{
			Value:          nil,
			Evaluated_path: []string{"development"},
			Path:           []string{"app1", "db", "url"},
			State:          Missing,
			Variables:      map[string]string{},
		}},
		Result{Key: "app3.db.url", State: PatternState{
			Value:          nil,
			Evaluated_path: []string{"development"},
			Path:           []string{"app1", "db", "url"},
			State:          Missing,
			Variables:      map[string]string{},
		}},
	}, t)
	expect(err, nil, t)
}

func TestLookupServiceLookupCollision(t *testing.T) {
	data := `
{
	"development": {
		"app1,app2.db.url" : "jdbc:h2:/first/path",
		"app2,app3.db.url" : "jdbc:h2:/other/path"
	}
}
	`
	rawConfig := map[string]interface{}{}
	json.Unmarshal([]byte(data), &rawConfig)
	lookupService := NewLookupService()
	start := NewPatternState([]string{"development", "app2", "db", "url"})
	results, err := lookupService.Lookup(start, rawConfig)

	expect(results, []Result{
		Result{Key: "app1,app2.db.url", State: PatternState{
			Value:          "jdbc:h2:/first/path",
			Evaluated_path: []string{"development", "app2", "db", "url"},
			Path:           []string{},
			State:          Collision,
			Variables:      map[string]string{},
		}},
		Result{Key: "app2,app3.db.url", State: PatternState{
			Value:          "jdbc:h2:/other/path",
			Evaluated_path: []string{"development", "app2", "db", "url"},
			Path:           []string{},
			State:          Collision,
			Variables:      map[string]string{},
		}},
	}, t)
	expect(err, nil, t)
}
