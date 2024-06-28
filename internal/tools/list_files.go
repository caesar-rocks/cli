package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caesar-rocks/cli/util/inform"
)

const (
	EXCLUDED_PATHS = "node_modules|.git|.idea|vendor|.env|.DS_Store"
)

// ListFiles lists all files in the current directory recursively, excluding node_modules.
func (wrapper *ToolsWrapper) ListFiles() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	var fileNames []string
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, excludedPath := range strings.Split(EXCLUDED_PATHS, "|") {
			if strings.Contains(path, excludedPath) {
				return nil
			}
		}

		if !d.IsDir() {
			fileNames = append(fileNames, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	wrapper.Inform(inform.Info, fmt.Sprintf("Found the following files: %v", fileNames))

	return nil
}
