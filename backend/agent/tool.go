package agent

import "github.com/furisto/construct/backend/tool/codeact"

type ToolResult interface {
	kind() string
}

type InterpreterToolResult struct {
	ID            string                 `json:"id"`
	Output        string                 `json:"output"`
	FunctionCalls []codeact.FunctionCall `json:"function_calls"`
	Error         error                  `json:"error"`
}

func (r *InterpreterToolResult) kind() string {
	return "interpreter"
}

type NativeToolResult struct {
	ID     string `json:"id"`
	Output string `json:"output"`
	Error  error  `json:"error"`
}

func (r *NativeToolResult) kind() string {
	return "native"
}
