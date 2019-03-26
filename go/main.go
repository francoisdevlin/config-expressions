package main

import (
	"confexpr"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

const DEFAULT_FILE_PATH = "conf.jsonw"

func Lookup(lineArgs []string) {
	rawConfig := map[string]interface{}{}

	flagSet := flag.NewFlagSet("BACON", flag.ContinueOnError)
	confFilePath := ""
	flagSet.StringVar(&confFilePath, "file", DEFAULT_FILE_PATH, "Bacon")
	flagSet.Parse(lineArgs)
	restArgs := flagSet.Args()

	raw, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	json.Unmarshal(raw, &rawConfig)

	lookupService := confexpr.NewLookupService()
	entries := regexp.MustCompile("\\.").Split(restArgs[0], -1)
	start := confexpr.NewPatternState(entries)
	results, err := lookupService.Lookup(start, rawConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	results = confexpr.FilterMissingResults(results)

	if len(results) == 0 {
		fmt.Printf("Could not find %v\n", restArgs[0])
		os.Exit(1)
	}

	highestResult := results[0]

	if highestResult.State.State == confexpr.Complete {
		fmt.Println(highestResult.State.Value)
		os.Exit(0)
	}
	if highestResult.State.State == confexpr.Collision {
		repeats := confexpr.FilterCollidingResults(results)
		fmt.Printf("Ambiguous match for '%v', the following expressions have the same priority:\n", restArgs[0])
		for _, repeat := range repeats {
			fmt.Printf("'%v'", repeat.Key)
		}
		os.Exit(1)
	}

	fmt.Printf("Could not find %v\n", restArgs[0])
	os.Exit(1)

}

func PrintExplain(result confexpr.Result, conf map[string]interface{}) string {
	state := result.State
	pattern := result.Key
	if state.State == confexpr.Complete {
		return fmt.Sprintf("HIT    : %v VALUE: '%v'", pattern, state.Value)
	} else if state.State == confexpr.Collision {
		return fmt.Sprintf("COLLIDE: %v VALUE: '%v'", pattern, state.Value)
	} else {
		var value interface{} = conf
		for _, key := range state.Evaluated_path {
			mapValue, _ := value.(map[string]interface{})
			child, present := mapValue[key]
			if present {
				mapChild, okay := child.(map[string]interface{})
				if okay {
					value = mapChild
				}

			}
		}
		_, hasChildren := value.(map[string]interface{})
		message := ": " + pattern
		if hasChildren {
			message += ", ignoring children"
		}
		return "MISS   " + message
	}
	return ""
}

func Explain(lineArgs []string) {
	rawConfig := map[string]interface{}{}
	fmt.Println("The following rules were evaluated in this order, the first hit is returned")

	flagSet := flag.NewFlagSet("BACON", flag.ContinueOnError)
	confFilePath := ""
	noColor := false
	flagSet.StringVar(&confFilePath, "file", DEFAULT_FILE_PATH, "Bacon")
	flagSet.BoolVar(&noColor, "no-color", false, "Bacon")
	flagSet.Parse(lineArgs)
	restArgs := flagSet.Args()

	raw, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	json.Unmarshal(raw, &rawConfig)

	lookupService := confexpr.NewLookupService()
	entries := regexp.MustCompile("\\.").Split(restArgs[0], -1)
	start := confexpr.NewPatternState(entries)
	results, err := lookupService.Lookup(start, rawConfig)

	for index, result := range results {
		fmt.Printf("Rule %5d: %v\n", index+1, PrintExplain(result, rawConfig))
	}
	os.Exit(0)
}

func main() {
	restArgs := os.Args[1:]
	command := restArgs[0]
	restArgs = restArgs[1:]
	switch command {
	case "lookup":
		Lookup(restArgs)
	case "explain":
		Explain(restArgs)
	default:
		fmt.Println("Command not recognized")
		os.Exit(1)
	}
}
