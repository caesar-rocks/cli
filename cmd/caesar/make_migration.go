package main

import (
	"os"

	"github.com/caesar-rocks/cli/internal/tools"
	"github.com/caesar-rocks/cli/util/inform"
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

		wrapper := tools.NewToolsWrapper(os.Stdout)
		if err := wrapper.MakeMigration(tools.MakeMigrationOpts{
			MigrationName: migrationName,
		}); err != nil {
			inform.Inform(os.Stdout, inform.Error, err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(makeMigrationCmd)
}
