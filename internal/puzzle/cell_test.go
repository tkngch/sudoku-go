package puzzle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestCellReplace(t *testing.T) {
	cell := puzzle.NewCell(puzzle.NewPosition(1, 2), puzzle.NewCandidate(5))
	replaced := cell.Replace(puzzle.NewCandidate(7))

	// The position is preserved and the candidates are replaced.
	assert.Equal(t, puzzle.NewPosition(1, 2), replaced.Position())
	assert.Equal(t, puzzle.NewCandidate(7), replaced.Candidates())

	// The original cell is left unchanged.
	assert.Equal(t, puzzle.NewCandidate(5), cell.Candidates())
}
