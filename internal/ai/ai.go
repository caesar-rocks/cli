package ai

import (
	"context"
	"encoding/json"
	"fmt"
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
			Model:    openai.GPT4o,
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
		if msg.FinishReason == openai.FinishReasonToolCalls {
			for _, toolCall := range msg.Message.ToolCalls {
				tool := gen.tools[toolCall.Function.Name]
				args := make(map[string]any)
				json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
				fmt.Println(tool.Invoke(args))
				break
			}

		}
		break
	}

	return nil
}
