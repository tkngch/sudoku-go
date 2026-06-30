package puzzle_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestNewLayoutFromCellCount(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		input            int
		expectedError    error
		expectedGridSize int
	}{
		{name: "16 cells", input: 16, expectedGridSize: 4},
		{name: "36 cells", input: 36, expectedGridSize: 6},
		{name: "81 cells", input: 81, expectedGridSize: 9},
		{name: "144 cells", input: 144, expectedGridSize: 12},
		{name: "zero", input: 0, expectedError: puzzle.ErrInvalidCellCount},
		{name: "255", input: 255, expectedError: puzzle.ErrInvalidCellCount},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				layout, err := puzzle.NewLayoutFromCellCount(testCase.input)

				if testCase.expectedError != nil {
					require.Error(t, err)
					require.ErrorIs(t, err, testCase.expectedError)
					assert.Equal(t, puzzle.Layout{}, layout)

					return
				}

				require.NoError(t, err)
				assert.Equal(t, testCase.expectedGridSize, layout.GridSize())
			},
		)
	}
}

func TestLayoutIsOnGrid(t *testing.T) {
	t.Parallel()

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
				t.Parallel()

				assert.Equal(t, testCase.expected, layout.IsOnGrid(testCase.position))
			},
		)
	}
}

func TestLayoutRowMajorIndex(t *testing.T) {
	t.Parallel()

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
				t.Parallel()

				layout, err := puzzle.NewLayoutFromCellCount(testCase.cellCount)
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, layout.RowMajorIndex(testCase.position))
			},
		)
	}
}

func TestLayoutPeersOf(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		cellCount int
		pos       puzzle.Position
		expected  puzzle.Peers
	}{
		{
			name:      "4x4 corner (0,0)",
			cellCount: 16,
			pos:       puzzle.NewPosition(0, 0),
			expected: puzzle.NewPeers(
				[]puzzle.Position{
					puzzle.NewPosition(0, 1),
					puzzle.NewPosition(0, 2),
					puzzle.NewPosition(0, 3),
				},
				[]puzzle.Position{
					puzzle.NewPosition(1, 0),
					puzzle.NewPosition(2, 0),
					puzzle.NewPosition(3, 0),
				},
				[]puzzle.Position{
					puzzle.NewPosition(0, 1),
					puzzle.NewPosition(1, 0),
					puzzle.NewPosition(1, 1),
				},
			),
		},
		{
			name:      "6x6 with 2x3 block (2,2)",
			cellCount: 36,
			pos:       puzzle.NewPosition(2, 2),
			expected: puzzle.NewPeers(
				[]puzzle.Position{
					puzzle.NewPosition(2, 0),
					puzzle.NewPosition(2, 1),
					puzzle.NewPosition(2, 3),
					puzzle.NewPosition(2, 4),
					puzzle.NewPosition(2, 5),
				},
				[]puzzle.Position{
					puzzle.NewPosition(0, 2),
					puzzle.NewPosition(1, 2),
					puzzle.NewPosition(3, 2),
					puzzle.NewPosition(4, 2),
					puzzle.NewPosition(5, 2),
				},
				[]puzzle.Position{
					puzzle.NewPosition(2, 0),
					puzzle.NewPosition(2, 1),
					puzzle.NewPosition(3, 2),
					puzzle.NewPosition(3, 0),
					puzzle.NewPosition(3, 1),
				},
			),
		},
		{
			name:      "4x4 outside the grid",
			cellCount: 16,
			pos:       puzzle.NewPosition(0, 4),
			expected:  puzzle.NewEmptyPeers(),
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				layout, err := puzzle.NewLayoutFromCellCount(testCase.cellCount)
				require.NoError(t, err)

				peers := layout.PeersOf(testCase.pos)
				assert.ElementsMatch(
					t,
					slices.Collect(testCase.expected.Row()),
					slices.Collect(peers.Row()),
					"row peers",
				)
				assert.ElementsMatch(
					t,
					slices.Collect(testCase.expected.Col()),
					slices.Collect(peers.Col()),
					"column peers",
				)
				assert.ElementsMatch(
					t,
					slices.Collect(testCase.expected.Block()),
					slices.Collect(peers.Block()),
					"block peers",
				)
			},
		)
	}
}
