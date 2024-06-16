package make

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/caesar-rocks/cli/util"
	"github.com/go-git/go-git/v5"
)

const (
	STARTER_KIT_GITHUB_REPO = "https://github.com/caesar-rocks/starter-kit.git"
)

func SetupApp(appName string, appNameSnake string) error {
	_, err := git.PlainClone("./"+appNameSnake, false, &git.CloneOptions{
		URL:      STARTER_KIT_GITHUB_REPO,
		Progress: &util.Discard{},
		Depth:    1,
	})
	if err != nil {
		return err
	}

	if err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName)); err != nil {
		return err
	}

	envExampleFile, err := os.ReadFile("./" + appNameSnake + "/.env.example")
	if err != nil {
		return err
	}

	envFileContents := string(envExampleFile)
	envFileContents = strings.ReplaceAll(envFileContents, "Caesar App", appName)
	envFileContents = strings.ReplaceAll(envFileContents, "<replace_by_app_key>", util.GenerateAppKey())

	if err := os.WriteFile("./"+appNameSnake+"/.env", []byte(envFileContents), 0644); err != nil {
		return err
	}

	renameGoModuleInGoModFile(appNameSnake)
	renameGoModuleInGoAndTemplFiles(appNameSnake)

	cmd := exec.Command("task", "install")
	cmd.Dir = "./" + appNameSnake
	if _, err := cmd.Output(); err != nil {
		return err
	}

	cleanKeepFiles(appNameSnake)

	return nil
}

func renameGoModuleInGoModFile(appNameSnake string) {
	goMod := fmt.Sprintf("./%s/go.mod", appNameSnake)
	goModFile, err := os.ReadFile(goMod)
	if err != nil {
		util.ExitAndCleanUp(appNameSnake, err)
	}

	goModContents := string(goModFile)
	goModContents = strings.ReplaceAll(goModContents, "starter_kit", appNameSnake)

	if err = os.WriteFile(goMod, []byte(goModContents), 0644); err != nil {
		util.ExitAndCleanUp(appNameSnake, err)
	}
}

func renameGoModuleInGoAndTemplFiles(appNameSnake string) error {
	return filepath.Walk("./"+appNameSnake, func(path string, info os.FileInfo, err error) error {
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
}

// cleanKeepFiles deletes all .keep files in the app directory
func cleanKeepFiles(appNameSnake string) {
	filepath.Walk("./"+appNameSnake, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".keep") {
			if err := os.Remove(path); err != nil {
				return err
			}
		}

		return nil
	})
}
