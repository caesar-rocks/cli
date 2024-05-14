package main

import (
	"github.com/caesar-rocks/cli/util"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:     "dev",
	Short:   "Start the development server",
	GroupID: "general",
	Run:     runDevServer,
}

func runDevServer(cmd *cobra.Command, args []string) {
	util.Exec("task", "air", "css", "templ")
}

func init() {
	rootCmd.AddCommand(devCmd)
}
