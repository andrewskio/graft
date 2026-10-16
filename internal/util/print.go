package util

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

var (
	green = lipgloss.Color("#04B575")
	red   = lipgloss.Color("#D4634C")
)

func output(colorText string) {
	fmt.Println(colorText)
}

func PrintGreen(str ...string) {
	output(lipgloss.NewStyle().Foreground(green).Render(str...))
}
func PrintRed(str ...string) {
	output(lipgloss.NewStyle().Foreground(red).Render(str...))
}
