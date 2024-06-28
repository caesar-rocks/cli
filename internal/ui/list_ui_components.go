package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	CAESAR_UI_GITHUB_API_BASE_URL = "https://api.github.com/repos/caesar-rocks/ui/contents"
)

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
