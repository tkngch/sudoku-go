package solver_test

import (
	"context"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tkngch/sudoku-go/internal/puzzle"
	"github.com/tkngch/sudoku-go/internal/solver"
)

// benchmarkingPuzzle is a hard 9x9 case that forces substantial backtracking.
type benchmarkingPuzzle struct {
	name     string
	input    string
	expected string
}

// benchmarkingPuzzles backs both TestSolveBenchmarkingPuzzles and
// BenchmarkSolve, so the puzzles the benchmark times are exactly the ones
// proven here to solve correctly. Each expected grid was verified independently
// of the solver: it agrees with the givens, and every row, column, and box is a
// permutation of 1-9.
func benchmarkingPuzzles() []benchmarkingPuzzle {
	return []benchmarkingPuzzle{
		{
			name: "diabolical - 1",
			input: "000100000" + "008090300" + "170800000" +
				"020000067" + "061050930" + "930000040" +
				"000002056" + "003040700" + "000001000",
			expected: "392165874" + "658794321" + "174823695" +
				"825439167" + "461257938" + "937618542" +
				"719382456" + "283546719" + "546971283",
		},
		{
			name: "diabolical - 2",
			input: "000607000" + "003080600" + "010309080" +
				"004000300" + "200000006" + "008050100" +
				"609108407" + "080000010" + "000705000",
			expected: "892647531" + "473581629" + "516329784" +
				"164872395" + "257913846" + "938456172" +
				"629138457" + "785264913" + "341795268",
		},
		{
			name: "diabolical - 3",
			input: "050090000" + "018400000" + "370008000" +
				"007002530" + "530000082" + "082500900" +
				"000900045" + "000003860" + "000080020",
			expected: "254397618" + "618425397" + "379618254" +
				"197842536" + "536179482" + "482536971" +
				"823961745" + "741253869" + "965784123",
		},
		{
			name: "diabolical - 4",
			input: "600305000" + "080200051" + "000000800" +
				"007010060" + "400703008" + "050040700" +
				"002000000" + "940007020" + "000109004",
			expected: "614385279" + "789264351" + "235971846" +
				"897512463" + "426793518" + "351846792" +
				"162458937" + "948637125" + "573129684",
		},
	}
}

func TestSolveBenchmarkingPuzzles(t *testing.T) {
	t.Parallel()

	for _, testCase := range benchmarkingPuzzles() {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				grid, err := puzzle.Parse(testCase.input)
				require.NoErrorf(t, err, "could not parse [%v]", testCase.input)

				expected, err := puzzle.Parse(testCase.expected)
				require.NoErrorf(t, err, "could not parse [%v]", testCase.expected)

				// Prove expected is a genuine solution consistent with the
				// givens, without consulting the solver.
				requireValidSolution(t, grid, expected)

				actual, err := solver.Solve(context.Background(), grid)
				require.NoError(t, err, "could not find a solution")

				assert.Equalf(
					t,
					expected,
					actual,
					"expected\n%s\nactual\n%s",
					expected.Render(),
					actual.Render(),
				)
			},
		)
	}
}

func BenchmarkSolve(b *testing.B) {
	for _, testCase := range benchmarkingPuzzles() {
		b.Run(
			testCase.name,
			func(b *testing.B) {
				grid, err := puzzle.Parse(testCase.input)
				require.NoError(b, err)

				expectedGrid, err := puzzle.Parse(testCase.expected)
				require.NoError(b, err)

				// Confirm the solver returns the correct grid once, outside the
				// timed loop. A regression that returns a fast but wrong answer
				// then fails the benchmark instead of masquerading as a
				// speed-up. TestSolveBenchmarkingPuzzles covers the same
				// puzzles under -race; this keeps the reported numbers honest.
				solution, err := solver.Solve(context.Background(), grid)
				require.NoError(b, err)
				require.Equal(b, expectedGrid, solution)

				for b.Loop() {
					_, _ = solver.Solve(context.Background(), grid)
				}
			},
		)
	}
}

// requireValidSolution verifies expected independently of solver.Solve, so a
// mistyped expected grid cannot slip through TestSolveBenchmarkingPuzzles by
// happening to match a solver bug. It confirms two properties: 1. solution
// keeps every given from input; and 2. solution is well-formed: every cell
// holds a single value and no row, column, or box repeats one.
func requireValidSolution(t *testing.T, input, solution *puzzle.Grid) {
	t.Helper()

	givens := slices.Collect(input.Cells())
	solved := slices.Collect(solution.Cells())
	require.Len(t, solved, len(givens))

	for idx, given := range givens {
		if given.Candidates().Count() != 1 {
			continue // input left this cell blank, so it is not a given
		}

		assert.Equalf(
			t,
			given.Candidates(),
			solved[idx].Candidates(),
			"expected disagrees with the given at %v",
			given.Position(),
		)
	}

	for cell := range solution.Cells() {
		assert.Equalf(
			t,
			1,
			cell.Candidates().Count(),
			"cell %v is not filled in the expected grid",
			cell.Position(),
		)

		for peer := range solution.AllPeersOf(cell.Position()) {
			assert.NotEqualf(
				t,
				cell.Candidates(),
				peer.Candidates(),
				"cells %v and %v share a value in the expected grid",
				cell.Position(),
				peer.Position(),
			)
		}
	}
}
