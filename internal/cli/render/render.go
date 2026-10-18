/*
Package render contains utility commands for rending Lipgloss formatted
strings to be printed, as well as a Writer for formatting streams.
*/
package render

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
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
