package web

import (
	"context"
	"errors"
	"fmt"

	"github.com/tkngch/sudoku-go/internal/puzzle"
	"github.com/tkngch/sudoku-go/internal/solver"
)

// The error kinds that ErrorKind returns. The browser selects one message per
// kind, so these labels are part of the API and stay stable.
const (
	// KindSize reports an input whose cell count does not matches a supported
	// layout.
	KindSize = "size"

	// KindCharacter reports an input that holds a character which no cell
	// accepts.
	KindCharacter = "character"

	// KindUnsolvable reports a well-formed puzzle that has no solution.
	KindUnsolvable = "unsolvable"

	// KindUnknown reports any other error.
	KindUnknown = "unknown"
)

// Solve parses the puzzle, solves it, and returns the solution in the compact,
// one-line form that Grid.String produces. Solve accepts the same input as
// puzzle.Parse, so it ignores whitespace and it selects the layout from the
// number of cells.
//
// Solve applies no deadline. The js/wasm runtime has neither a sysmon thread
// nor asynchronous preemption, so a timer fires late during a solve. The caller
// owns the whole timeout policy: the browser stops a solve by terminating the
// worker.
//
// Pass the error to ErrorKind to obtain a short label for the browser.
func Solve(input string) (string, error) {
	grid, err := puzzle.Parse(input)
	if err != nil {
		return "", fmt.Errorf("web solve: %w", err)
	}

	solution, err := solver.Solve(context.Background(), grid)
	if err != nil {
		return "", fmt.Errorf("web solve: %w", err)
	}

	return solution.String(), nil
}

// ErrorKind maps an error from Solve to a short label. ErrorKind returns
// KindUnknown for a nil error, and for an error that it does not recognize.
func ErrorKind(err error) string {
	switch {
	case errors.Is(err, puzzle.ErrInvalidCellCount):
		return KindSize

	case errors.Is(err, puzzle.ErrInvalidCharacter):
		return KindCharacter

	case errors.Is(err, solver.ErrSolutionNotFound):
		return KindUnsolvable

	default:
		return KindUnknown
	}
}
