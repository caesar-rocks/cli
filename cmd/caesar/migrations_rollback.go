package main

import (
	"github.com/caesar-rocks/cli/util"

	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:     "migrations:rollback",
	Short:   "Rollback the last migration",
	GroupID: "migrations",
	Run: func(cmd *cobra.Command, args []string) {
		util.Exec("go", "run", ".", "--rollback")
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}
