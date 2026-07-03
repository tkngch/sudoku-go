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
	// The benchmarking puzzles are taken from the article: Singh, M K. (2025).
	// Top 9 Hardest Sudoku Puzzles Ever Published.
	// https://sudokutimes.com/hardest-sudoku-puzzles/ (retrieved on July 2026).
	// Despite the article title, the article lists 8 unique puzzles, not 9. The
	// 9th puzzle in the article is identical to the 8th one. "Platinum Blonde"
	// in the article does not have a solution, and "Platinum Blonde" below was
	// taken from https://groups.google.com/g/sci.math/c/HXZfuNub7Cs
	return []benchmarkingPuzzle{
		{
			name: "AI Escargot",
			input: "100007090" + "030020008" + "009600500" +
				"005300900" + "010080002" + "600004000" +
				"300000010" + "040000007" + "007000300",
			expected: "162857493" + "534129678" + "789643521" +
				"475312986" + "913586742" + "628794135" +
				"356478219" + "241935867" + "897261354",
		},
		{
			name: "Golden Nugget",
			input: "020000009" + "000700100" + "000003000" +
				"000507800" + "500000030" + "000034067" +
				"005608000" + "608000040" + "000005600",
			expected: "726451389" + "853796124" + "491283756" +
				"364527891" + "517869432" + "289134567" +
				"945678213" + "678312945" + "132945678",
		},
		{
			name: "AI Escargot II",
			input: "023000700" + "006000000" + "700020406" +
				"000000800" + "007800234" + "000034060" +
				"000070002" + "000000000" + "010000600",
			expected: "123468759" + "846759123" + "759123486" +
				"234697815" + "967815234" + "581234967" +
				"698371542" + "472586391" + "315942678",
		},
		{
			name: "World's Hardest Sudoku",
			input: "005300000" + "800000020" + "070010500" +
				"400005300" + "010070006" + "003200080" +
				"060500009" + "004000030" + "000009700",
			expected: "145327698" + "839654127" + "672918543" +
				"496185372" + "218473956" + "753296481" +
				"367542819" + "984761235" + "521839764",
		},
		{
			name: "Platinum Blonde",
			input: "000000012" + "000000003" + "002300400" +
				"001800005" + "060070800" + "000009000" +
				"008500000" + "900040500" + "470006000",
			expected: "839465712" + "146782953" + "752391486" +
				"391824675" + "564173829" + "287659341" +
				"628537194" + "913248567" + "475916238",
		},
		{
			name: "Inescapable Puzzle",
			input: "000406700" + "406000000" + "000020406" +
				"004067000" + "000090030" + "801004067" +
				"305608010" + "000000000" + "010005008",
			expected: "928416753" + "476583129" + "153729486" +
				"534267891" + "762891534" + "891354267" +
				"345678912" + "689142375" + "217935648",
		},
		{
			name: "Impossible Sudoku",
			input: "023000009" + "006009000" + "089000000" +
				"000007800" + "507800004" + "001034007" +
				"000000000" + "000000005" + "012000000",
			expected: "123475689" + "456389172" + "789126453" +
				"234597861" + "597861234" + "861234597" +
				"645713928" + "978642315" + "312958746",
		},
		{
			name: "Everest Puzzle",
			input: "800000000" + "003600000" + "070090200" +
				"050007000" + "000045700" + "000100030" +
				"001000068" + "008500010" + "090000400",
			expected: "812753649" + "943682175" + "675491283" +
				"154237896" + "369845721" + "287169534" +
				"521974368" + "438526917" + "796318452",
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

		peers := solution.AllPeersOf(cell.Position())
		for idx := range peers.Len() {
			assert.NotEqualf(
				t,
				cell.Candidates(),
				solution.CandidatesAt(peers.At(idx)),
				"cells %v and %v share a value in the expected grid",
				cell.Position(),
				peers.At(idx),
			)
		}
	}
}
