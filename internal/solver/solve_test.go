package solver_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tkngch/sudoku-go/internal/puzzle"
	"github.com/tkngch/sudoku-go/internal/solver"
)

func TestSolve(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		puzzle   []string
		solution []string
	}{
		{
			name: "4x4 without backtracking",
			puzzle: []string{
				".234",
				"3.12",
				"43.1",
				"214.",
			},
			solution: []string{
				"1234",
				"3412",
				"4321",
				"2143",
			},
		},
		{
			name: "9x9 without backtracking",
			puzzle: []string{
				"812753649",
				"940080170",
				"005491203",
				"054237090",
				"369845020",
				"207069504",
				"521970360",
				"438526917",
				"796318452",
			},
			solution: []string{
				"812753649",
				"943682175",
				"675491283",
				"154237896",
				"369845721",
				"287169534",
				"521974368",
				"438526917",
				"796318452",
			},
		},
		{
			name: "9x9 with backtracking",
			puzzle: []string{
				"800000000",
				"003600000",
				"070090200",
				"050007000",
				"000045700",
				"000100030",
				"001000068",
				"008500010",
				"090000400",
			},
			solution: []string{
				"812753649",
				"943682175",
				"675491283",
				"154237896",
				"369845721",
				"287169534",
				"521974368",
				"438526917",
				"796318452",
			},
		},
		{
			name: "4x4 with no solution",
			puzzle: []string{
				"11..",
				"....",
				"....",
				"....",
			},
			solution: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				grid, err := puzzle.Parse(strings.Join(testCase.puzzle, ""))
				require.NoErrorf(t, err, "could not parse [%v]", testCase.puzzle)

				actual, err := solver.Solve(grid)
				if testCase.solution == nil {
					require.ErrorIs(t, err, solver.ErrSolutionNotFound)

					return
				}

				require.NoError(t, err, "could not find a solution")

				expected, err := puzzle.Parse(strings.Join(testCase.solution, ""))
				require.NoErrorf(t, err, "could not parse [%v]", testCase.solution)

				assert.Equalf(t, expected.String(), actual.String(), "expected\n%s\nactual\n%s", expected.Render(), actual.Render())
			},
		)
	}
}
