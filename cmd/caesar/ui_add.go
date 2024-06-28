package main

import (
	"os"

	"github.com/caesar-rocks/cli/internal/tools"
	"github.com/caesar-rocks/cli/internal/ui"
	"github.com/caesar-rocks/cli/util"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var uiAddCmd = &cobra.Command{
	Use:     "ui:add",
	Short:   "Add a CaesarUI component to your codebase",
	GroupID: "ui",
	Run: func(cmd *cobra.Command, args []string) {
		var componentName string

		if len(args) > 0 {
			componentName = args[0]
		} else {
			uiComponents, err := ui.ListRemoteUIComponents()
			if err == nil {
				options := make([]huh.Option[string], len(uiComponents))
				for i, c := range uiComponents {
					options[i] = huh.NewOption(util.SnakeToCamel(c), c)
				}
				huh.NewSelect[string]().
					Title("What component would you like to add?").
					Options(options...).
					Value(&componentName).
					Run()
			} else {
				huh.
					NewInput().
					Title("What is the name of the component you want to add (e.g. \"button\") ?").
					Value(&componentName).
					Run()
			}
		}

		wrapper := tools.NewToolsWrapper(os.Stdout)
		if err := wrapper.AddUIComponent(tools.AddUIComponentOpts{
			ComponentName: componentName,
		}); err != nil {
			util.ExitWithError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiAddCmd)
}
