package puzzle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestPositionString(t *testing.T) {
	pos := puzzle.NewPosition(3, 5)

	assert.Contains(t, pos.String(), "3")
	assert.Contains(t, pos.String(), "5")
}
