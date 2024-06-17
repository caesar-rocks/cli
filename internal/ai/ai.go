package ai

import (
	"context"
	"encoding/json"
	"os"

	"github.com/caesar-rocks/cli/internal/ai/tool"
	"github.com/sashabaranov/go-openai"
)

type LlmGeneration struct {
	dialogue []openai.ChatCompletionMessage
	tools    map[string]tool.Tool
}

func NewLlmGeneration() *LlmGeneration {
	return &LlmGeneration{
		tools: make(map[string]tool.Tool),
	}
}

func (gen *LlmGeneration) AddTool(tool tool.Tool) {
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
		if msg.FinishReason == openai.FinishReasonStop {
			break
		}

		gen.dialogue = append(gen.dialogue, msg.Message)

		if msg.FinishReason == openai.FinishReasonToolCalls {
			for _, toolCall := range msg.Message.ToolCalls {
				tool := gen.tools[toolCall.Function.Name]
				args := make(map[string]any)
				json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
				res, err := tool.Invoke(args)
				if err != nil {
					return err
				}

				var errString string
				err, ok := res.(error)
				if ok {
					errString = "Error: " + err.Error()
				} else {
					errString = "No error occurred."
				}

				gen.dialogue = append(gen.dialogue, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    errString,
					Name:       toolCall.Function.Name,
					ToolCallID: toolCall.ID,
				})
			}
		}
	}

	return nil
}
