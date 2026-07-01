package main

import (
	"os"

	"github.com/tkngch/sudoku-go/internal/sudoku"
)

func main() {
	exitCode := sudoku.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	os.Exit(int(exitCode))
}
