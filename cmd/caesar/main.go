package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "caesar",
}

func main() {
	rootCmd.AddGroup(&cobra.Group{
		ID:    "general",
		Title: "General Commands",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "list",
		Title: "Make Commands",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "make",
		Title: "Make Commands",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "migrations",
		Title: "Migration Commands",
	})
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
