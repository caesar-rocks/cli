package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/caesar-rocks/cli/util"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

const (
	STARTER_KIT_GITHUB_REPO = "https://github.com/caesar-rocks/starter-kit.git"
)

var (
	appName      string
	appNameSnake string
)

var newCmd = &cobra.Command{
	Use:     "new [app-name]",
	Short:   "Initialize a new Caesar app",
	GroupID: "general",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			appName = args[0]
		} else {
			huh.NewInput().Title("How should we name your app, civis Romanus?").Value(&appName).Run()
		}

		appNameSnake = util.ConvertToSnakeCase(appName)

		err := spinner.New().
			Title("Setting up your new Caesar application").
			Action(setupApp).
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

func setupApp() {
	appNameSnake = util.ConvertToSnakeCase(appName)

	// Clone the starter kit repository
	_, err := git.PlainClone("./"+appNameSnake, false, &git.CloneOptions{
		URL:      STARTER_KIT_GITHUB_REPO,
		Progress: &util.Discard{},
		Depth:    1,
	})
	if err != nil {
		exitAndCleanUp(err)
	}

	// Remove the .git directory (we don't want to mess with the original repository)
	err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
	if err != nil {
		exitAndCleanUp(err)
	}

	// Write the .env file
	envExampleFile, err := os.ReadFile("./" + appNameSnake + "/.env.example")
	if err != nil {
		exitAndCleanUp(err)
	}

	envFileContents := string(envExampleFile)
	envFileContents = strings.ReplaceAll(envFileContents, "Caesar App", appName)
	envFileContents = strings.ReplaceAll(envFileContents, "<replace_by_app_key>", util.GenerateAppKey())

	err = os.WriteFile("./"+appNameSnake+"/.env", []byte(envFileContents), 0644)
	if err != nil {
		exitAndCleanUp(err)
	}

	renameGoModuleInGoModFile()
	renameGoModuleInGoAndTemplFiles("./" + appNameSnake)

	// Run `task install` to install dependencies
	cmd := exec.Command("task", "install")
	cmd.Dir = "./" + appNameSnake
	_, err = cmd.Output()
	if err != nil {
		exitAndCleanUp(err)
	}
}

func renameGoModuleInGoModFile() {
	goMod := fmt.Sprintf("./%s/go.mod", appNameSnake)
	goModFile, err := os.ReadFile(goMod)
	if err != nil {
		exitAndCleanUp(err)
	}

	goModContents := string(goModFile)
	goModContents = strings.ReplaceAll(goModContents, "starter_kit", appNameSnake)

	err = os.WriteFile(goMod, []byte(goModContents), 0644)
	if err != nil {
		exitAndCleanUp(err)
	}
}

func renameGoModuleInGoAndTemplFiles(basePath string) {
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		matched := strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".templ")

		if matched {
			read, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			newContents := strings.ReplaceAll(string(read), "starter_kit", appNameSnake)

			err = os.WriteFile(path, []byte(newContents), 0)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		exitAndCleanUp(err)
	}
}

func exitAndCleanUp(err error) {
	_ = os.RemoveAll(appNameSnake)
	util.ExitWithError(err)
}

func init() {
	rootCmd.AddCommand(newCmd)
}
