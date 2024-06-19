package main

import (
	"github.com/caesar-rocks/cli/internal/make"
	"github.com/caesar-rocks/cli/util"
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

		if err := make.MakeRepository(make.MakeRepositoryOpts{
			ModelName: input,
		}); err != nil {
			util.PrintWithPrefix("error", "#FF0000", err.Error())
		}

		if err := make.MakeController(make.MakeControllerOpts{
			Input:    input,
			Resource: true,
		}); err != nil {
			util.PrintWithPrefix("error", "#FF0000", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(makeResourceCmd)
}
