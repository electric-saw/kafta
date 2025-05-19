package util

import (
	"encoding/json"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

// DiffJSONs calculates the diff between two JSON strings and prints it in a pretty format.
func DiffJSONs(json1, json2 string) {
	// Parse the first JSON string into a generic map
	var obj1 map[string]interface{}
	if err := json.Unmarshal([]byte(json1), &obj1); err != nil {
		log.Fatalf("Failed to parse JSON1: %v", err)
	}

	// Parse the second JSON string into a generic map
	var obj2 map[string]interface{}
	if err := json.Unmarshal([]byte(json2), &obj2); err != nil {
		log.Fatalf("Failed to parse JSON2: %v", err)
	}

	// Use gojsondiff to calculate the diff
	differ := gojsondiff.New()
	diff := differ.CompareObjects(obj1, obj2)

	// Format the diff using a formatter
	formatter := formatter.NewAsciiFormatter(obj1, formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       true,
	})
	diffString, err := formatter.Format(diff)
	if err != nil {
		log.Fatalf("Failed to format diff: %v", err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Schema Diff"})
	t.AppendRow(table.Row{diffString})
	t.Render()
}
