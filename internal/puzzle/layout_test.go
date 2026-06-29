package puzzle_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestNewLayoutFromCellCount(t *testing.T) {
	testCases := []struct {
		name              string
		input             int
		expectedError     error
		expectedGridSize  int
		expectedPeerCount int
	}{
		{name: "16 cells", input: 16, expectedGridSize: 4, expectedPeerCount: 7},
		{name: "36 cells", input: 36, expectedGridSize: 6, expectedPeerCount: 12},
		{name: "81 cells", input: 81, expectedGridSize: 9, expectedPeerCount: 20},
		{name: "144 cells", input: 144, expectedGridSize: 12, expectedPeerCount: 28},
		{name: "zero", input: 0, expectedError: puzzle.ErrInvalidCellCount},
		{name: "255", input: 255, expectedError: puzzle.ErrInvalidCellCount},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				layout, err := puzzle.NewLayoutFromCellCount(testCase.input)

				if testCase.expectedError != nil {
					require.Error(t, err)
					assert.ErrorIs(t, err, testCase.expectedError)
					assert.Equal(t, puzzle.Layout{}, layout)
					return
				}

				peers := slices.Collect(layout.PeersOf(puzzle.NewPosition(0, 0)))

				require.NoError(t, err)
				assert.Equal(t, testCase.expectedGridSize, layout.GridSize())
				assert.Equal(t, testCase.expectedPeerCount, len(peers))
			},
		)
	}
}

func TestLayoutIsOnGrid(t *testing.T) {
	layout := Must(puzzle.NewLayoutFromCellCount(16)) // 4x4 grid
	testCases := []struct {
		name     string
		position puzzle.Position
		expected bool
	}{
		{name: "top left", position: puzzle.NewPosition(0, 0), expected: true},
		{name: "bottom right", position: puzzle.NewPosition(3, 3), expected: true},
		{name: "row out of range", position: puzzle.NewPosition(4, 0), expected: false},
		{name: "col out of range", position: puzzle.NewPosition(0, 4), expected: false},
		{name: "both out of range", position: puzzle.NewPosition(4, 4), expected: false},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				assert.Equal(t, testCase.expected, layout.IsOnGrid(testCase.position))
			},
		)
	}
}

func TestLayoutRowMajorIndex(t *testing.T) {
	testCases := []struct {
		name      string
		cellCount int
		position  puzzle.Position
		expected  int
	}{
		{
			name:      "4x4 top left",
			cellCount: 16,
			position:  puzzle.NewPosition(0, 0),
			expected:  0,
		},
		{
			name:      "4x4 top right",
			cellCount: 16,
			position:  puzzle.NewPosition(0, 3),
			expected:  3,
		},
		{
			name:      "4x4 second row",
			cellCount: 16,
			position:  puzzle.NewPosition(1, 0),
			expected:  4,
		},
		{
			name:      "4x4 bottom right",
			cellCount: 16,
			position:  puzzle.NewPosition(3, 3),
			expected:  15,
		},
		{
			name:      "9x9 third row",
			cellCount: 81,
			position:  puzzle.NewPosition(2, 3),
			expected:  21,
		},
		{
			name:      "9x9 bottom right",
			cellCount: 81,
			position:  puzzle.NewPosition(8, 8),
			expected:  80,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				layout := Must(puzzle.NewLayoutFromCellCount(testCase.cellCount))
				assert.Equal(t, testCase.expected, layout.RowMajorIndex(testCase.position))
			},
		)
	}
}
