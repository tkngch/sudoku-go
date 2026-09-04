package web_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tkngch/sudoku-go/internal/web"
)

func TestSolve(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "4x4",
			input:    ".234" + "3.12" + "43.1" + "214.",
			expected: "1234" + "3412" + "4321" + "2143",
		},
		{
			name:     "multiline input",
			input:    ".234\n3.12\n43.1\n214.\n",
			expected: "1234" + "3412" + "4321" + "2143",
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				solution, errKind := web.Solve(testCase.input)
				require.Empty(t, errKind)
				assert.Equal(t, testCase.expected, solution)
			},
		)
	}
}

func TestSolveErrorKind(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		input        string
		expectedKind web.ErrorKind
	}{
		{
			name:         "empty",
			input:        "",
			expectedKind: web.ErrorKindSize,
		},
		{
			name:         "cell count matches no layout",
			input:        "123",
			expectedKind: web.ErrorKindSize,
		},
		{
			name:         "unexpected character",
			input:        "z234" + "3.12" + "43.1" + "214.",
			expectedKind: web.ErrorKindCharacter,
		},
		{
			// A paste carries a byte order mark before a complete grid. The
			// cell count is correct, so the character is the fault.
			name:         "byte order mark before a full grid",
			input:        "\uFEFF" + ".234" + "3.12" + "43.1" + "214.",
			expectedKind: web.ErrorKindCharacter,
		},
		{
			name:         "value too large for the layout",
			input:        "9234" + "3.12" + "43.1" + "214.",
			expectedKind: web.ErrorKindValue,
		},
		{
			// The first row holds the value 1 twice, so no solution exists.
			name:         "repeated given",
			input:        "11.." + "...." + "...." + "....",
			expectedKind: web.ErrorKindUnsolvable,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				solution, errKind := web.Solve(testCase.input)
				assert.Empty(t, solution)
				assert.Equal(t, testCase.expectedKind, errKind)
			},
		)
	}
}
