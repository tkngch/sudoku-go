package web

import (
	"context"
	"errors"

	"github.com/tkngch/sudoku-go/internal/puzzle"
	"github.com/tkngch/sudoku-go/internal/solver"
)

// ErrorKind is a short label for an error from Solve. The browser selects one
// message per kind, so a label is part of the API and stays stable.
type ErrorKind string

const (
	// ErrorKindSize reports an input whose cell count does not match a
	// supported layout.
	ErrorKindSize ErrorKind = "size"

	// ErrorKindCharacter reports an input that holds a character which no cell
	// accepts.
	ErrorKindCharacter ErrorKind = "character"

	// ErrorKindValue reports a valid digit or letter that is too large for the
	// layout, such as '9' in a 4x4 puzzle.
	ErrorKindValue ErrorKind = "value"

	// ErrorKindUnsolvable reports a well-formed puzzle that has no solution.
	ErrorKindUnsolvable ErrorKind = "unsolvable"

	// ErrorKindUnknown reports any other error.
	ErrorKindUnknown ErrorKind = "unknown"
)

// Solve parses the puzzle, solves it, and returns the solution in the compact,
// one-line form that Grid.String produces. Solve accepts the same input as
// puzzle.Parse, so it ignores whitespace and it selects the layout from the
// number of cells.
//
// Solve applies no timeout. The js/wasm runtime omits the sysmon thread and
// asynchronous preemption, so the runtime delays a timer for an unbounded time
// during a solve. The caller owns the whole timeout policy. To stop a solve,
// the browser terminates the worker. Therefore, run Solve in a worker and not
// on the main thread.
//
// When a puzzle admits more than one solution, Solve returns one of them and
// does not detect or report non-uniqueness.
//
// Solve returns an empty ErrorKind after a success. Otherwise Solve returns an
// empty solution and an ErrorKind that labels the fault.
func Solve(input string) (string, ErrorKind) {
	grid, err := puzzle.Parse(input)
	if err != nil {
		return "", newErrorKind(err)
	}

	solution, err := solver.Solve(context.Background(), grid)
	if err != nil {
		return "", newErrorKind(err)
	}

	return solution.String(), ""
}

// newErrorKind maps an error from Solve to a short label. newErrorKind returns
// KindUnknown for an error that it does not recognize.
func newErrorKind(err error) ErrorKind {
	switch {
	case errors.Is(err, puzzle.ErrInvalidCellCount):
		return ErrorKindSize

	case errors.Is(err, puzzle.ErrInvalidCharacter):
		return ErrorKindCharacter

	case errors.Is(err, puzzle.ErrValueOutOfRange):
		return ErrorKindValue

	case errors.Is(err, solver.ErrSolutionNotFound):
		return ErrorKindUnsolvable

	default:
		return ErrorKindUnknown
	}
}
