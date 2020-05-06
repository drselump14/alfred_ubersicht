package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/deanishe/awgo"
	"os/exec"
)

var (
	query      = flag.String("query", "Google", "Search term")
	maxResults = flag.Int64("max-results", 10, "Max YouTube results")
)

// Widget does
// widget struct exported from main.scpt
//
type Widget struct {
	ID     string
	Hidden bool
}

// Workflow is the main API
var wf *aw.Workflow

func init() {
	// Create a new Workflow using default settings.
	// Critical settings are provided by Alfred via environment variables,
	// so this *will* die in flames if not run in an Alfred-like environment.
	wf = aw.New()
}

func run() {
	out, err := exec.Command("/usr/bin/osascript", "main.scpt").Output()

	if err != nil {
		log.Fatal(err)
	}
	var widgets []Widget

	jsonString := string(out)
	err = json.Unmarshal([]byte(jsonString), &widgets)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range widgets {
		wf.NewItem(item.ID).
			Subtitle(fmt.Sprintf("hidden: %v", item.Hidden)).
			Autocomplete(item.ID).
			Arg(item.ID).
			Valid(true)
	}

	wf.SendFeedback()

}

func main() {
	wf.Run(run)
}
