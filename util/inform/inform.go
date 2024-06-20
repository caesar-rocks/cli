package inform

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
)

func InformWithPrefix(w io.Writer, prefix string, prefixColor string, value string) {
	prefixStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(prefixColor)).
		Bold(true)
	fmt.Printf("%s %s\n", prefixStyle.Render(prefix), value)
}

type Level string

const (
	Created Level = "created"
	Deleted Level = "deleted"
	Updated Level = "updated"
	Success Level = "success"
	Info    Level = "info"
	Error   Level = "error"
)

func (l Level) String() string {
	return string(l)
}

func Inform(w io.Writer, level Level, message string, noColor ...bool) {
	if len(noColor) > 0 && noColor[0] {
		fmt.Printf("level: %s, message: %s\n", level, message)
		return
	}

	switch level {
	case Success:
	case Created:
		InformWithPrefix(w, level.String(), "#00c900", message)
	case Info:
		InformWithPrefix(w, level.String(), "#00c9ff", message)
	case Updated:
		InformWithPrefix(w, level.String(), "#6C757D", message)
	case Error:
	case Deleted:
		InformWithPrefix(w, level.String(), "#ff0000", message)
	}
}
