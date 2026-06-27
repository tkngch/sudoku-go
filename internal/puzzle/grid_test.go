package puzzle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestGridSet(t *testing.T) {
	t.Run(
		"cell is updated in place",
		func(t2 *testing.T) {
			grid := newGrid(
				[][]uint8{{1, 2}, {2, 1}},
				puzzle.NewGridSize(2),
				puzzle.NewRegionSize(1, 2),
			)

			newCandidate := puzzle.NewCandidate(9)
			grid.Set(puzzle.NewPosition(1, 0), newCandidate)

			cells := grid.Cells()
			assert.Equal(t2, newCandidate, cells[2].Candidates(), "cell is updated in place")
			// Every other cell is untouched.
			assert.Equal(t2, puzzle.NewCandidate(1), cells[0].Candidates())
			assert.Equal(t2, puzzle.NewCandidate(2), cells[1].Candidates())
			assert.Equal(t2, puzzle.NewCandidate(1), cells[3].Candidates())
		},
	)

	t.Run(
		"does not panic when an out-of-bounds position is provided",
		func(t2 *testing.T) {
			grid := newGrid(
				[][]uint8{{1, 2}, {2, 1}},
				puzzle.NewGridSize(2),
				puzzle.NewRegionSize(1, 2),
			)

			assert.NotPanics(t2, func() {
				grid.Set(puzzle.NewPosition(0, 2), puzzle.NewCandidate(9))
			})
		},
	)
}

func TestGridWith(t *testing.T) {
	original := newGrid(
		[][]uint8{{1, 2}, {2, 1}},
		puzzle.NewGridSize(2),
		puzzle.NewRegionSize(1, 2),
	)

	t.Run(
		"cell is updated after copy",
		func(t2 *testing.T) {
			modified := original.With(puzzle.NewPosition(0, 0), puzzle.NewCandidate(9))

			// The returned grid holds the new value.
			assert.Equal(t2, puzzle.NewCandidate(9), modified.Cells()[0].Candidates())
			// The original is left untouched at the changed position.
			assert.Equal(t2, puzzle.NewCandidate(1), original.Cells()[0].Candidates())
			// Every other cell is shared/equal between the two grids.
			for i := 1; i < 4; i++ {
				assert.Equal(t2, original.Cells()[i].Candidates(), modified.Cells()[i].Candidates())
			}
		},
	)

	t.Run(
		"noop when an out-of-bounds position is provided",
		func(t2 *testing.T) {
			modified := original.With(puzzle.NewPosition(0, 2), puzzle.NewCandidate(9))
			assert.Equal(t2, original, modified)
		},
	)

}

func TestGridClone(t *testing.T) {
	rows := [][]uint8{{1}}
	gridSize := puzzle.NewGridSize(1)
	regionSize := puzzle.NewRegionSize(1, 1)

	t.Run("mutating the clone leaves the original unchanged", func(t2 *testing.T) {
		original := newGrid(rows, gridSize, regionSize)
		clone := original.Clone()

		clone.Set(puzzle.NewPosition(0, 0), puzzle.NewCandidate(2))

		assert.Equal(t2, puzzle.NewCandidate(1), original.Cells()[0].Candidates())
		assert.Equal(t2, puzzle.NewCandidate(2), clone.Cells()[0].Candidates())
	})

	t.Run("mutating the original leaves the clone unchanged", func(t2 *testing.T) {
		original := newGrid(rows, gridSize, regionSize)
		clone := original.Clone()

		original.Set(puzzle.NewPosition(0, 0), puzzle.NewCandidate(2))

		assert.Equal(t2, puzzle.NewCandidate(2), original.Cells()[0].Candidates())
		assert.Equal(t2, puzzle.NewCandidate(1), clone.Cells()[0].Candidates())
	})
}

func TestGridGetPeers(t *testing.T) {
	testCases := []struct {
		name       string
		gridSize   puzzle.GridSize
		regionSize puzzle.RegionSize
		pos        puzzle.Position
		expected   []puzzle.Position
	}{
		{
			name:       "4x4 corner (0,0)",
			gridSize:   puzzle.NewGridSize(4),
			regionSize: puzzle.NewRegionSize(2, 2),
			pos:        puzzle.NewPosition(0, 0),
			expected: []puzzle.Position{
				puzzle.NewPosition(0, 1), puzzle.NewPosition(0, 2), puzzle.NewPosition(0, 3), // row
				puzzle.NewPosition(1, 0), puzzle.NewPosition(2, 0), puzzle.NewPosition(3, 0), // column
				puzzle.NewPosition(1, 1), // region
			},
		},
		{
			name:       "6x6 non-square region interior (2,2)",
			gridSize:   puzzle.NewGridSize(6),
			regionSize: puzzle.NewRegionSize(2, 3),
			pos:        puzzle.NewPosition(2, 2),
			expected: []puzzle.Position{
				puzzle.NewPosition(2, 0), puzzle.NewPosition(2, 1), puzzle.NewPosition(2, 3), puzzle.NewPosition(2, 4), puzzle.NewPosition(2, 5), // row
				puzzle.NewPosition(0, 2), puzzle.NewPosition(1, 2), puzzle.NewPosition(3, 2), puzzle.NewPosition(4, 2), puzzle.NewPosition(5, 2), // column
				puzzle.NewPosition(3, 0), puzzle.NewPosition(3, 1), // region
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t2 *testing.T) {
				rows := make([][]uint8, testCase.gridSize.RowCount())
				for i := range testCase.gridSize.RowCount() {
					rows[i] = make([]uint8, testCase.gridSize.ColCount())
				}
				grid := newGrid(rows, testCase.gridSize, testCase.regionSize)

				positions := make([]puzzle.Position, 0, len(testCase.expected))
				for cell := range grid.GetPeers(testCase.pos) {
					positions = append(positions, cell.Position())
				}

				assert.ElementsMatch(t2, testCase.expected, positions)
			},
		)
	}
}

func newGrid(rows [][]uint8, grid puzzle.GridSize, region puzzle.RegionSize) puzzle.Grid {
	cells := make([]puzzle.Cell, 0, grid.RowCount()*grid.ColCount())
	for row, rowValues := range rows {
		for col, value := range rowValues {
			position := puzzle.NewPosition(uint8(row), uint8(col))
			cells = append(cells, puzzle.NewCell(position, puzzle.NewCandidate(value)))
		}
	}
	return puzzle.NewGrid(cells, grid, region)
}
