package main

import (
	"github.com/caesar-rocks/cli/util"
	"github.com/spf13/cobra"
)

var listRoutesCmd = &cobra.Command{
	Use:     "list:routes",
	Short:   "List all routes of your application",
	GroupID: "list",
	Run: func(cmd *cobra.Command, args []string) {
		util.Exec("go", "run", ".", "list:routes")
	},
}

func init() {
	rootCmd.AddCommand(listRoutesCmd)
}
