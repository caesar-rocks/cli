package main

import (
	"github.com/caesar-rocks/cli/util"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:     "migrations:migrate",
	Short:   "Run all pending migrations",
	GroupID: "migrations",
	Run: func(cmd *cobra.Command, args []string) {
		util.Exec("go", "run", ".", "--migrate")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
