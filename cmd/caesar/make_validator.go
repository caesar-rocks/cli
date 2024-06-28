package main

import (
	"os"

	"github.com/caesar-rocks/cli/internal/tools"
	"github.com/caesar-rocks/cli/util"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var makeValidatorCmd = &cobra.Command{
	Use:     "make:validator",
	Short:   "Create a new validator",
	GroupID: "make",
	Run: func(cmd *cobra.Command, args []string) {
		input := ""

		if len(args) > 0 {
			input = args[0]
		} else {
			huh.NewInput().Title("How should we name your validator, civis Romanus?").Value(&input).Run()
		}

		wrapper := tools.NewToolsWrapper(os.Stdout)
		if err := wrapper.MakeValidator(tools.MakeValidatorOpts{
			ValidatiorName: input,
		}); err != nil {
			util.ExitWithError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(makeValidatorCmd)
}
