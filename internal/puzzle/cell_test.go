package puzzle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestNewCell(t *testing.T) {
	position := puzzle.NewPosition(1, 2)
	candidates := puzzle.NewSingleCandidate(5)
	cell := puzzle.NewCell(position, candidates)

	assert.Equal(t, position, cell.Position())
	assert.Equal(t, candidates, cell.Candidates())

	assert.Contains(t, cell.String(), position.String())
	assert.Contains(t, cell.String(), candidates.String())
}
