package puzzle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestPosition(t *testing.T) {
	pos := puzzle.NewPosition(3, 5)

	assert.Equal(t, uint8(3), pos.Row())
	assert.Equal(t, uint8(5), pos.Col())
}
