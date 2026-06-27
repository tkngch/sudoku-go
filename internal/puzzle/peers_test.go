package puzzle_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestPeersOf(t *testing.T) {
	testCases := []struct {
		name     string
		grid     puzzle.GridSize
		region   puzzle.RegionSize
		pos      puzzle.Position
		expected []puzzle.Position
	}{
		{
			name:   "4x4 corner (0,0)",
			grid:   puzzle.NewGridSize(4),
			region: puzzle.NewRegionSize(2, 2),
			pos:    puzzle.NewPosition(0, 0),
			expected: []puzzle.Position{
				puzzle.NewPosition(0, 1), puzzle.NewPosition(0, 2), puzzle.NewPosition(0, 3), // row
				puzzle.NewPosition(1, 0), puzzle.NewPosition(2, 0), puzzle.NewPosition(3, 0), // column
				puzzle.NewPosition(1, 1), // region
			},
		},
		{
			name:   "4x4 non-corner (1,2)",
			grid:   puzzle.NewGridSize(4),
			region: puzzle.NewRegionSize(2, 2),
			pos:    puzzle.NewPosition(1, 2),
			expected: []puzzle.Position{
				puzzle.NewPosition(1, 0), puzzle.NewPosition(1, 1), puzzle.NewPosition(1, 3), // row
				puzzle.NewPosition(0, 2), puzzle.NewPosition(2, 2), puzzle.NewPosition(3, 2), // column
				puzzle.NewPosition(0, 3), // region
			},
		},
		{
			name:   "6x6 non-square region (2,2)",
			grid:   puzzle.NewGridSize(6),
			region: puzzle.NewRegionSize(2, 3),
			pos:    puzzle.NewPosition(2, 2),
			expected: []puzzle.Position{
				puzzle.NewPosition(2, 0), puzzle.NewPosition(2, 1), puzzle.NewPosition(2, 3), puzzle.NewPosition(2, 4), puzzle.NewPosition(2, 5), // row
				puzzle.NewPosition(0, 2), puzzle.NewPosition(1, 2), puzzle.NewPosition(3, 2), puzzle.NewPosition(4, 2), puzzle.NewPosition(5, 2), // column
				puzzle.NewPosition(3, 0), puzzle.NewPosition(3, 1), // region
			},
		},
		{
			name:   "6x6 non-square region (0,0)",
			grid:   puzzle.NewGridSize(6),
			region: puzzle.NewRegionSize(2, 3),
			pos:    puzzle.NewPosition(0, 0),
			expected: []puzzle.Position{
				puzzle.NewPosition(0, 1), puzzle.NewPosition(0, 2), puzzle.NewPosition(0, 3), puzzle.NewPosition(0, 4), puzzle.NewPosition(0, 5), // row
				puzzle.NewPosition(1, 0), puzzle.NewPosition(2, 0), puzzle.NewPosition(3, 0), puzzle.NewPosition(4, 0), puzzle.NewPosition(5, 0), // column
				puzzle.NewPosition(1, 1), puzzle.NewPosition(1, 2), // region
			},
		},
		{
			name:     "4x4 invalid position",
			grid:     puzzle.NewGridSize(4),
			region:   puzzle.NewRegionSize(2, 2),
			pos:      puzzle.NewPosition(99, 99),
			expected: []puzzle.Position{},
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t2 *testing.T) {
				peers := puzzle.NewPeers(testCase.grid, testCase.region)
				actual := slices.Collect(peers.Of(testCase.pos))
				assert.ElementsMatch(t2, testCase.expected, actual)
			},
		)
	}
}
