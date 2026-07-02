package sudoku_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/sudoku"
)

// terminalWithNoStdIn looks like an interactive terminal: Stat reports a
// character device, so Run should exit without reading.
type terminalWithNoStdIn struct{}

func (terminalWithNoStdIn) Read([]byte) (int, error)   { return 0, io.EOF }
func (terminalWithNoStdIn) Stat() (os.FileInfo, error) { return charDevice{}, nil }

// charDevice is a minimal os.FileInfo whose mode carries the char-device bit.
type charDevice struct{}

func (charDevice) Name() string       { return "" }
func (charDevice) Size() int64        { return 0 }
func (charDevice) Mode() os.FileMode  { return os.ModeCharDevice }
func (charDevice) ModTime() time.Time { return time.Time{} }
func (charDevice) IsDir() bool        { return false }
func (charDevice) Sys() any           { return nil }

// blockingReader closes reading when Read is called, then blocks until its
// release channel is closed, and reports EOF. It models a stdin that never
// yields data.
type blockingReader struct {
	reading chan struct{}
	release chan struct{}
}

func (b blockingReader) Read([]byte) (int, error) {
	close(b.reading)
	<-b.release

	return 0, io.EOF
}

// failingReader fails the test if it is ever read.
type failingReader struct {
	t *testing.T
}

func (r failingReader) Read([]byte) (int, error) {
	r.t.Error("stdin was read using the failing reader")

	return 0, io.EOF
}

var errSimulated = errors.New("simulated error")

type testSpec struct {
	name                   string
	context                func() (context.Context, context.CancelFunc)
	args                   []string
	stdin                  io.Reader
	expectedCode           sudoku.ExitCode
	expectedStdout         string
	expectedStderrContains string
}

const (
	unsolved4x4 = ".234" + "3.12" + "43.1" + "214."
	solved4x4   = "1234" + "3412" + "4321" + "2143"
)

func TestRun(t *testing.T) {
	t.Parallel()

	testCases := []testSpec{
		{
			name:                   "puzzle in an argument",
			args:                   []string{unsolved4x4},
			expectedCode:           sudoku.ExitOK,
			expectedStdout:         solved4x4 + "\n",
			expectedStderrContains: "Solution",
		},
		{
			name:                   "timeout disabled",
			args:                   []string{"-timeout", "0", unsolved4x4},
			expectedCode:           sudoku.ExitOK,
			expectedStdout:         solved4x4 + "\n",
			expectedStderrContains: "Solution",
		},
		{
			name:                   "timeout argument",
			args:                   []string{"-timeout", "5s", unsolved4x4},
			expectedCode:           sudoku.ExitOK,
			expectedStdout:         solved4x4 + "\n",
			expectedStderrContains: "Solution",
		},
		{
			name:                   "multilined puzzle in stdin",
			args:                   []string{},
			stdin:                  strings.NewReader(unsolved4x4 + "\n"),
			expectedCode:           sudoku.ExitOK,
			expectedStdout:         solved4x4 + "\n",
			expectedStderrContains: "Solution",
		},
		{
			name:                   "-h succeeds",
			args:                   []string{"-h"},
			expectedCode:           sudoku.ExitOK,
			expectedStderrContains: "usage:",
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()
				testRun(t, testCase)
			},
		)
	}
}

func TestRunError(t *testing.T) {
	t.Parallel()

	testCases := []testSpec{
		{
			name:                   "empty stdin",
			args:                   []string{},
			stdin:                  strings.NewReader(""),
			expectedCode:           sudoku.ExitError,
			expectedStderrContains: "invalid cell count",
		},
		{
			name:                   "invalid character",
			args:                   []string{"z234123412341234"},
			expectedCode:           sudoku.ExitError,
			expectedStderrContains: "invalid character",
		},
		{
			name:                   "invalid cell count",
			args:                   []string{"123"},
			expectedCode:           sudoku.ExitError,
			expectedStderrContains: "invalid cell count",
		},
		{
			name:                   "unsolvable puzzle",
			args:                   []string{"11.." + "...." + "...." + "...."},
			expectedCode:           sudoku.ExitError,
			expectedStderrContains: "solution not found",
		},
		{
			name:                   "stdin error",
			args:                   []string{},
			stdin:                  iotest.ErrReader(errSimulated),
			expectedCode:           sudoku.ExitError,
			expectedStderrContains: errSimulated.Error(),
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()
				testRun(t, testCase)
			},
		)
	}
}

func TestRunMisuseError(t *testing.T) {
	t.Parallel()

	testCases := []testSpec{
		{
			name:                   "too many arguments",
			args:                   []string{solved4x4, solved4x4},
			expectedCode:           sudoku.ExitMisuse,
			expectedStderrContains: "usage:",
		},
		{
			name:                   "unknown flag",
			args:                   []string{"-x"},
			expectedCode:           sudoku.ExitMisuse,
			expectedStderrContains: "not defined",
		},
		{
			name:                   "no args no pipe",
			args:                   []string{},
			stdin:                  terminalWithNoStdIn{},
			expectedCode:           sudoku.ExitMisuse,
			expectedStderrContains: "usage:",
		},
		{
			name:                   "negative timeout",
			args:                   []string{"-timeout", "-5s", solved4x4},
			expectedCode:           sudoku.ExitMisuse,
			expectedStderrContains: "timeout",
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()
				testRun(t, testCase)
			},
		)
	}
}

func TestRunContext(t *testing.T) {
	t.Parallel()

	testCases := []testSpec{
		{
			name:                   "timed-out",
			context:                alreadyTimedoutContext,
			args:                   []string{solved4x4},
			expectedCode:           sudoku.ExitTimeout,
			expectedStderrContains: "timed out",
		},
		{
			name:                   "already canceled",
			context:                alreadyCanceledContext,
			args:                   []string{solved4x4},
			expectedCode:           sudoku.ExitInterrupted,
			expectedStderrContains: "interrupted",
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()
				testRun(t, testCase)
			},
		)
	}
}

func TestRunCanceledBeforeReadingStdin(t *testing.T) {
	t.Parallel()

	// With an already-done context, Run must report the interruption
	// deterministically without reading stdin at all: it must not enter the
	// select where a already-completed read could win and change the exit code.
	// failingReader fails the test if stdin is read.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var stdout, stderr bytes.Buffer

	code := sudoku.Run(ctx, []string{}, failingReader{t: t}, &stdout, &stderr)

	assert.Equal(t, sudoku.ExitInterrupted, code)
	assert.Contains(t, stderr.String(), "interrupted")
}

func TestRunInterruptedDuringRead(t *testing.T) {
	t.Parallel()

	// A read that never completes must still be abortable through the context, so a
	// cancellation (Ctrl-C or SIGTERM) arriving mid-read yields ExitInterrupted
	// rather than hanging.
	release := make(chan struct{})
	defer close(release)

	reading := make(chan struct{})
	stdin := blockingReader{reading: reading, release: release}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stdout, stderr bytes.Buffer

	code := make(chan sudoku.ExitCode, 1)

	go func() {
		code <- sudoku.Run(ctx, []string{}, stdin, &stdout, &stderr)
	}()

	<-reading // the read is in progress and the context was live when it started
	cancel()  // interrupt mid-read

	assert.Equal(t, sudoku.ExitInterrupted, <-code)
	assert.Contains(t, stderr.String(), "interrupted")
}

func testRun(
	t *testing.T,
	testCase testSpec,
) {
	t.Helper()

	var (
		stdout, stderr bytes.Buffer
		ctx            context.Context
	)

	if testCase.context != nil {
		var cancel context.CancelFunc

		ctx, cancel = testCase.context()
		defer cancel()
	} else {
		ctx = context.Background()
	}

	code := sudoku.Run(ctx, testCase.args, testCase.stdin, &stdout, &stderr)

	assert.Equal(t, testCase.expectedCode, code)
	assert.Equal(t, testCase.expectedStdout, stdout.String())
	assert.Contains(t, stderr.String(), testCase.expectedStderrContains)
}

func alreadyTimedoutContext() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Unix(0, 0))
}

func alreadyCanceledContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	return ctx, cancel
}
