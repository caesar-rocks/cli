package tools

import (
	"fmt"
	"os"

	"github.com/caesar-rocks/cli/util/inform"
)

type WriteContentsToFileOpts struct {
	FilePath    string `description:"The path to the file whose contents should be replaced"`
	NewContents string `description:"The new contents of the file"`
}

// WriteContentsToFile replaces the contents of the specified file with newContents.
func (wrapper *ToolsWrapper) WriteContentsToFile(opts WriteContentsToFileOpts) error {
	err := os.WriteFile(opts.FilePath, []byte(opts.NewContents), 0644)
	if err != nil {
		return err
	}

	wrapper.Inform(inform.Updated, fmt.Sprintf("Replaced the contents of the file: %s", opts.FilePath))
	return nil
}
