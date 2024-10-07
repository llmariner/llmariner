package ui

import (
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

// FormatTable formats the table.
func FormatTable(t table.Table) {
	if !enableColor() {
		// Do nothing.
		return
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	t.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
}
