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
		name        string
		input       string
		expected    puzzle.Grid
		errExpected bool
	}{
		{name: "empty", input: "", expected: puzzle.Grid{}, errExpected: true},
		{name: "too short", input: "123", expected: puzzle.Grid{}, errExpected: true},
		{
			name:        "too large",
			input:       strings.Repeat(".", 255),
			expected:    puzzle.Grid{},
			errExpected: true,
		},
		{
			name:        "unexpectedly large value",
			input:       "9234123412341234", // 9 is unexpected for 4x4 grid
			expected:    puzzle.Grid{},
			errExpected: true,
		},
		{
			name:        "unexpected value",
			input:       "z234123412341234", // z is unexpected
			expected:    puzzle.Grid{},
			errExpected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t2 *testing.T) {
				grid, err := puzzle.Parse(testCase.input)
				if testCase.errExpected {
					require.Error(t2, err)
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
		{name: "4x4", input: strings.Repeat("1234", 4)},
		{name: "9x9 all givens", input: strings.Repeat("123456789", 9)},
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
		name       string
		grid       [][]uint8
		gridSize   puzzle.GridSize
		regionSize puzzle.RegionSize
		expected   string
	}{
		{
			name: "4x4 grid, 2x2 region",
			grid: [][]uint8{
				{1, 2, 3, 4},
				{3, 4, 1, 2},
				{4, 3, 2, 1},
				{2, 1, 4, 0},
			},
			gridSize:   puzzle.NewGridSize(4),
			regionSize: puzzle.NewRegionSize(2, 2),
			expected: ("+-----+-----+\n" +
				"| 1 2 | 3 4 |\n" +
				"| 3 4 | 1 2 |\n" +
				"+-----+-----+\n" +
				"| 4 3 | 2 1 |\n" +
				"| 2 1 | 4 . |\n" +
				"+-----+-----+\n" +
				"Grid Size: 4 by 4\n" +
				"Region Size: 2 by 2"),
		},
		{
			name: "6x6 grid, 2x3 region",
			grid: [][]uint8{
				{0, 2, 3, 4, 5, 6},
				{4, 5, 6, 1, 2, 3},
				{2, 3, 1, 5, 6, 4},
				{5, 6, 4, 2, 3, 1},
				{3, 1, 2, 6, 4, 5},
				{6, 4, 5, 3, 1, 2},
			},
			gridSize:   puzzle.NewGridSize(6),
			regionSize: puzzle.NewRegionSize(2, 3),
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
				"Grid Size: 6 by 6\n" +
				"Region Size: 2 by 3"),
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t2 *testing.T) {
				grid := newGrid(testCase.grid, testCase.gridSize, testCase.regionSize)
				assert.Equal(t2, testCase.expected, grid.Render())
			},
		)
	}
}
