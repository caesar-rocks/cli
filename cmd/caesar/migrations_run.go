package main

import (
	"github.com/caesar-rocks/cli/util"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "migrations:run",
	Short:   "Run all pending migrations",
	GroupID: "migrations",
	Run: func(cmd *cobra.Command, args []string) {
		util.Exec("go", "run", ".", "migrations:run")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
