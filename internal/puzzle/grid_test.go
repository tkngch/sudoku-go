package puzzle_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestGridSet(t *testing.T) {
	t.Run(
		"cell is updated in place",
		func(t2 *testing.T) {
			grid := newGrid(
				slices.Repeat([][]uint8{{1, 2, 3, 4}}, 4),
				Must(puzzle.NewLayoutFromCellCount(16)),
			)

			newCandidate := puzzle.NewCandidate(9)
			grid.Set(puzzle.NewPosition(1, 0), newCandidate)

			cells := slices.Collect(grid.Cells())
			assert.Equal(t2, newCandidate, cells[4].Candidates(), "cell is updated in place")
			// Every other cell is untouched.
			assert.Equal(t2, puzzle.NewCandidate(1), cells[0].Candidates())
			assert.Equal(t2, puzzle.NewCandidate(2), cells[1].Candidates())
			assert.Equal(t2, puzzle.NewCandidate(3), cells[2].Candidates())
			assert.Equal(t2, puzzle.NewCandidate(4), cells[3].Candidates())
			for i := uint8(5); i < 16; i++ {
				assert.Equalf(t2, puzzle.NewCandidate(i%4+1), cells[i].Candidates(), "cell %d", i)
			}
		},
	)

	t.Run(
		"does not panic when an out-of-bounds position is provided",
		func(t2 *testing.T) {
			grid := newGrid(
				slices.Repeat([][]uint8{{1, 2, 3, 4}}, 4),
				Must(puzzle.NewLayoutFromCellCount(16)),
			)

			assert.NotPanics(t2, func() {
				grid.Set(puzzle.NewPosition(0, 2), puzzle.NewCandidate(9))
			})
		},
	)
}

func TestGridWith(t *testing.T) {
	original := newGrid(
		slices.Repeat([][]uint8{{1, 2, 3, 4}}, 4),
		Must(puzzle.NewLayoutFromCellCount(16)),
	)

	t.Run(
		"cell is updated after copy",
		func(t2 *testing.T) {
			modified := original.With(puzzle.NewPosition(0, 0), puzzle.NewCandidate(9))

			originalCells := slices.Collect(original.Cells())
			modifiedCells := slices.Collect(modified.Cells())
			// The returned grid holds the new value.
			assert.Equal(t2, puzzle.NewCandidate(9), modifiedCells[0].Candidates())
			// The original is left untouched at the changed position.
			assert.Equal(t2, puzzle.NewCandidate(1), originalCells[0].Candidates())
			// Every other cell is shared/equal between the two grids.
			for i := 1; i < 4; i++ {
				assert.Equal(t2, originalCells[i].Candidates(), modifiedCells[i].Candidates())
			}
		},
	)

	t.Run(
		"noop when an out-of-bounds position is provided",
		func(t2 *testing.T) {
			modified := original.With(puzzle.NewPosition(0, 4), puzzle.NewCandidate(9))
			assert.Equal(t2, original, modified)
		},
	)

}

func TestGridClone(t *testing.T) {
	rows := slices.Repeat([][]uint8{{1, 2, 3, 4}}, 4)
	layout := Must(puzzle.NewLayoutFromCellCount(16))

	t.Run("mutating the clone leaves the original unchanged", func(t2 *testing.T) {
		original := newGrid(rows, layout)
		clone := original.Clone()

		clone.Set(puzzle.NewPosition(0, 0), puzzle.NewCandidate(2))

		originalCells := slices.Collect(original.Cells())
		cloneCells := slices.Collect(clone.Cells())
		assert.Equal(t2, puzzle.NewCandidate(1), originalCells[0].Candidates())
		assert.Equal(t2, puzzle.NewCandidate(2), cloneCells[0].Candidates())
	})

	t.Run("mutating the original leaves the clone unchanged", func(t2 *testing.T) {
		original := newGrid(rows, layout)
		clone := original.Clone()

		original.Set(puzzle.NewPosition(0, 0), puzzle.NewCandidate(2))

		originalCells := slices.Collect(original.Cells())
		cloneCells := slices.Collect(clone.Cells())
		assert.Equal(t2, puzzle.NewCandidate(2), originalCells[0].Candidates())
		assert.Equal(t2, puzzle.NewCandidate(1), cloneCells[0].Candidates())
	})
}

func TestGridPeersOf(t *testing.T) {
	testCases := []struct {
		name     string
		layout   puzzle.Layout
		pos      puzzle.Position
		expected []puzzle.Position
	}{
		{
			name:   "4x4 corner (0,0)",
			layout: Must(puzzle.NewLayoutFromCellCount(16)),
			pos:    puzzle.NewPosition(0, 0),
			expected: []puzzle.Position{
				puzzle.NewPosition(0, 1), puzzle.NewPosition(0, 2), puzzle.NewPosition(0, 3), // row
				puzzle.NewPosition(1, 0), puzzle.NewPosition(2, 0), puzzle.NewPosition(3, 0), // column
				puzzle.NewPosition(1, 1), // block
			},
		},
		{
			name:   "6x6 with 2x3 block (2,2)",
			layout: Must(puzzle.NewLayoutFromCellCount(36)),
			pos:    puzzle.NewPosition(2, 2),
			expected: []puzzle.Position{
				puzzle.NewPosition(2, 0), puzzle.NewPosition(2, 1), puzzle.NewPosition(2, 3), puzzle.NewPosition(2, 4), puzzle.NewPosition(2, 5), // row
				puzzle.NewPosition(0, 2), puzzle.NewPosition(1, 2), puzzle.NewPosition(3, 2), puzzle.NewPosition(4, 2), puzzle.NewPosition(5, 2), // column
				puzzle.NewPosition(3, 0), puzzle.NewPosition(3, 1), // block
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t2 *testing.T) {
				rows := make([][]uint8, testCase.layout.GridSize())
				for i := range testCase.layout.GridSize() {
					rows[i] = make([]uint8, testCase.layout.GridSize())
				}
				grid := newGrid(rows, testCase.layout)

				positions := make([]puzzle.Position, 0, len(testCase.expected))
				for cell := range grid.PeersOf(testCase.pos) {
					positions = append(positions, cell.Position())
				}

				assert.ElementsMatch(t2, testCase.expected, positions)
			},
		)
	}
}

func newGrid(rows [][]uint8, layout puzzle.Layout) puzzle.Grid {
	cells := make([]puzzle.Cell, 0, int(layout.GridSize())*int(layout.GridSize()))
	for row, rowValues := range rows {
		for col, value := range rowValues {
			position := puzzle.NewPosition(row, col)
			cells = append(cells, puzzle.NewCell(position, puzzle.NewCandidate(value)))
		}
	}
	return Must(puzzle.NewGrid(cells, layout))
}
