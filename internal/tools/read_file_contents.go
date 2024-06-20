package tools

import (
	"fmt"
	"os"

	"github.com/caesar-rocks/cli/util/inform"
)

type ReadFileContentsOpts struct {
	FilePath string `description:"The path to the file whose contents should be read"`
}

// ReadFileContents reads the contents of the specified file and returns them as a string.
func (wrapper *ToolsWrapper) ReadFileContents(opts ReadFileContentsOpts) error {
	contents, err := os.ReadFile(opts.FilePath)
	if err != nil {
		return err
	}

	wrapper.Inform(inform.Info, fmt.Sprintf("Contents of %s: %s", opts.FilePath, string(contents)))

	return nil
}
