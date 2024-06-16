package tool

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const (
	TOOL_CALLS_PREFIX_STR = "[TOOL_CALLS]"
)

type Tool struct {
	*openai.Tool
	Fn any
}

func NewTool(fn any, description string) Tool {
	return Tool{
		Tool: &openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        GetFunctionName(fn),
				Description: description,
				Parameters:  GetFunctionParameters(fn),
			},
		},
		Fn: fn,
	}
}

func (t *Tool) Invoke(args map[string]any) (any, error) {
	fnValue := reflect.ValueOf(t.Fn)
	if fnValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("not a function")
	}

	fnType := fnValue.Type()
	if fnType.NumIn() != 1 {
		return nil, fmt.Errorf("function must have exactly one parameter")
	}

	paramType := fnType.In(0)
	paramValue := reflect.New(paramType).Elem()

	for name := range t.Function.Parameters.(jsonschema.Definition).Properties {
		arg, ok := args[name]
		if !ok {
			return nil, fmt.Errorf("missing argument: %s", name)
		}

		field := paramValue.FieldByName(name)
		if !field.IsValid() {
			return nil, fmt.Errorf("no such field: %s in struct", name)
		}
		if !field.CanSet() {
			return nil, fmt.Errorf("cannot set field: %s", name)
		}

		argValue := reflect.ValueOf(arg)
		if argValue.Type() != field.Type() {
			return nil, fmt.Errorf("argument type mismatch for %s: expected %s, got %s", name, field.Type(), argValue.Type())
		}

		field.Set(argValue)
	}

	results := fnValue.Call([]reflect.Value{paramValue})
	if len(results) == 0 {
		return nil, nil
	}

	return results[0].Interface(), nil
}

type Parameter struct {
	Name        string
	Type        string
	Description string
}

func (p *Parameter) String() string {
	sb := strings.Builder{}

	sb.WriteString("\"")
	sb.WriteString(p.Name)
	sb.WriteString("\": {\"type\": \"")
	sb.WriteString(p.Type)
	if p.Description != "" {
		sb.WriteString("\", \"description\": \"")
		sb.WriteString(p.Description)
	}
	sb.WriteString("\"}")

	return sb.String()
}

func GetFunctionName(fn any) string {
	valueOfFn := reflect.ValueOf(fn)
	if valueOfFn.Kind() != reflect.Func {
		return ""
	}

	pc := valueOfFn.Pointer()

	function := runtime.FuncForPC(pc)
	if function == nil {
		return ""
	}

	name := function.Name()

	splitedName := strings.Split(name, ".")

	return splitedName[len(splitedName)-1]
}

func GetFunctionParameters(fn any) any {
	fnType := reflect.TypeOf(fn)

	optsStruct := fnType.In(0)

	required := []string{}
	properties := map[string]jsonschema.Definition{}

	for i := 0; i < optsStruct.NumField(); i++ {
		field := optsStruct.Field(i)
		properties[field.Name] = jsonschema.Definition{
			Type:        jsonschema.String,
			Description: field.Tag.Get("description"),
		}
		required = append(required, field.Name)
	}

	return jsonschema.Definition{
		Type:       jsonschema.Object,
		Properties: properties,
		Required:   required,
	}
}

type ToolCall struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

func IsToolCallsStrStarted(str string) bool {
	return strings.HasPrefix(str, TOOL_CALLS_PREFIX_STR) &&
		len(str) > len(TOOL_CALLS_PREFIX_STR)
}

func IsToolCallsStrFinished(str string) bool {
	return IsToolCallsStrStarted(str) &&
		strings.HasSuffix(str, "]")
}

func ParseToolCallsStr(str string) ([]ToolCall, error) {
	var toolCallsStr string

	re := regexp.MustCompile(regexp.QuoteMeta(TOOL_CALLS_PREFIX_STR) + `(.*)`)
	matches := re.FindStringSubmatch(str)
	if len(matches) == 2 {
		toolCallsStr = matches[1]
	}

	var toolCalls []ToolCall
	if err := json.Unmarshal([]byte(toolCallsStr), &toolCalls); err != nil {
		return nil, err
	}

	return toolCalls, nil
}
