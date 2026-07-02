package sudoku

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/tkngch/sudoku-go/internal/puzzle"
	"github.com/tkngch/sudoku-go/internal/solver"
)

type ExitCode int

const (
	ExitOK          ExitCode = 0
	ExitError       ExitCode = 1
	ExitMisuse      ExitCode = 2
	ExitTimeout     ExitCode = 124
	ExitInterrupted ExitCode = 130
)

const appName = "sudoku"

// errUsage signals CLI misuse (wrong argument count, or no puzzle when stdin is
// a terminal). Run responds by printing the usage message and returning
// ExitMisuse.
var errUsage = errors.New("usage error")

// Run parses a single Sudoku puzzle (from the sole argument, or from stdin when
// no argument is given), solves it, and writes the results to the given
// streams: the compact one-line solution to stdout, and human-readable
// diagnostics (the rendered input and solution, or an error message) to stderr.
func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) ExitCode {
	flags := flag.NewFlagSet(appName, flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		_, _ = fmt.Fprintf(stderr, "usage: %s [-timeout duration] [puzzle]\n", appName)

		flags.PrintDefaults()
	}
	timeoutPtr := flags.Duration(
		"timeout",
		30*time.Second,
		"maximum time to spend on solving; 0 disables the timeout",
	)

	err := flags.Parse(args)
	switch {
	case errors.Is(err, flag.ErrHelp):
		return ExitOK
	case err != nil:
		return ExitMisuse
	}

	timeout := *timeoutPtr
	if timeout < 0 {
		_, _ = fmt.Fprintf(stderr, "invalid -timeout %s: must not be negative\n", timeout)

		return ExitMisuse
	}

	input, err := resolveInput(ctx, flags, stdin)
	switch {
	case errors.Is(err, errUsage):
		flags.Usage()

		return ExitMisuse
	case err != nil:
		return reportError(stderr, err, timeout)
	}

	grid, err := puzzle.Parse(input)
	if err != nil {
		return fail(stderr, err)
	}

	// Start the solve deadline here so it measures time spent solving, not time
	// spent reading input.
	if timeout > 0 {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	_, _ = fmt.Fprintf(stderr, "Sudoku\n%s\n", grid.Render())

	solution, err := solver.Solve(ctx, grid)
	if err != nil {
		return reportError(stderr, err, timeout)
	}

	_, _ = fmt.Fprintf(stderr, "Solution\n%s\n", solution.Render())
	_, _ = fmt.Fprintf(stdout, "%s\n", solution.String())

	return ExitOK
}

// resolveInput returns the sanitized puzzle input from the sole argument or,
// when no argument is given, from stdin. It returns errUsage when the arguments
// are misused (no puzzle while stdin is a terminal, or more than one argument),
// or the read error (including a ctx cancellation) when stdin cannot be read.
func resolveInput(
	ctx context.Context,
	flags *flag.FlagSet,
	stdin io.Reader,
) (string, error) {
	switch flags.NArg() {
	case 0:
		if isTerminal(stdin) {
			return "", fmt.Errorf("resolve-input: %w", errUsage)
		}

		data, err := readAll(ctx, stdin)
		if err != nil {
			return "", fmt.Errorf("resolve-input: %w", err)
		}

		return sanitize(string(data)), nil
	case 1:
		return sanitize(flags.Arg(0)), nil
	default:
		return "", fmt.Errorf("resolve-input: %w", errUsage)
	}
}

// readAll reads all of reader, but stops early and returns ctx.Err() if ctx is
// done first. The read itself is not cancellable, so on early return the
// reading goroutine stays blocked until reader yields. The result channel is
// buffered so that goroutine can always send and exit, even after readAll has
// already returned.
func readAll(ctx context.Context, reader io.Reader) ([]byte, error) {
	// Give an already-fired cancellation priority over a read that has already
	// been completed.
	err := ctx.Err()
	if err != nil {
		return nil, fmt.Errorf("read-all: %w", err)
	}

	type result struct {
		data []byte
		err  error
	}

	done := make(chan result, 1)

	go func() {
		data, err := io.ReadAll(reader)
		done <- result{data: data, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("read-all: %w", ctx.Err())
	case res := <-done:
		return res.data, res.err
	}
}

// reportError reports err on stderr and returns the matching exit code. timeout
// is used only to annotate the timed-out message.
func reportError(stderr io.Writer, err error, timeout time.Duration) ExitCode {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		_, _ = fmt.Fprintf(stderr, "timed out after %s\n", timeout)

		return ExitTimeout
	case errors.Is(err, context.Canceled):
		_, _ = fmt.Fprintln(stderr, "interrupted")

		return ExitInterrupted
	default:
		return fail(stderr, err)
	}
}

// isTerminal reports whether r is an interactive character device (a TTY).
func isTerminal(r io.Reader) bool {
	f, ok := r.(interface{ Stat() (os.FileInfo, error) })
	if !ok {
		return false
	}

	fi, err := f.Stat()

	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}

// sanitize strips all whitespace so a puzzle may be supplied across multiple
// lines (for example pasted as a grid) and still match the compact,
// one-character-per-cell form that puzzle.Parse expects.
func sanitize(raw string) string {
	return strings.Join(strings.Fields(raw), "")
}

// fail reports err on stderr and returns the error exit code.
func fail(stderr io.Writer, err error) ExitCode {
	_, _ = fmt.Fprintf(stderr, "error: %v\n", err)

	return ExitError
}
