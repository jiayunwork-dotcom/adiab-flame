package main

import (
	"fmt"
	"os"

	"adiab-flame/internal/cli"
	"adiab-flame/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		runServe()
		return
	}
	if os.Args[1] == "serve" {
		runServe()
		return
	}
	os.Exit(cli.RunCommand(os.Args[1:]))
}

func runServe() {
	addr := ":8080"
	fmt.Fprintf(os.Stdout, "adiab-flame server listening on %s\n", addr)
	if err := server.ListenAndServe(server.Config{Addr: addr}); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
