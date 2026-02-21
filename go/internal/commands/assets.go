package commands

import _ "embed"

//go:embed assets/COMMANDS.md
var commandsMd string

//go:embed assets/clean
var cleanScript string

//go:embed assets/CONTRACT.md
var contractMd string
