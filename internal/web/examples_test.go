package web_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tkngch/sudoku-go/internal/web"
)

// examplesPath locates the data file that the browser downloads. The path is
// relative to this package directory.
const examplesPath = "../../web/examples.json"

// cellCount is the number of cells in a 9x9 grid. The web interface supports
// the 9x9 grid only, so every example holds this many characters.
const cellCount = 81

// example is one entry in web/examples.json. The web interface reads the same
// file, and it shows Name on the control that loads Puzzle.
//
// Each puzzle has exactly one solution, which a separate search confirmed. The
// test below proves that a solution exists. It does not prove that the solution
// is unique.
type example struct {
	Name   string `json:"name"`
	Puzzle string `json:"puzzle"`
}

// TestExamples keeps web/examples.json true to the solver. This test aims to
// prevent a broken example from reaching the browser.
func TestExamples(t *testing.T) {
	t.Parallel()

	examples := readExamples(t)
	require.NotEmpty(t, examples)

	seen := make(map[string]bool, len(examples))
	for _, entry := range examples {
		require.NotEmpty(t, entry.Name, "an example holds an empty name")
		require.NotContains(t, seen, entry.Name, "two examples hold the same name")

		seen[entry.Name] = true
	}

	for _, entry := range examples {
		t.Run(
			entry.Name,
			func(t *testing.T) {
				t.Parallel()

				assert.Len(t, entry.Puzzle, cellCount)

				solution, errKind := web.Solve(entry.Puzzle)
				require.Empty(t, errKind)
				assert.Len(t, solution, cellCount)
			},
		)
	}
}

func readExamples(t *testing.T) []example {
	t.Helper()

	content, err := os.ReadFile(examplesPath)
	require.NoError(t, err)

	var examples []example
	require.NoError(t, json.Unmarshal(content, &examples))

	return examples
}
