package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeValidationCmd = &cobra.Command{
	Use:     "make:validator",
	Short:   "Create a new validator",
	GroupID: "make",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("make:validator called")
	},
}

func init() {
	rootCmd.AddCommand(makeValidationCmd)
}
