package prompt

import (
	_ "embed"
)

//go:embed architect.md
var Plan string

//go:embed coder.md
var Edit string
