package puzzle_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expected      puzzle.Grid
		expectedError error
	}{
		{
			name:          "empty",
			input:         "",
			expected:      puzzle.Grid{},
			expectedError: puzzle.ErrInvalidCellCount,
		},
		{
			name:          "too short",
			input:         "123",
			expected:      puzzle.Grid{},
			expectedError: puzzle.ErrInvalidCellCount,
		},
		{
			name:          "too large",
			input:         strings.Repeat(".", 255),
			expected:      puzzle.Grid{},
			expectedError: puzzle.ErrInvalidCellCount,
		},
		{
			name:          "unexpectedly large value",
			input:         "9234123412341234", // 9 is unexpected for 4x4 grid
			expected:      puzzle.Grid{},
			expectedError: puzzle.ErrInvalidCharacter,
		},
		{
			name:          "unexpected value",
			input:         "z234123412341234", // z is unexpected
			expected:      puzzle.Grid{},
			expectedError: puzzle.ErrInvalidCharacter,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t2 *testing.T) {
				grid, err := puzzle.Parse(testCase.input)
				if testCase.expectedError != nil {
					require.ErrorIs(t2, err, testCase.expectedError)
				}
				assert.Equal(t2, testCase.expected, grid)
			},
		)
	}
}

func TestParseStringRoundTrip(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{name: "4x4 all missing", input: strings.Repeat(".", 16)},
		{name: "4x4", input: strings.Repeat("1234", 4)},
		{name: "9x9", input: strings.Repeat("123456789", 9)},
		{name: "12x12 with hex digits", input: strings.Repeat("123456789abc", 12)},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t2 *testing.T) {
				grid, err := puzzle.Parse(testCase.input)
				require.NoError(t2, err)
				assert.Equal(t2, testCase.input, grid.String())
			},
		)
	}
}

func TestGridRender(t *testing.T) {
	testCases := []struct {
		name     string
		grid     [][]uint8
		layout   puzzle.Layout
		expected string
	}{
		{
			name: "4x4 grid, 2x2 blocks",
			grid: [][]uint8{
				{1, 2, 3, 4},
				{3, 4, 1, 2},
				{4, 3, 2, 1},
				{2, 1, 4, 0},
			},
			layout: Must(puzzle.NewLayoutFromCellCount(16)),
			expected: ("+-----+-----+\n" +
				"| 1 2 | 3 4 |\n" +
				"| 3 4 | 1 2 |\n" +
				"+-----+-----+\n" +
				"| 4 3 | 2 1 |\n" +
				"| 2 1 | 4 . |\n" +
				"+-----+-----+\n" +
				"4-by-4 grid with 4 2-by-2 blocks"),
		},
		{
			name: "6x6 grid, 2x3 blocks",
			grid: [][]uint8{
				{0, 2, 3, 4, 5, 6},
				{4, 5, 6, 1, 2, 3},
				{2, 3, 1, 5, 6, 4},
				{5, 6, 4, 2, 3, 1},
				{3, 1, 2, 6, 4, 5},
				{6, 4, 5, 3, 1, 2},
			},
			layout: Must(puzzle.NewLayoutFromCellCount(36)),
			expected: ("+-------+-------+\n" +
				"| . 2 3 | 4 5 6 |\n" +
				"| 4 5 6 | 1 2 3 |\n" +
				"+-------+-------+\n" +
				"| 2 3 1 | 5 6 4 |\n" +
				"| 5 6 4 | 2 3 1 |\n" +
				"+-------+-------+\n" +
				"| 3 1 2 | 6 4 5 |\n" +
				"| 6 4 5 | 3 1 2 |\n" +
				"+-------+-------+\n" +
				"6-by-6 grid with 6 2-by-3 blocks"),
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t2 *testing.T) {
				grid := newGrid(testCase.grid, testCase.layout)
				assert.Equal(t2, testCase.expected, grid.Render())
			},
		)
	}
}

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
