package main

import (
	"os"

	"adiab-flame/internal/cli"
)

// main only dispatches to the CLI package; the solver lives in internal/.
func main() {
	os.Exit(cli.RunCommand(os.Args[1:]))
}
