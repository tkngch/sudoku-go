package puzzle_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestNewPeers(t *testing.T) {
	t.Parallel()

	rowPeers := []puzzle.Position{
		puzzle.NewPosition(0, 1),
		puzzle.NewPosition(0, 2),
	}
	colPeers := []puzzle.Position{
		puzzle.NewPosition(1, 0),
		puzzle.NewPosition(2, 0),
	}
	// (0,1) also appears in rowPeers and (1,0) in colPeers, so All must report
	// each of them once.
	blockPeers := []puzzle.Position{
		puzzle.NewPosition(0, 1),
		puzzle.NewPosition(1, 0),
		puzzle.NewPosition(1, 1),
	}

	peers := puzzle.NewPeers(rowPeers, colPeers, blockPeers)

	assert.Equal(t, rowPeers, slices.Collect(peers.Row()), "row peers")
	assert.Equal(t, colPeers, slices.Collect(peers.Col()), "column peers")
	assert.Equal(t, blockPeers, slices.Collect(peers.Block()), "block peers")

	expected := []puzzle.Position{
		puzzle.NewPosition(0, 1),
		puzzle.NewPosition(0, 2),
		puzzle.NewPosition(1, 0),
		puzzle.NewPosition(2, 0),
		puzzle.NewPosition(1, 1),
	}
	actual := slices.Collect(peers.All())

	assert.ElementsMatch(t, expected, actual, "deduplicated peers")
	assert.Len(t, actual, len(expected), "All must not contain duplicates")
}

func TestNewEmptyPeers(t *testing.T) {
	t.Parallel()

	peers := puzzle.NewEmptyPeers()

	assert.Empty(t, slices.Collect(peers.Row()))
	assert.Empty(t, slices.Collect(peers.Col()))
	assert.Empty(t, slices.Collect(peers.Block()))
	assert.Empty(t, slices.Collect(peers.All()))
}
