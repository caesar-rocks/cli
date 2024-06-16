package main

import (
	"github.com/caesar-rocks/cli/internal/ai"
	"github.com/caesar-rocks/cli/internal/ai/tool"
	"github.com/caesar-rocks/cli/internal/make"
	"github.com/caesar-rocks/cli/util"
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
	gen.AddTool(tool.NewTool(make.MakeController, "Generate a new controller"))
	if err := gen.Generate(prompt); err != nil {
		util.PrintWithPrefix("error", "#FF0000", err.Error())
	}
}

func init() {
	rootCmd.AddCommand(aiCmd)
}
