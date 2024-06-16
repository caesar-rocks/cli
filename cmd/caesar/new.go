package main

import (
	"fmt"

	"github.com/caesar-rocks/cli/internal/make"
	"github.com/caesar-rocks/cli/util"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:     "new [app-name]",
	Short:   "Initialize a new Caesar app",
	GroupID: "general",
	Run: func(cmd *cobra.Command, args []string) {
		appName := ""

		if len(args) > 0 {
			appName = args[0]
		} else {
			huh.NewInput().Title("How should we name your app, civis Romanus?").Value(&appName).Run()
		}

		appNameSnake := util.ConvertToSnakeCase(appName)

		err := spinner.New().
			Title("Setting up your new Caesar application").
			Action(func() {
				if err := make.SetupApp(appName, appNameSnake); err != nil {
					util.ExitAndCleanUp(appNameSnake, err)
				}
			}).
			Run()
		if err != nil {
			util.ExitWithError(err)
		}

		fmt.Println("Ave, Caesar!")
		fmt.Println("")
		fmt.Println("Your new Caesar app is ready to conquer the world.")
		fmt.Println()
		fmt.Println("To get started, change into your app's directory:")
		fmt.Println(prefixCommand(), styleCommand("cd "+util.ConvertToSnakeCase(appName)))
		fmt.Println()
		fmt.Println("To start your development setup, run the following commands (in separate terminals):")
		fmt.Println(prefixCommand(), styleCommand("task air"))
		fmt.Println()
		fmt.Println(prefixCommand(), styleCommand("task css"))
		fmt.Println()
		fmt.Println(prefixCommand(), styleCommand("task templ"))
		fmt.Println("")
		fmt.Println("For the glory of Rome! 🏛️")
	},
}

func prefixCommand() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6C757D")).
		Bold(true)

	return style.Render("❯")
}

func styleCommand(cmd string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#61AFEF"))

	return style.Render(cmd)
}

func init() {
	rootCmd.AddCommand(newCmd)
}
