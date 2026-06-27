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
		cellCount   int
		expected    puzzle.GridSize
		errExpected bool
	}{
		{cellCount: 144, expected: puzzle.NewGridSize(12), errExpected: false},
		{cellCount: 81, expected: puzzle.NewGridSize(9), errExpected: false},
		{cellCount: 36, expected: puzzle.NewGridSize(6), errExpected: false},
		{cellCount: 16, expected: puzzle.NewGridSize(4), errExpected: false},
		{cellCount: 99, expected: puzzle.NewGridSize(0), errExpected: true},
		{cellCount: 0, expected: puzzle.NewGridSize(0), errExpected: true},
		{cellCount: -16, expected: puzzle.NewGridSize(0), errExpected: true},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("NewGridSizeFromCellCount(%d)", testCase.cellCount),
			func(t2 *testing.T) {
				size, err := puzzle.NewGridSizeFromCellCount(testCase.cellCount)
				if testCase.errExpected {
					require.Error(t2, err)
				} else {
					require.NoError(t2, err)
				}
				assert.Equal(t2, testCase.expected, size)
			},
		)
	}
}

func TestNewRegionSizeFromCellCount(t *testing.T) {
	testCases := []struct {
		cellCount   int
		expected    puzzle.RegionSize
		errExpected bool
	}{
		{cellCount: 144, expected: puzzle.NewRegionSize(4, 3), errExpected: false},
		{cellCount: 81, expected: puzzle.NewRegionSize(3, 3), errExpected: false},
		{cellCount: 36, expected: puzzle.NewRegionSize(2, 3), errExpected: false},
		{cellCount: 16, expected: puzzle.NewRegionSize(2, 2), errExpected: false},
		{cellCount: 99, expected: puzzle.NewRegionSize(0, 0), errExpected: true},
		{cellCount: 0, expected: puzzle.NewRegionSize(0, 0), errExpected: true},
		{cellCount: -16, expected: puzzle.NewRegionSize(0, 0), errExpected: true},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("NewRegionSizeFromCellCount(%d)", testCase.cellCount),
			func(t2 *testing.T) {
				size, err := puzzle.NewRegionSizeFromCellCount(testCase.cellCount)
				if testCase.errExpected {
					require.Error(t2, err)
				} else {
					require.NoError(t2, err)
				}
				assert.Equal(t2, testCase.expected, size)
			},
		)
	}
}
