package main

import (
	"os"

	"github.com/caesar-rocks/cli/internal/tools"
	"github.com/caesar-rocks/cli/util/inform"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var makeRepositoryCmd = &cobra.Command{
	Use:     "make:repository",
	Short:   "Create a new controller",
	GroupID: "make",
	Run: func(cmd *cobra.Command, args []string) {
		var input string

		if len(args) > 0 {
			input = args[0]
		} else {
			huh.NewInput().Title("Which model should we create a repository for, civis Romanus?").Value(&input).Run()
		}

		wrapper := tools.NewToolsWrapper(os.Stdout)
		if err := wrapper.MakeRepository(tools.MakeRepositoryOpts{
			ModelName: input,
		}); err != nil {
			inform.Inform(os.Stdout, inform.Error, err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(makeRepositoryCmd)
}
