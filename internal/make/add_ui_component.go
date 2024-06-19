package make

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/caesar-rocks/cli/util"
	"github.com/charmbracelet/huh"
)

const (
	CAESAR_UI_GITHUB_RAW_BASE_URL = "https://raw.githubusercontent.com/caesar-rocks/ui/master"
	CAESAR_UI_GITHUB_API_BASE_URL = "https://api.github.com/repos/caesar-rocks/ui/contents"

	DEFAULT_CSS_PATH                = "./views/app.css"
	DEFAULT_UI_COMPONENTS_BASE_PATH = "./views/ui"
)

type AddUIComponentOpts struct {
	ComponentName string `description:"The name of the component you want to add to your codebase (lowercase, matching the .templ file name)"`
}

func AddUIComponent(opts AddUIComponentOpts) error {
	componentContents, err := retrieveComponentContents(opts.ComponentName)
	if err != nil {
		return err
	}

	if err = saveComponentContents(opts.ComponentName, componentContents); err != nil {
		return err
	}

	componentStylesContents, err := retrieveComponentStyles(opts.ComponentName)
	if err == nil {
		if err := saveComponentStyles(componentStylesContents); err != nil {
			return err
		}
	}

	return nil
}

func retrieveComponentContents(componentName string) (string, error) {
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
		return "", fmt.Errorf("failed to retrieve component %s", componentName)
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

func ListRemoteUIComponents() ([]string, error) {
	resp, err := http.Get(CAESAR_UI_GITHUB_API_BASE_URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to retrieve UI components list")
	}

	var files []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	var components []string
	for _, file := range files {
		if file.Type == "file" && strings.HasSuffix(file.Name, ".templ") {
			componentName := strings.TrimSuffix(file.Name, ".templ")
			components = append(components, componentName)
		}
	}

	return components, nil
}
