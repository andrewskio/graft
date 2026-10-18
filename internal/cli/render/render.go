/*
Package render contains utility commands for rending Lipgloss formatted
strings to be printed, as well as a Writer for formatting streams.
*/
package render

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

var (
	green  = lipgloss.Color("#04B575")
	red    = lipgloss.Color("#D4634C")
	yellow = lipgloss.Color("#FFFF00")
)

func renderText(color color.Color, text string) string {
	return lipgloss.NewStyle().
		Foreground(color).
		Render(text)
}

func Success(format string, a ...any) string {
	return renderText(green, fmt.Sprintf(format, a...))
}

func Error(format string, a ...any) string {
	return renderText(red, fmt.Sprintf(format, a...))
}

func Warning(format string, a ...any) string {
	return renderText(yellow, fmt.Sprintf(format, a...))
}

/*
Table creates a rendered Table from a list of headers (strings), a list of
items of any type, and a callback mapping those item's fields to a list of strings

Example:

```

	render.Table(
		[]string{"Name", "Type"},
		srvConfs,
		func(s server.Server) []string {
			return []string{s.ServerName, s.ServerType}
		},
	)

```
*/
func Table[T any](headers []string, items []T, mapRow func(T) []string) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(headers...)

	for _, item := range items {
		t.Row(mapRow(item)...)
	}

	return t.Render()
}
