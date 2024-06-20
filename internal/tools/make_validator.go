package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/caesar-rocks/cli/util"
	"github.com/caesar-rocks/cli/util/inform"
)

func (wrapper *ToolsWrapper) MakeValidator(input string) error {
	validatorNameSnake := util.ConvertToSnakeCase(input)

	packageName := "validators"
	inputSnakeParts := strings.Split(validatorNameSnake, "/")

	var (
		subfolders        []string
		validatorFilePath string
	)

	if len(inputSnakeParts) > 1 {
		for i := 0; i < len(inputSnakeParts)-1; i++ {
			os.MkdirAll("app/validators/"+strings.Join(inputSnakeParts[:i+1], "/"), 0755)
		}

		validatorNameSnake = inputSnakeParts[len(inputSnakeParts)-1]

		subfolders = inputSnakeParts[:len(inputSnakeParts)-1]

		packageName = strings.Join(inputSnakeParts[:len(inputSnakeParts)-1], "_")

		validatorFilePath = fmt.Sprintf("./app/validators/%s/%s_validator.go", strings.Join(subfolders, "/"), validatorNameSnake)
	} else {
		validatorFilePath = fmt.Sprintf("./app/validators/%s_validator.go", validatorNameSnake)
	}

	validatorName := util.ConvertToUpperCamelCase(validatorNameSnake)
	if !strings.HasSuffix(validatorName, "Validator") {
		validatorName += "Validator"
	}

	validatorContents := fmt.Sprintf(`package %s

// For full reference, see: https://github.com/go-playground/validator
type %s struct {
	// Add your fields here.
}
`, packageName, validatorName)

	if err := os.WriteFile(validatorFilePath, []byte(validatorContents), 0644); err != nil {
		return err
	}

	wrapper.Inform(inform.Created, "Validator created successfully.")

	return nil
}
