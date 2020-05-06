package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/deanishe/awgo"
	"os/exec"
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

var (
	filterFlag = flag.String("filterFlag", "none", "get enabled widgets")
)

func init() {
	// Create a new Workflow using default settings.
	// Critical settings are provided by Alfred via environment variables,
	// so this *will* die in flames if not run in an Alfred-like environment.
	wf = aw.New()
}

func listWidgets() []Widget {
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

	return widgets

}

func filter(widgets []Widget, test func(Widget) bool) (ret []Widget) {
	for _, widget := range widgets {
		if test(widget) {
			ret = append(ret, widget)
		}
	}
	return
}

func run() {
	flag.Parse()

	widgets := listWidgets()
	filteredWidgets := []Widget{}

	var filterFunc func(widget Widget) bool

	switch *filterFlag {
	case "enabled":
		filterFunc = func(widget Widget) bool { return !widget.Hidden }
	case "disabled":
		filterFunc = func(widget Widget) bool { return widget.Hidden }
	case "none":
		filterFunc = func(widget Widget) bool { return true }
	}

	filteredWidgets = filter(widgets, filterFunc)

	for _, item := range filteredWidgets {
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
