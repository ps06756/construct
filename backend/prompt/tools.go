package prompt

import (
	_ "embed"
)

//go:embed tools.md
var toolInstructions string

func ToolInstructions() string {
	return toolInstructions
}
