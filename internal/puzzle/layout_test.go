package puzzle_test

import (
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
			func(t2 *testing.T) {
				layout, err := puzzle.NewLayoutFromCellCount(testCase.input)

				if testCase.expectedError != nil {
					require.Error(t2, err)
					assert.ErrorIs(t2, err, testCase.expectedError)
					assert.Equal(t2, puzzle.Layout{}, layout)
					return
				}

				require.NoError(t2, err)
				assert.Equal(t2, testCase.expectedGridSize, layout.GridSize())
			},
		)
	}
}

func TestLayoutAreInSameBlock(t *testing.T) {
	testCases := []struct {
		name      string
		cellCount int
		this      puzzle.Position
		that      puzzle.Position
		expected  bool
	}{
		{
			name:      "4x4 same cell",
			cellCount: 16,
			this:      puzzle.NewPosition(0, 0),
			that:      puzzle.NewPosition(0, 0),
			expected:  true,
		},
		{
			name:      "4x4 same row, same block",
			cellCount: 16,
			this:      puzzle.NewPosition(0, 0),
			that:      puzzle.NewPosition(0, 1),
			expected:  true,
		},
		{
			name:      "4x4 diagonal, same block",
			cellCount: 16,
			this:      puzzle.NewPosition(0, 0),
			that:      puzzle.NewPosition(1, 1),
			expected:  true,
		},
		{
			name:      "4x4 same row, different block",
			cellCount: 16,
			this:      puzzle.NewPosition(0, 0),
			that:      puzzle.NewPosition(0, 2),
			expected:  false,
		},
		{
			name:      "4x4 same col, different block",
			cellCount: 16,
			this:      puzzle.NewPosition(0, 0),
			that:      puzzle.NewPosition(2, 0),
			expected:  false,
		},
		{
			name:      "6x6 same block",
			cellCount: 36,
			this:      puzzle.NewPosition(0, 0),
			that:      puzzle.NewPosition(1, 2),
			expected:  true,
		},
		{
			name:      "6x6 block below",
			cellCount: 36,
			this:      puzzle.NewPosition(0, 0),
			that:      puzzle.NewPosition(2, 0),
			expected:  false,
		},
		{
			name:      "6x6 block to the right",
			cellCount: 36,
			this:      puzzle.NewPosition(0, 0),
			that:      puzzle.NewPosition(0, 3),
			expected:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t2 *testing.T) {
				layout := Must(puzzle.NewLayoutFromCellCount(testCase.cellCount))
				assert.Equal(t2, testCase.expected, layout.AreInSameBlock(testCase.this, testCase.that))
				// The relation is symmetric.
				assert.Equal(t2, testCase.expected, layout.AreInSameBlock(testCase.that, testCase.this))
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
			func(t2 *testing.T) {
				assert.Equal(t2, testCase.expected, layout.IsOnGrid(testCase.position))
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
			func(t2 *testing.T) {
				layout := Must(puzzle.NewLayoutFromCellCount(testCase.cellCount))
				assert.Equal(t2, testCase.expected, layout.RowMajorIndex(testCase.position))
			},
		)
	}
}
