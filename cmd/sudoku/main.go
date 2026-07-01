package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tkngch/sudoku-go/internal/sudoku"
)

func main() {
	exitCode := run()
	os.Exit(int(exitCode))
}

// run executes the program and returns an exit code. This short function is not
// inlined to main, to ensure context is cancelled before termination: `defer`
// fires on function returns but not os.Exit.
func run() sudoku.ExitCode {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return sudoku.Run(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
}
