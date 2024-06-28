package main

import (
	"github.com/caesar-rocks/cli/internal/make"
	"github.com/caesar-rocks/cli/util"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var makeControllerCmd = &cobra.Command{
	Use:     "make:controller",
	Short:   "Create a new controller",
	GroupID: "make",
	Run: func(cmd *cobra.Command, args []string) {
		var input string

		if len(args) > 0 {
			input = args[0]
		} else {
			huh.NewInput().Title("How should we name your controller, civis Romanus?").Value(&input).Run()
		}

		if err := make.MakeController(make.MakeControllerOpts{
			Input:    input,
			Resource: false,
		}); err != nil {
			util.PrintWithPrefix("error", "#FF0000", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(makeControllerCmd)
}
