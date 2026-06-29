package puzzle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestCellWith(t *testing.T) {
	cell := puzzle.NewCell(puzzle.NewPosition(1, 2), puzzle.NewSingleCandidate(5))
	replaced := cell.With(puzzle.NewSingleCandidate(7))

	// The position is preserved and the candidates are replaced.
	assert.Equal(t, puzzle.NewPosition(1, 2), replaced.Position())
	assert.Equal(t, puzzle.NewSingleCandidate(7), replaced.Candidates())

	// The original cell is left unchanged.
	assert.Equal(t, puzzle.NewSingleCandidate(5), cell.Candidates())
}
