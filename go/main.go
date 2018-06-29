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
	fmt.Println("Uh-oh")
	os.Exit(1)

}

func Explain(restArgs []string) {
	os.Exit(1)
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
