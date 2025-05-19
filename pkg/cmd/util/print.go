package util

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintTable(header table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetStyle(table.StyleDefault)

	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)

	for _, r := range rows {
		t.AppendRow(r)
	}

	t.AppendSeparator()
	t.Render()
}

func PrettyJSON(data []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "  ")
	if err != nil {
		return string(data)
	}
	return out.String()
}
