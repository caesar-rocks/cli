package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/caesar-rocks/cli/util"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

const (
	CAESAR_UI_GITHUB_RAW_BASE_URL = "https://raw.githubusercontent.com/caesar-rocks/ui/master"

	DEFAULT_CSS_PATH                = "./views/app.css"
	DEFAULT_UI_COMPONENTS_BASE_PATH = "./views/ui"
)

var uiAddCmd = &cobra.Command{
	Use:     "ui:add",
	Short:   "Add a CaesarUI component to your codebase",
	GroupID: "ui",
	Run: func(cmd *cobra.Command, args []string) {
		var componentName string

		if len(args) > 0 {
			componentName = args[0]
		} else {
			huh.NewInput().Title("What is the name of the component you want to add (e.g. \"button\") ?").Value(&componentName).Run()
		}

		componentContents, err := retrieveComponentContents(componentName)
		if err != nil {
			util.ExitWithError(err)
		}

		err = saveComponentContents(componentName, componentContents)
		if err != nil {
			util.ExitWithError(err)
		}

		componentStylesContents, err := retrieveComponentStyles(componentName)
		if err == nil {
			if err := saveComponentStyles(componentStylesContents); err != nil {
				util.ExitWithError(err)
			}
		}

		util.PrintWithPrefix("success", "#00c900", fmt.Sprintf("Component \"%s\" added successfully!", componentName))
	},
}

func retrieveComponentContents(componentName string) (string, error) {
	url := fmt.Sprintf("%s/%s.templ", CAESAR_UI_GITHUB_RAW_BASE_URL, componentName)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Failed to retrieve component %s", componentName)
	}

	componentContents, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(componentContents), nil
}

func saveComponentContents(componentName string, contents string) error {
	uiComponentsBasePath := DEFAULT_UI_COMPONENTS_BASE_PATH
	huh.NewInput().Title("In which folder would you save the component?").Value(&uiComponentsBasePath).Run()

	if _, err := os.Stat(uiComponentsBasePath); os.IsNotExist(err) {
		os.MkdirAll(uiComponentsBasePath, 0755)
	}

	componentFilePath := fmt.Sprintf("%s/%s.templ", uiComponentsBasePath, componentName)
	err := os.WriteFile(componentFilePath, []byte(contents), 0644)
	if err != nil {
		return err
	}

	util.PrintWithPrefix("created", "#6C757D", componentFilePath)

	return nil
}

func retrieveComponentStyles(componentName string) (string, error) {
	url := fmt.Sprintf("%s/%s.css", CAESAR_UI_GITHUB_RAW_BASE_URL, componentName)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Failed to retrieve component %s", componentName)
	}

	componentStylesContents, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(componentStylesContents), nil
}

func saveComponentStyles(contents string) error {
	cssPath := DEFAULT_CSS_PATH
	huh.NewInput().Title("In which CSS file would you save the component styles?").Value(&cssPath).Run()

	if _, err := os.Stat(cssPath); os.IsNotExist(err) {
		return os.WriteFile(cssPath, []byte(contents), 0644)
	}

	f, err := os.OpenFile(cssPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString("\n" + contents); err != nil {
		return err
	}

	util.PrintWithPrefix("updated", "#6C757D", cssPath)

	return nil
}

func init() {
	rootCmd.AddCommand(uiAddCmd)
}
