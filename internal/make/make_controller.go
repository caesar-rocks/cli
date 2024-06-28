package make

import (
	"fmt"
	"os"
	"strings"

	"github.com/caesar-rocks/cli/util"
)

type MakeControllerOpts struct {
	Input    string `description:"The name of the controller to create"`
	Resource bool   `description:"Whether or not this controller should be built with base resource"`
}

func MakeController(opts MakeControllerOpts) error {
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

	if opts.Resource {
		createResourceFile(packageName, controllerNameUpperCamel, controllerFilePath)
		controllerNameUpperCamel = controllerNameUpperCamel + "Resources"
	} else {
		createControllerFile(packageName, controllerNameUpperCamel, controllerFilePath)
	}
	registerController(packageName, controllerNameUpperCamel)

	util.PrintWithPrefix("success", "#00c900", "Controller created successfully.")

	return nil
}

func createControllerFile(packageName string, controllerNameUpperCamel string, controllerFilePath string) {
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
	util.PrintWithPrefix("created", "#6C757D", controllerFilePath)
}

func createResourceFile(packageName string, controllerNameUpperCamel string, controllerFilePath string) {
	// Form the contents for the controller file
	controllerTemplate := `package controllers

import (
	caesar "github.com/caesar-rocks/core"
	"citadel/app/models"
	"citadel/app/repositories"
	"net/http"
)


// Update config.Routes with the following:
// 1.
// Pass ApplicationsRepository *repositories.ApplicationsRepository, into RegisterRoutes()
// 
// 2. 
// Within RegisterRoutes() function:
// Applications Controller
// controller := controllers.NewApplicationsController(
// 	ApplicationsRepository,
// )
// router.Get("/routeplaceholder/{id}", func(ctx *caesar.Context) error {
// 	return controller.Index(ctx)
// })
// router.Get("/routeplaceholder/all", func(ctx *caesar.Context) error {
// 	return controller.Show(ctx)
// })
// router.Post("/routeplaceholder/create", func(ctx *caesar.Context) error {
// 	return controller.Create(ctx)
// })
// router.Delete("/routeplaceholder/{id}/delete", func(ctx *caesar.Context) error {
// 	return controller.Delete(ctx)
// })
// router.Patch("/routeplaceholder/update", func(ctx *caesar.Context) error {
// 	return controller.Update(ctx)
// })

type ApplicationsController struct {
	repo *repositories.ApplicationsRepository
}

func NewApplicationsController(repo *repositories.ApplicationsRepository) *ApplicationsController {
	return &ApplicationsController{repo: repo}
}

// Simple Response Serializers
type serializeApplicationsSingleResponse struct {
	Results *models.MyModel
}

type serializeApplicationsMultipleResponse struct {
	Results []models.MyModel
}

func (c *ApplicationsController) Index(ctx *caesar.Context) error {
	// curl localhost:3000/routeplaceholder/1/ -H "Content-Type: application/json"
	model, err := c.repo.FindOneBy(ctx.Context(), "id", ctx.PathValue("id"))
	if err != nil {
		return ctx.SendJSON(http.StatusInternalServerError)
	}
	serializedResponse := serializeSingleResponse{
		Results: model,
	}
	return ctx.SendJSON(serializedResponse)
}

func (c *ApplicationsController) Show(ctx *caesar.Context) error {
	// curl localhost:3000/routeplaceholders/all -H "Content-Type: application/json"
	models, err := c.repo.FindAll(ctx.Context())
	if err != nil {
		return caesar.NewError(http.StatusInternalServerError)
	}
	serializedResponse := serializeMultipleResponse{
		Results: models,
	}
	return ctx.SendJSON(serializedResponse)
}

func (c *ApplicationsController) Create(ctx *caesar.Context) error {
	// curl -X POST localhost:3000/routeplaceholder/create/ -H "Content-Type: application/json" -d '{"name": "myname"}'
	var data struct {
		Name string ` + "`json:\"name\"`" + `
	}
	if err := ctx.DecodeJSON(&data); err != nil {
		return caesar.NewError(http.StatusBadRequest)
	}
	model := &models.StorageBucket{
		Name: data.Name,
	}
	if err := c.repo.Create(ctx.Context(), model); err != nil {
		return caesar.NewError(http.StatusInternalServerError)
	}
	serializedResponse := serializeSingleResponse{
		Results: model,
	}
	return ctx.SendJSON(serializedResponse)
}

func (c *ApplicationsController) Delete(ctx *caesar.Context) error {
	// curl -X DELETE localhost:3000/routeplaceholder/1/delete/ -H "Content-Type: application/json"
	err := c.repo.DeleteOneWhere(ctx.Context(), "id", ctx.PathValue("id"))
	if err != nil {
		return caesar.NewError(http.StatusInternalServerError)
	}
	return ctx.SendJSON("Row Deleted")
}

func (c *ApplicationsController) Update(ctx *caesar.Context) error {
	// curl -X PATCH localhost:3000/routeplaceholder/update/ -H "Content-Type: application/json" -d '{"id": "1" , "name": "newname"}'
	var data struct {
		ID   string ` + "`json:\"id\"`" + `
		Name string ` + "`json:\"name\"`" + `
	}
	if err := ctx.DecodeJSON(&data); err != nil {
		return caesar.NewError(http.StatusBadRequest)
	}
	model := &models.StorageBucket{
		Name: data.Name,
	}
	if err := c.repo.UpdateOneWhere(ctx.Context(), "id", data.ID, model); err != nil {
		return caesar.NewError(http.StatusInternalServerError)
	}
	serializedResponse := serializeSingleResponse{
		Results: model,
	}
	return ctx.SendJSON(serializedResponse)
}

`
	// replace route path
	controllerTemplate = strings.ReplaceAll(controllerTemplate, "routeplaceholder", strings.ToLower(controllerNameUpperCamel))
	// replace app name
	controllerTemplate = strings.ReplaceAll(controllerTemplate, "citadel", util.RetrieveModuleName())
	// replace model name
	controllerTemplate = strings.ReplaceAll(controllerTemplate, "MyModel", controllerNameUpperCamel+"Resource")
	// replace controller + repository name
	controllerTemplate = strings.ReplaceAll(controllerTemplate, "Applications", controllerNameUpperCamel+"Resources")
	controllerTemplate = strings.ReplaceAll(controllerTemplate, "package controllers", "package "+packageName)

	// Write the contents to the file
	err := os.WriteFile(controllerFilePath, []byte(controllerTemplate), 0644)
	if err != nil {
		panic(err)
	}
	util.PrintWithPrefix("created", "#6C757D", controllerFilePath)
}

func registerController(packageName string, controllerNameUpperCamel string) {
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
