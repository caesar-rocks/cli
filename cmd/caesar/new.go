package main

import (
	"fmt"

	"github.com/caesar-rocks/cli/internal"
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
				if err := internal.SetupApp(appName, appNameSnake); err != nil {
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
		fmt.Println("To start the development server, run the following commands:")
		fmt.Println(prefixCommand(), styleCommand("cd "+util.ConvertToSnakeCase(appName)))
		fmt.Println(prefixCommand(), styleCommand("caesar dev"))
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
