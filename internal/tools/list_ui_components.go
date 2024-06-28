package tools

import (
	"fmt"

	"github.com/caesar-rocks/cli/internal/ui"
	"github.com/caesar-rocks/cli/util/inform"
)

func (wrapper *ToolsWrapper) ListUiComponents() error {
	components, err := ui.ListRemoteUIComponents()
	if err != nil {
		return err
	}

	wrapper.Inform(inform.Updated, fmt.Sprintf("Found the following UI components: %s", components))

	return nil
}
