package puzzle_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestNewGridSizeFromCellCount(t *testing.T) {
	testCases := []struct {
		cellCount int
		expected  puzzle.GridSize
	}{
		{144, puzzle.NewGridSize(12)},
		{81, puzzle.NewGridSize(9)},
		{36, puzzle.NewGridSize(6)},
		{16, puzzle.NewGridSize(4)},
		{99, puzzle.NewGridSize(0)},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("NewGridSizeFromCellCount(%d)", testCase.cellCount),
			func(t2 *testing.T) {
				size, err := puzzle.NewGridSizeFromCellCount(testCase.cellCount)
				if testCase.expected.RowCount() != 0 {
					require.NoError(t2, err)
				}
				assert.Equal(t2, testCase.expected, size)
			},
		)
	}
}

func TestNewRegionSizeFromCellCount(t *testing.T) {
	testCases := []struct {
		cellCount int
		expected  puzzle.RegionSize
	}{
		{144, puzzle.NewRegionSize(4, 3)},
		{81, puzzle.NewRegionSize(3, 3)},
		{36, puzzle.NewRegionSize(2, 3)},
		{16, puzzle.NewRegionSize(2, 2)},
		{99, puzzle.NewRegionSize(0, 0)},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("NewGridSizeFromCellCount(%d)", testCase.cellCount),
			func(t2 *testing.T) {
				size, err := puzzle.NewRegionSizeFromCellCount(testCase.cellCount)
				if testCase.expected.RowCount() != 0 {
					require.NoError(t2, err)
				}
				assert.Equal(t2, testCase.expected, size)
			},
		)
	}
}
