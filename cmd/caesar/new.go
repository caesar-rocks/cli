package main

import (
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:     "new",
	Short:   "Initialize a new Caesar project",
	GroupID: "general",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			projectName string
			useTailwind bool   = true
			useHTMX     bool   = true
			dbms        string = "sqlite"
		)
		huh.NewInput().Title("What is your project named?").Value(&projectName).Run()
		huh.NewConfirm().Title("Would you like to use TailwindCSS?").Value(&useTailwind).Run()
		huh.NewConfirm().Title("Would you like to use HTMX?").Value(&useHTMX).Run()
		huh.NewSelect[string]().Title("Which database management system would you like to use?").Value(&dbms).Options(
			huh.NewOption("SQLite", "sqlite"),
			huh.NewOption("PostgreSQL", "postgres"),
			huh.NewOption("MySQL", "mysql"),
			huh.NewOption("Skip", ""),
		).Run()

		// Do something with the input below.
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
