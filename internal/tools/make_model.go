package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caesar-rocks/cli/util"
	"github.com/caesar-rocks/cli/util/inform"
)

type MakeModelOpts struct {
	ModelName      string `description:"The name of the model to create"`
	WithRepository bool   `description:"Whether to create a repository for the model"`
}

func (wrapper *ToolsWrapper) MakeModel(opts MakeModelOpts) error {
	modelNameCamelCase := opts.ModelName
	modelNameSnakeCase := util.CamelToSnake(modelNameCamelCase)

	structName := strings.ToUpper(modelNameCamelCase[:1]) + modelNameCamelCase[1:]

	modelFileContents := fmt.Sprintf(`package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// Refer to Bun's documentation: https://bun.uptrace.dev/
type %s struct {
	ID        int64     `+"`bun:\"id,pk,autoincrement\"`"+`

	// Add your own fields here...

	CreatedAt time.Time `+"`bun:\"created_at,notnull,default:current_timestamp\"`"+`
	UpdatedAt time.Time `+"`bun:\"updated_at,notnull,default:current_timestamp\"`"+`
}

var _ bun.BeforeAppendModelHook = (*%s)(nil)

func (m *%s) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
`, structName, structName, structName)

	// Define the path for the model file
	modelsDir := "./app/models"
	modelFilePath := filepath.Join(modelsDir, fmt.Sprintf("%s.go", modelNameSnakeCase))
	// Check if the models directory exists, if not create it
	if _, err := os.Stat(modelsDir); os.IsNotExist(err) {
		err := os.MkdirAll(modelsDir, 0755)
		if err != nil {
			return fmt.Errorf("unable to create models directory: %v", err)
		}
	}

	// Create and write the model file
	err := os.WriteFile(modelFilePath, []byte(modelFileContents), 0644)
	if err != nil {
		return err
	}
	wrapper.Inform(inform.Created, modelFilePath)

	// Create the repository if requested
	if opts.WithRepository {
		return wrapper.MakeRepository(MakeRepositoryOpts{
			ModelName: modelNameCamelCase,
		})
	}

	return nil
}
