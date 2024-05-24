package main

import (
	"fmt"

	"github.com/caesar-rocks/cli/internal"
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
			huh.NewInput().Title("What is the name of the component you want to add (e.g. \"button\") ?").Value(&componentName).Run()
		}

		if err := internal.AddUIComponent(componentName); err != nil {
			util.ExitWithError(err)
		}

		util.PrintWithPrefix("success", "#00c900", fmt.Sprintf("Component \"%s\" added successfully!", componentName))
	},
}

func init() {
	rootCmd.AddCommand(uiAddCmd)
}
