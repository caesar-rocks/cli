package tools

import (
	"io"

	"github.com/caesar-rocks/cli/util/inform"
)

type ToolsWrapper struct {
	w io.Writer
}

func NewToolsWrapper(w io.Writer) *ToolsWrapper {
	return &ToolsWrapper{w}
}

func (wrapper *ToolsWrapper) Inform(level inform.Level, message string, noColor ...bool) {
	inform.Inform(wrapper.w, level, message, noColor...)
}
