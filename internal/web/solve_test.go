package web_test

import (
	"context"
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

				solution, err := web.Solve(testCase.input)
				require.NoError(t, err)
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
		expectedKind string
	}{
		{
			name:         "empty",
			input:        "",
			expectedKind: web.KindSize,
		},
		{
			name:         "cell count matches no layout",
			input:        "123",
			expectedKind: web.KindSize,
		},
		{
			name:         "unexpected character",
			input:        "z234" + "3.12" + "43.1" + "214.",
			expectedKind: web.KindCharacter,
		},
		{
			name:         "value too large for the layout",
			input:        "9234" + "3.12" + "43.1" + "214.",
			expectedKind: web.KindCharacter,
		},
		{
			// The first row holds the value 1 twice, so no solution exists.
			name:         "repeated given",
			input:        "11.." + "...." + "...." + "....",
			expectedKind: web.KindUnsolvable,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				solution, err := web.Solve(testCase.input)
				require.Error(t, err)
				assert.Empty(t, solution)
				assert.Equal(t, testCase.expectedKind, web.ErrorKind(err))
			},
		)
	}
}

// TestErrorKindUnknown covers the errors that Solve does not produce.
func TestErrorKindUnknown(t *testing.T) {
	t.Parallel()

	assert.Equal(t, web.KindUnknown, web.ErrorKind(nil))
	assert.Equal(t, web.KindUnknown, web.ErrorKind(context.Canceled))
}
