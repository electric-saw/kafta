package util

import (
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
