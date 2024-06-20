package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/caesar-rocks/cli/internal/ai/tool"
	"github.com/caesar-rocks/cli/util/inform"
	"github.com/sashabaranov/go-openai"
)

const (
	INITIAL_PROMPT = `

	You are Geppetto, a senior web developer, working with the Caesar Go web framework.

	The Caesar Go web framework is a MVC framework taking inspiration from Ruby on Rails and Django, Laravel, ...

	It is wrapping the FX DI lib, the Bun ORM, Templ (a language for writing HTML user interfaces in Go), AlpineJS, HTMX, and Tailwind CSS.

	Templ syntax looks like this:
		package main

		templ Hello(name string) {
			<div>Hello, { name }</div>
		}

		templ Greeting(person Person) {
			<div class="greeting">
				@Hello(person.Name)
			</div>
		}

	In a Caesar web app,
	- routes are defined in the ./config/routes.go file (feel free to use the ReadFileContents tool to read the contents of this file)
	- FX providers and invokers are defined in the ./config/app.go file ; 
	- controllers are defined in the ./app/controllers directory ;
	- models are defined in the ./app/models directory ;
	- repositories are defined in the ./app/repositories directory ;
	- migrations are defined in the ./database/migrations directory ;
	- Templ files are defined in the ./views directory (with Caesar UI components in the ./views/components directory) ;
	- Templ layouts are defined in the ./views/layouts directory ;
	- Templ pages are defined in the ./views/pages directory ;
	- static assets are defined in the ./public directory.

	Your client asks you some web development tasks. You can use the tools provided to generate code for your Caesar web app.

	To accomplish your tasks, don't hesistate to get your hands dirty, and make use of all the tools at your disposal.
`
)

type LlmGeneration struct {
	dialogue []openai.ChatCompletionMessage
	tools    map[string]tool.AiTool

	// stringsBuilder is used to capture the output of the AI tool calls.
	StringsBuilder *strings.Builder
}

func NewLlmGeneration() *LlmGeneration {
	return &LlmGeneration{
		tools:          make(map[string]tool.AiTool),
		StringsBuilder: &strings.Builder{},
	}
}

func (gen *LlmGeneration) AddTool(tool tool.AiTool) {
	gen.tools[tool.Function.Name] = tool
}

func (gen *LlmGeneration) Generate(prompt string) error {
	ctx := context.Background()
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	var tools []openai.Tool
	for _, tool := range gen.tools {
		tools = append(tools, *tool.Tool)
	}

	gen.dialogue = []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: INITIAL_PROMPT},
		{Role: openai.ChatMessageRoleUser, Content: prompt},
	}

	for {
		resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:    openai.GPT4TurboPreview,
			Tools:    tools,
			Messages: gen.dialogue,
		})
		if err != nil {
			return err
		}

		msg := resp.Choices[0]

		gen.dialogue = append(gen.dialogue, msg.Message)

		if msg.FinishReason == openai.FinishReasonToolCalls {
			for _, toolCall := range msg.Message.ToolCalls {
				tool := gen.tools[toolCall.Function.Name]
				args := make(map[string]any)
				json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

				inform.Inform(os.Stdout, inform.Info, "\"Calling tool: "+toolCall.Function.Name+"\"")
				res, err := tool.Invoke(args)
				if err != nil {
					inform.Inform(os.Stdout, inform.Error, err.Error())
					return err
				}

				var errString string
				err, ok := res.(error)
				if ok {
					errString = "Error: " + err.Error()
				} else {
					errString = "No error occurred."
				}
				inform.Inform(os.Stdout, inform.Info, "\"Tool call completed.\"")
				inform.Inform(os.Stdout, inform.Info, "\"Output: "+gen.StringsBuilder.String()+"\"")

				content := fmt.Sprintf("Error: %s\nOutput: %s", errString, gen.StringsBuilder.String())
				fmt.Println(gen.StringsBuilder.String())
				gen.dialogue = append(gen.dialogue, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    content,
					Name:       toolCall.Function.Name,
					ToolCallID: toolCall.ID,
				})
				gen.StringsBuilder.Reset()
			}

			continue
		}

		if msg.FinishReason == openai.FinishReasonStop {
			inform.Inform(os.Stdout, inform.Ai, msg.Message.Content)
			break
		}

		if msg.FinishReason == openai.FinishReasonLength {
			inform.Inform(os.Stdout, inform.Ai, "The response is too long.")
			break
		}

		if msg.FinishReason == openai.FinishReasonContentFilter {
			inform.Inform(os.Stdout, inform.Error, "The response was flagged by the content filter.")
			break
		}

		if msg.FinishReason == openai.FinishReasonNull {
			inform.Inform(os.Stdout, inform.Error, "The response was null.")
			break
		}

		continue
	}

	return nil
}
