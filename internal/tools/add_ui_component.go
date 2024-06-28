package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/caesar-rocks/cli/util/inform"
	"github.com/charmbracelet/huh"
)

const (
	CAESAR_UI_GITHUB_RAW_BASE_URL = "https://raw.githubusercontent.com/caesar-rocks/ui/master"

	DEFAULT_CSS_PATH                = "./views/app.css"
	DEFAULT_UI_COMPONENTS_BASE_PATH = "./views/ui"
)

type AddUIComponentOpts struct {
	ComponentName string `description:"The name of the component you want to add to your codebase (lowercase, matching the .templ file name)"`
}

func (wrapper *ToolsWrapper) AddUIComponent(opts AddUIComponentOpts) error {
	componentContents, err := wrapper.retrieveComponentContents(opts.ComponentName)
	if err != nil {
		return err
	}

	if err = wrapper.saveComponentContents(opts.ComponentName, componentContents); err != nil {
		return err
	}

	componentStylesContents, err := wrapper.retrieveComponentStyles(opts.ComponentName)
	if err == nil {
		if err := wrapper.saveComponentStyles(componentStylesContents); err != nil {
			return err
		}
	}

	return nil
}

func (wrapper *ToolsWrapper) retrieveComponentContents(componentName string) (string, error) {
	url := fmt.Sprintf("%s/%s.templ", CAESAR_UI_GITHUB_RAW_BASE_URL, componentName)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to retrieve component %s", componentName)
	}

	componentContents, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(componentContents), nil
}

func (wrapper *ToolsWrapper) saveComponentContents(componentName string, contents string) error {
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

	wrapper.Inform(inform.Info, fmt.Sprintf("created %s", componentFilePath))

	return nil
}

func (wrapper *ToolsWrapper) retrieveComponentStyles(componentName string) (string, error) {
	url := fmt.Sprintf("%s/%s.css", CAESAR_UI_GITHUB_RAW_BASE_URL, componentName)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to retrieve component %s", componentName)
	}

	componentStylesContents, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(componentStylesContents), nil
}

func (wrapper *ToolsWrapper) saveComponentStyles(contents string) error {
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

	wrapper.Inform(inform.Updated, cssPath)

	return nil
}
