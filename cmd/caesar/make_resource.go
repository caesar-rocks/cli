package main

import (
	"os"

	"github.com/caesar-rocks/cli/internal/tools"
	"github.com/caesar-rocks/cli/util/inform"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var makeResourceCmd = &cobra.Command{
	Use:     "make:resource",
	Short:   "Create a new resource",
	GroupID: "make",
	Run: func(cmd *cobra.Command, args []string) {
		var input string

		if len(args) > 0 {
			input = args[0]
		} else {
			huh.NewInput().Title("How should we name your resource, civis Romanus?").Value(&input).Run()
		}
		wrapper := tools.NewToolsWrapper(os.Stdout)

		if err := wrapper.MakeRepository(tools.MakeRepositoryOpts{
			ModelName: input + "Resources",
		}); err != nil {
			inform.Inform(os.Stdout, inform.Error, err.Error())
		}

		if err := wrapper.MakeController(tools.MakeControllerOpts{
			Input:    input,
			Resource: true,
		}); err != nil {
			inform.Inform(os.Stdout, inform.Error, err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(makeResourceCmd)
}
