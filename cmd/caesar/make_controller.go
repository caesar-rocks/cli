package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/caesar-rocks/cli/util"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	input      string
	inputSnake string

	controllerNameUpperCamel string
	controllerNameSnake      string
	controllerFilePath       string
)

var makeControllerCmd = &cobra.Command{
	Use:     "make:controller",
	Short:   "Create a new controller",
	GroupID: "make",
	Run: func(cmd *cobra.Command, args []string) {
		// Get input.
		if len(args) > 0 {
			input = args[0]
		} else {
			huh.NewInput().Title("How should we name your controller, civis Romanus?").Value(&input).Run()
		}

		inputSnake = util.ConvertToSnakeCase(input)

		packageName := "controllers"
		inputSnakeParts := strings.Split(inputSnake, "/")

		var subfolders []string

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

		controllerNameUpperCamel = util.ConvertToUpperCamelCase(controllerNameSnake)

		createControllerFile(packageName)
		registerController(packageName)

		util.PrintWithPrefix("success", "#00c900", "Controller created successfully.")
	},
}

func createControllerFile(packageName string) {
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
// func (c *ApplicationsController) Index(ctx *caesar.CaesarCtx) error {
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
	util.PrintWithPrefix("created", "#6C757D", controllerFilePath)
}

func registerController(packageName string) {
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
		util.PrintWithPrefix("modified", "#6C757D", "./config/app.go")
	}
}

func init() {
	rootCmd.AddCommand(makeControllerCmd)
}
