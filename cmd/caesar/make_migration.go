package main

import (
	"github.com/caesar-rocks/cli/internal/make"
	"github.com/caesar-rocks/cli/util"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var makeMigrationCmd = &cobra.Command{
	Use:     "make:migration",
	Short:   "Create a new migration",
	GroupID: "make",
	Run: func(cmd *cobra.Command, args []string) {
		var migrationName string
		if len(args) > 0 {
			migrationName = args[0]
		} else {
			huh.NewInput().Title("How is your migration named?").Value(&migrationName).Run()
		}

		if err := make.MakeMigration(make.MakeMigrationOpts{
			MigrationName: migrationName,
		}); err != nil {
			util.PrintWithPrefix("error", "#FF0000", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(makeMigrationCmd)
}
