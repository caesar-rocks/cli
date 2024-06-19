package main

import (
	"github.com/caesar-rocks/cli/internal/make"
	"github.com/caesar-rocks/cli/util"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var makeModelCmd = &cobra.Command{
	Use:     "make:model",
	Short:   "Create a new model",
	GroupID: "make",
	Run: func(cmd *cobra.Command, args []string) {
		var input string

		if len(args) > 0 {
			input = args[0]
		} else {
			huh.NewInput().Title("How should we name your model, civis Romanus?").Value(&input).Run()
		}

		if err := make.MakeModel(make.MakeModelOpts{
			ModelName: input,
		}); err != nil {
			util.PrintWithPrefix("error", "#FF0000", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(makeModelCmd)
}
