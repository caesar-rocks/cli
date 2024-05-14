package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeControllerCmd = &cobra.Command{
	Use:     "make:controller",
	Short:   "Create a new controller",
	GroupID: "make",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("make:controller called")
	},
}

func init() {
	rootCmd.AddCommand(makeControllerCmd)
}
