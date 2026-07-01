package sudoku_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/sudoku"
)

func TestRun(t *testing.T) {
	t.Parallel()

	const solved4x4 = "1234" + "3412" + "4321" + "2143"

	testCases := []struct {
		name           string
		args           []string
		stdin          string
		expectedCode   sudoku.ExitCode
		expectedStdout string
		stderrContains string
	}{
		{
			name:           "puzzle in an argument",
			args:           []string{".234" + "3.12" + "43.1" + "214."},
			expectedCode:   0,
			expectedStdout: solved4x4 + "\n",
			stderrContains: "Solution",
		},
		{
			name:           "multilined puzzle in stdin",
			args:           []string{},
			stdin:          ".234\n3.12\n43.1\n214.\n",
			expectedCode:   0,
			expectedStdout: solved4x4 + "\n",
			stderrContains: "Solution",
		},
		{
			name:           "empty stdin",
			args:           []string{},
			expectedCode:   1,
			stderrContains: "invalid cell count",
		},
		{
			name:           "too many arguments",
			args:           []string{solved4x4, solved4x4},
			expectedCode:   2,
			stderrContains: "usage:",
		},
		{
			name:           "invalid character",
			args:           []string{"z234123412341234"},
			expectedCode:   1,
			stderrContains: "invalid character",
		},
		{
			name:           "invalid cell count",
			args:           []string{"123"},
			expectedCode:   1,
			stderrContains: "invalid cell count",
		},
		{
			name:           "unsolvable puzzle",
			args:           []string{"11.." + "...." + "...." + "...."},
			expectedCode:   1,
			stderrContains: "solution not found",
		},
		{
			name:           "-h succeeds",
			args:           []string{"-h"},
			expectedCode:   0,
			stderrContains: "usage:",
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				var stdout, stderr bytes.Buffer

				code := sudoku.Run(
					testCase.args,
					strings.NewReader(testCase.stdin),
					&stdout,
					&stderr,
				)

				assert.Equal(t, testCase.expectedCode, code)
				assert.Equal(t, testCase.expectedStdout, stdout.String())
				assert.Contains(t, stderr.String(), testCase.stderrContains)
			},
		)
	}
}
