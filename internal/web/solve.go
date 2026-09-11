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
	// supported layout. It also reports an input that is longer than
	// MaxInputLength.
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

// maxInputLength is the longest input that Solve accepts. The unit is the
// byte.
//
// The largest layout holds 144 cells. Solve ignores whitespace, so it also
// accepts a grid form with a line break and a space between the cells. 4096
// keeps a large margin above that form.
//
// Solve needs this bound because puzzle.Parse allocates memory for each
// character before it counts the cells. A long input therefore costs much more
// memory than the input itself. In the browser this stops the WebAssembly
// runtime with an out-of-memory error.
const maxInputLength = 4096

// Solve parses the puzzle, solves it, and returns the solution in the compact,
// one-line form that Grid.String produces. Solve accepts the same input as
// puzzle.Parse, so it ignores whitespace and it selects the layout from the
// number of cells.
//
// Solve rejects an input that is longer than MaxInputLength, and it reports the
// kind ErrorKindSize. Solve applies this bound before it parses the input.
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
	if len(input) > maxInputLength {
		return "", ErrorKindSize
	}

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
