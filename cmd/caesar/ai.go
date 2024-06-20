package main

import (
	"os"

	"github.com/caesar-rocks/cli/internal/ai"
	"github.com/caesar-rocks/cli/internal/ai/tool"
	"github.com/caesar-rocks/cli/internal/tools"
	"github.com/caesar-rocks/cli/util/inform"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var aiCmd = &cobra.Command{
	Use:     "ai",
	Short:   "Interact with your Caesar app making use of LLMs.",
	GroupID: "general",
	Run:     runAi,
}

func runAi(cmd *cobra.Command, args []string) {
	var prompt string

	if len(args) > 0 {
		prompt = args[0]
	} else {
		huh.NewInput().Title("What do you want to do, civis Romanus?").Value(&prompt).Run()
	}

	gen := ai.NewLlmGeneration()
	wrapper := tools.NewToolsWrapper(gen.StringsBuilder)
	gen.AddTool(tool.NewTool(wrapper.SetupApp, "Setup a new Caesar app"))
	gen.AddTool(tool.NewTool(wrapper.MakeController, "Generate a new Caesar controller"))
	gen.AddTool(tool.NewTool(wrapper.MakeMigration, "Generate a new Caesar migration"))
	gen.AddTool(tool.NewTool(wrapper.MakeModel, "Generate a new Bun model struct"))
	gen.AddTool(tool.NewTool(wrapper.MakeRepository, "Generate a new Caesar repository"))
	gen.AddTool(tool.NewTool(wrapper.AddUIComponent, "Add a new Caesar UI component"))
	if err := gen.Generate(prompt); err != nil {
		inform.Inform(os.Stdout, inform.Error, err.Error())
	}
}

func init() {
	rootCmd.AddCommand(aiCmd)
}
