package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/caesar-rocks/cli/util"
	"github.com/caesar-rocks/cli/util/inform"
)

type MakeControllerOpts struct {
	Input string `description:"The name of the controller to create"`
}

func (wrapper *ToolsWrapper) MakeController(opts MakeControllerOpts) error {
	inputSnake := util.ConvertToSnakeCase(opts.Input)

	packageName := "controllers"
	inputSnakeParts := strings.Split(inputSnake, "/")

	var subfolders []string
	var controllerNameSnake string
	var controllerFilePath string

	if len(inputSnakeParts) > 1 {
		for i := 0; i < len(inputSnakeParts)-1; i++ {
			os.MkdirAll("app/controllers/"+strings.Join(inputSnakeParts[:i+1], "/"), 0755)
		}

		controllerNameSnake = inputSnakeParts[len(inputSnakeParts)-1]

		subfolders = inputSnakeParts[:len(inputSnakeParts)-1]

		packageName = strings.Join(inputSnakeParts[:len(inputSnakeParts)-1], "_")

		controllerFilePath = fmt.Sprintf("app/controllers/%s/%s_controller.go", strings.Join(subfolders, "/"), controllerNameSnake)
	} else {
		controllerNameSnake = inputSnake

		controllerFilePath = fmt.Sprintf("app/controllers/%s_controller.go", controllerNameSnake)
	}

	controllerNameUpperCamel := util.ConvertToUpperCamelCase(controllerNameSnake)

	wrapper.createControllerFile(packageName, controllerNameUpperCamel, controllerFilePath)
	wrapper.registerController(packageName, controllerNameUpperCamel)

	wrapper.Inform(inform.Created, "Controller created successfully.")

	return nil
}

func (wrapper *ToolsWrapper) createControllerFile(packageName string, controllerNameUpperCamel string, controllerFilePath string) {
	// Form the contents for the controller file
	controllerTemplate := `package controllers

// Uncomment the following import statement once you implement your first controller method.
// import (
// 	caesar "github.com/caesar-rocks/core"
// )

type ApplicationsController struct {
	// Define dependencies here
}

func NewApplicationsController() *ApplicationsController {
	return &ApplicationsController{
		// Inject dependencies here (like repositories, services, etc.)
	}
}

// Define controller methods here
// func (c *ApplicationsController) Index(ctx *caesar.Context) error {
// 	// Implement the controller method here

// 	return nil
// }
`
	controllerTemplate = strings.ReplaceAll(controllerTemplate, "Applications", controllerNameUpperCamel)
	controllerTemplate = strings.ReplaceAll(controllerTemplate, "package controllers", "package "+packageName)

	// Write the contents to the file
	err := os.WriteFile(controllerFilePath, []byte(controllerTemplate), 0644)
	if err != nil {
		panic(err)
	}
	wrapper.Inform(inform.Created, controllerFilePath)
}

func (wrapper *ToolsWrapper) registerController(packageName string, controllerNameUpperCamel string) {
	// Read file's contents
	bytes, err := os.ReadFile("./config/app.go")
	if err != nil {
		panic(err)
	}
	contents := string(bytes)

	moduleName := util.RetrieveModuleName()

	// Format the package path
	var packagePath string
	if packageName == "controllers" {
		packagePath = moduleName + "/app/controllers"
	} else {
		packagePath = moduleName + "/app/controllers/" + packageName
	}

	if packageName != "controllers" {
		packageName += "Controllers"
	}

	// Add import statement
	if packageName == "controllers" {
		alreadyImported := strings.Contains(contents, "\""+packagePath+"\"")
		if !alreadyImported {
			contents = strings.ReplaceAll(contents, "import (", "import (\n\t\""+packagePath+"\"")
		}
	} else {
		alreadyImported := strings.Contains(contents, packageName+" \""+packagePath+"\"")
		if !alreadyImported {
			contents = strings.ReplaceAll(contents, "import (", "import (\n\t"+packageName+" \""+packagePath+"\"")
		}
	}

	// Add controller to app.RegisterProviders
	if strings.Contains(contents, "app.RegisterProviders(") {
		contents = strings.Replace(contents, "app.RegisterProviders(", "app.RegisterProviders(\n\t\t"+packageName+".New"+controllerNameUpperCamel+"Controller,", 1)

		err = os.WriteFile("./config/app.go", []byte(contents), 0644)
		if err != nil {
			panic(err)
		}
		wrapper.Inform(inform.Updated, "./config/app.go")
	}
}
