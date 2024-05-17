package util

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func PrintWithPrefix(prefix string, prefixColor string, value string) {
	prefixStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(prefixColor)).
		Bold(true)
	fmt.Printf("%s %s\n", prefixStyle.Render(prefix), value)
}
