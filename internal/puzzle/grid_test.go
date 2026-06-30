package puzzle_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestNewGrid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		cells         []puzzle.Candidates
		layout        puzzle.Layout
		expectedError error
	}{
		{
			name:          "not enough cells for the layout",
			cells:         []puzzle.Candidates{},
			layout:        Must(puzzle.NewLayoutFromCellCount(16)),
			expectedError: puzzle.ErrInvalidCells,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				_, err := puzzle.NewGrid(testCase.cells, testCase.layout)
				assert.ErrorIs(t, err, testCase.expectedError)
			},
		)
	}
}

func TestGridSet(t *testing.T) {
	t.Parallel()

	t.Run(
		"cell is updated in place",
		func(t *testing.T) {
			t.Parallel()

			grid := newGrid(
				slices.Repeat([][]int{{1, 2, 3, 4}}, 4),
				Must(puzzle.NewLayoutFromCellCount(16)),
			)

			newCandidate := puzzle.NewSingleCandidate(9)
			grid.Set(puzzle.NewPosition(1, 0), newCandidate)

			cells := slices.Collect(grid.Cells())
			assert.Equal(t, newCandidate, cells[4].Candidates(), "cell is updated in place")
			// Every other cell is untouched.
			assert.Equal(t, puzzle.NewSingleCandidate(1), cells[0].Candidates())
			assert.Equal(t, puzzle.NewSingleCandidate(2), cells[1].Candidates())
			assert.Equal(t, puzzle.NewSingleCandidate(3), cells[2].Candidates())
			assert.Equal(t, puzzle.NewSingleCandidate(4), cells[3].Candidates())

			for idx := 5; idx < 16; idx++ {
				assert.Equalf(
					t,
					puzzle.NewSingleCandidate(idx%4+1),
					cells[idx].Candidates(),
					"cell %d",
					idx,
				)
			}
		},
	)

	t.Run(
		"does not panic when an out-of-bounds position is provided",
		func(t *testing.T) {
			t.Parallel()

			grid := newGrid(
				slices.Repeat([][]int{{1, 2, 3, 4}}, 4),
				Must(puzzle.NewLayoutFromCellCount(16)),
			)

			assert.NotPanics(t, func() {
				grid.Set(puzzle.NewPosition(0, 2), puzzle.NewSingleCandidate(9))
			})
		},
	)
}

func TestGridClone(t *testing.T) {
	t.Parallel()

	rows := slices.Repeat([][]int{{1, 2, 3, 4}}, 4)
	layout := Must(puzzle.NewLayoutFromCellCount(16))

	t.Run("mutating the clone leaves the original unchanged", func(t *testing.T) {
		t.Parallel()

		original := newGrid(rows, layout)
		clone := original.Clone()

		clone.Set(puzzle.NewPosition(0, 0), puzzle.NewSingleCandidate(2))

		originalCells := slices.Collect(original.Cells())
		cloneCells := slices.Collect(clone.Cells())

		assert.Equal(t, puzzle.NewSingleCandidate(1), originalCells[0].Candidates())
		assert.Equal(t, puzzle.NewSingleCandidate(2), cloneCells[0].Candidates())
	})

	t.Run("mutating the original leaves the clone unchanged", func(t *testing.T) {
		t.Parallel()

		original := newGrid(rows, layout)
		clone := original.Clone()

		original.Set(puzzle.NewPosition(0, 0), puzzle.NewSingleCandidate(2))

		originalCells := slices.Collect(original.Cells())
		cloneCells := slices.Collect(clone.Cells())

		assert.Equal(t, puzzle.NewSingleCandidate(2), originalCells[0].Candidates())
		assert.Equal(t, puzzle.NewSingleCandidate(1), cloneCells[0].Candidates())
	})
}

func newGrid(rows [][]int, layout puzzle.Layout) *puzzle.Grid {
	cells := make([]puzzle.Candidates, 0, layout.GridSize()*layout.GridSize())

	for _, rowValues := range rows {
		for _, value := range rowValues {
			cells = append(cells, puzzle.NewSingleCandidate(value))
		}
	}

	return Must(puzzle.NewGrid(cells, layout))
}
