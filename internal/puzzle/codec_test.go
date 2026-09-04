package puzzle_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestParseErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "empty",
			input:         "",
			expectedError: puzzle.ErrInvalidCellCount,
		},
		{
			name:          "whitespace only",
			input:         " \t\n ",
			expectedError: puzzle.ErrInvalidCellCount,
		},
		{
			name:          "too short",
			input:         "123",
			expectedError: puzzle.ErrInvalidCellCount,
		},
		{
			name:          "too large",
			input:         strings.Repeat(".", 255),
			expectedError: puzzle.ErrInvalidCellCount,
		},
		{
			name:          "too short after ignoring whitespace",
			input:         "234 123412341234", // 16 characters, but only 15 cells
			expectedError: puzzle.ErrInvalidCellCount,
		},
		{
			name:          "unexpectedly large value",
			input:         "9234123412341234", // 9 is unexpected for 4x4 grid
			expectedError: puzzle.ErrValueOutOfRange,
		},
		{
			name:          "unexpected value",
			input:         "z234123412341234", // z is unexpected
			expectedError: puzzle.ErrInvalidCharacter,
		},
		{
			// This input holds both faults. The character takes precedence,
			// because the cell count of a malformed input means little.
			name:          "unexpected value in a short input",
			input:         "z23",
			expectedError: puzzle.ErrInvalidCharacter,
		},
		{
			name:          "zero width space before a full grid",
			input:         "\u200B" + strings.Repeat("123456789", 9),
			expectedError: puzzle.ErrInvalidCharacter,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				grid, err := puzzle.Parse(testCase.input)
				require.ErrorIs(t, err, testCase.expectedError)
				assert.Nil(t, grid)
			},
		)
	}
}

func TestParseEquivalence(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		input      string
		equivalent string
	}{
		{
			name:       "multiline",
			input:      ".234\n3.12\n43.1\n214.\n",
			equivalent: ".234" + "3.12" + "43.1" + "214.",
		},
		{
			name:       "spaced rows",
			input:      ". 2 3 4  3 . 1 2  4 3 . 1  2 1 4 .",
			equivalent: ".234" + "3.12" + "43.1" + "214.",
		},
		{
			name:       "leading and trailing space",
			input:      "  " + ".234" + "3.12" + "43.1" + "214." + "\n",
			equivalent: ".234" + "3.12" + "43.1" + "214.",
		},
		{
			name:       "no-break space between rows",
			input:      ".234" + "\u00A0" + "3.12" + "43.1" + "214.",
			equivalent: ".234" + "3.12" + "43.1" + "214.",
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				grid, err := puzzle.Parse(testCase.input)
				require.NoError(t, err)

				expected, err := puzzle.Parse(testCase.equivalent)
				require.NoError(t, err)

				assert.Equal(t, expected, grid)
			},
		)
	}
}

func TestParseStringRoundTrip(t *testing.T) {
	t.Parallel()

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
			func(t *testing.T) {
				t.Parallel()

				grid, err := puzzle.Parse(testCase.input)
				require.NoError(t, err)
				assert.Equal(t, testCase.input, grid.String())
			},
		)
	}
}

func TestGridRender(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		grid      [][]int
		cellCount int
		expected  string
	}{
		{
			name: "4x4 grid, 2x2 blocks",
			grid: [][]int{
				{1, 2, 3, 4},
				{3, 4, 1, 2},
				{4, 3, 2, 1},
				{2, 1, 4, 0},
			},
			cellCount: 16,
			expected: ("+-----+-----+\n" +
				"| 1 2 | 3 4 |\n" +
				"| 3 4 | 1 2 |\n" +
				"+-----+-----+\n" +
				"| 4 3 | 2 1 |\n" +
				"| 2 1 | 4 . |\n" +
				"+-----+-----+"),
		},
		{
			name: "6x6 grid, 2x3 blocks",
			grid: [][]int{
				{0, 2, 3, 4, 5, 6},
				{4, 5, 6, 1, 2, 3},
				{2, 3, 1, 5, 6, 4},
				{5, 6, 4, 2, 3, 1},
				{3, 1, 2, 6, 4, 5},
				{6, 4, 5, 3, 1, 2},
			},
			cellCount: 36,
			expected: ("+-------+-------+\n" +
				"| . 2 3 | 4 5 6 |\n" +
				"| 4 5 6 | 1 2 3 |\n" +
				"+-------+-------+\n" +
				"| 2 3 1 | 5 6 4 |\n" +
				"| 5 6 4 | 2 3 1 |\n" +
				"+-------+-------+\n" +
				"| 3 1 2 | 6 4 5 |\n" +
				"| 6 4 5 | 3 1 2 |\n" +
				"+-------+-------+"),
		},
		{
			name:      "empty grid",
			grid:      [][]int{},
			cellCount: 0,
			expected:  "",
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				grid := newGrid(t, testCase.grid, newLayoutForCellCount(t, testCase.cellCount))
				assert.Equal(t, testCase.expected, grid.Render())
			},
		)
	}
}

func newLayoutForCellCount(t *testing.T, cellCount int) puzzle.Layout {
	t.Helper()

	if cellCount == 0 { // empty grid: a zero-value layout renders as ""
		return puzzle.Layout{}
	}

	layout, err := puzzle.NewLayoutForCellCount(cellCount)
	require.NoError(t, err)

	return layout
}
