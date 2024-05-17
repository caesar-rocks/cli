package main

import (
	"github.com/caesar-rocks/cli/util"

	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:     "migrations:reset",
	Short:   "Reset all migrations",
	GroupID: "migrations",
	Run: func(cmd *cobra.Command, args []string) {
		util.Exec("go", "run", ".", "migrations:reset")
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
