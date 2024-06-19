package make

import (
	"fmt"
	"os"
	"strings"

	"github.com/caesar-rocks/cli/util"
	"github.com/gertd/go-pluralize"
)

const (
	REPOSITORIES_BASE_DIR = "./app/repositories"
)

type MakeRepositoryOpts struct {
	ModelName string `description:"The name of the model to create a repository for"`
}

func MakeRepository(opts MakeRepositoryOpts) error {
	modelNameSnake := util.CamelToSnake(opts.ModelName)

	packageName := "repositories"
	modelNameSnakeParts := strings.Split(modelNameSnake, "/")

	var subfolders []string
	var repositoryNameSnake string
	var repositoryFilePath string

	pluralize := pluralize.NewClient()

	// Create the repositories base directory if it doesn't exist
	if _, err := os.Stat(REPOSITORIES_BASE_DIR); os.IsNotExist(err) {
		err := os.MkdirAll(REPOSITORIES_BASE_DIR, 0755)
		if err != nil {
			return fmt.Errorf("unable to create repositories base directory: %v", err)
		}
	}

	if len(modelNameSnakeParts) > 1 {
		for i := 0; i < len(modelNameSnakeParts)-1; i++ {
			os.MkdirAll("app/repositories/"+strings.Join(modelNameSnakeParts[:i+1], "/"), 0755)
		}
		repositoryNameSnake = pluralize.Plural(modelNameSnakeParts[len(modelNameSnakeParts)-1])

		subfolders = modelNameSnakeParts[:len(modelNameSnakeParts)-1]

		packageName = strings.Join(modelNameSnakeParts[:len(modelNameSnakeParts)-1], "_")

		repositoryFilePath = fmt.Sprintf("app/repositories/%s/%s_repository.go", strings.Join(subfolders, "/"), repositoryNameSnake)
	} else {
		repositoryNameSnake = pluralize.Plural(modelNameSnake)

		repositoryFilePath = fmt.Sprintf("app/repositories/%s_repository.go", repositoryNameSnake)
	}

	repositoryNameUpperCamel := util.ConvertToUpperCamelCase(repositoryNameSnake)

	if err := createRepositoryFile(packageName, repositoryNameUpperCamel, repositoryFilePath); err != nil {
		return err
	}

	if err := registerRepository(packageName, repositoryNameUpperCamel); err != nil {
		return err
	}

	util.PrintWithPrefix("success", "#00c900", "Repository created successfully.")

	return nil
}

func createRepositoryFile(packageName string, repositoryNameUpperCamel string, repositoryFilePath string) error {
	// Form the contents for the repository file
	repositoryTemplate := `package repositories

import (
	"citadel/app/models"

	"github.com/caesar-rocks/orm"
)

type ApplicationsRepository struct {
	*orm.Repository[models.Application]
}

func NewApplicationsRepository(db *orm.Database) *ApplicationsRepository {
	return &ApplicationsRepository{Repository: &orm.Repository[models.Application]{
		Database: db,
	}}
}
`
	repositoryTemplate = strings.ReplaceAll(repositoryTemplate, "citadel", util.RetrieveModuleName())
	repositoryTemplate = strings.ReplaceAll(repositoryTemplate, "Applications", repositoryNameUpperCamel)
	repositoryTemplate = strings.ReplaceAll(repositoryTemplate, "Application", repositoryNameUpperCamel[:len(repositoryNameUpperCamel)-1])

	// Write the contents to the file
	err := os.WriteFile(repositoryFilePath, []byte(repositoryTemplate), 0644)
	if err != nil {
		return err
	}
	util.PrintWithPrefix("created", "#6C757D", repositoryFilePath)

	return nil
}

func registerRepository(packageName string, repositoryNameUpperCamel string) error {
	// Read the contents of the app.go file
	appFilePath := "./config/app.go"
	bytes, err := os.ReadFile(appFilePath)
	if err != nil {
		return fmt.Errorf("unable to read app.go file: %v", err)
	}
	contents := string(bytes)

	moduleName := util.RetrieveModuleName()
	repositoryStructName := repositoryNameUpperCamel + "Repository"

	// Format the package path
	var packagePath string
	if packageName == "repositories" {
		packagePath = moduleName + "/app/repositories"
	} else {
		packagePath = moduleName + "/app/repositories/" + packageName
	}

	if packageName != "repositories" {
		packageName += "Repositories"
	}

	// Add import statement if not already imported
	if packageName == "repositories" {
		if !strings.Contains(contents, "\""+packagePath+"\"") {
			contents = strings.ReplaceAll(contents, "import (", "import (\n\t\""+packagePath+"\"")
		}
	} else {
		if !strings.Contains(contents, packageName+" \""+packagePath+"\"") {
			contents = strings.ReplaceAll(contents, "import (", "import (\n\t"+packageName+" \""+packagePath+"\"")
		}
	}

	// Add repository to app.RegisterProviders
	if strings.Contains(contents, "app.RegisterProviders(") {
		contents = strings.Replace(
			contents,
			"app.RegisterProviders(",
			"app.RegisterProviders(\n\t\t"+packageName+".New"+repositoryStructName+",",
			1,
		)

		// Write the modified contents back to the file
		err = os.WriteFile(appFilePath, []byte(contents), 0644)
		if err != nil {
			return fmt.Errorf("unable to write to app.go file: %v", err)
		}
		util.PrintWithPrefix("modified", "#6C757D", appFilePath)
	}

	return nil
}
