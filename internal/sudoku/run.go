package sudoku

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/tkngch/sudoku-go/internal/puzzle"
	"github.com/tkngch/sudoku-go/internal/solver"
)

type ExitCode int

const (
	exitOK     ExitCode = 0
	exitError  ExitCode = 1
	exitMisuse ExitCode = 2
)

const appName = "sudoku"

// Run parses a single Sudoku puzzle (from the sole argument, or from stdin when
// no argument is given), solves it, and writes the results to the given
// streams: the compact one-line solution to stdout, and human-readable
// diagnostics (the rendered input and solution, or an error message) to stderr.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) ExitCode {
	flags := flag.NewFlagSet(appName, flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		_, _ = fmt.Fprintf(stderr, "usage: %s [puzzle]\n", appName)

		flags.PrintDefaults()
	}

	err := flags.Parse(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return exitOK
		}

		return exitMisuse
	}

	var input string

	switch flags.NArg() {
	case 0:
		data, err := io.ReadAll(stdin)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "error: %v\n", err)

			return exitError
		}

		input = sanitize(string(data))
	case 1:
		input = sanitize(flags.Arg(0))
	default:
		flags.Usage()

		return exitMisuse
	}

	sudoku, err := puzzle.Parse(input)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "error: %v\n", err)

		return exitError
	}

	_, _ = fmt.Fprintf(stderr, "Sudoku\n%s\n", sudoku.Render())

	solution, err := solver.Solve(sudoku)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "error: %v\n", err)

		return exitError
	}

	_, _ = fmt.Fprintf(stderr, "Solution\n%s\n", solution.Render())
	_, _ = fmt.Fprintf(stdout, "%s\n", solution.String())

	return exitOK
}

// sanitize strips all whitespace so a puzzle may be supplied across multiple
// lines (for example pasted as a grid) and still match the compact,
// one-character-per-cell form that puzzle.Parse expects.
func sanitize(raw string) string {
	return strings.Join(strings.Fields(raw), "")
}
