package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeModelCmd = &cobra.Command{
	Use:     "make:model",
	Short:   "Create a new model",
	GroupID: "make",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("make:model called")
	},
}

func init() {
	rootCmd.AddCommand(makeModelCmd)
}
