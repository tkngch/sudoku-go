package puzzle_test

import (
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

	assertEqualPositionList(t, puzzle.NewPositionList(rowPeers), peers.Row(), "row peers")
	assertEqualPositionList(t, puzzle.NewPositionList(colPeers), peers.Col(), "column peers")
	assertEqualPositionList(t, puzzle.NewPositionList(blockPeers), peers.Block(), "block peers")

	expected := puzzle.NewPositionList(
		[]puzzle.Position{
			puzzle.NewPosition(0, 1),
			puzzle.NewPosition(0, 2),
			puzzle.NewPosition(1, 0),
			puzzle.NewPosition(2, 0),
			puzzle.NewPosition(1, 1),
		})
	actual := peers.All()

	assertEqualPositionList(t, expected, actual, "peers all")
}

func TestNewEmptyPeers(t *testing.T) {
	t.Parallel()

	peers := puzzle.NewEmptyPeers()

	assert.Equal(t, 0, peers.Row().Len(), "row length")
	assert.Equal(t, 0, peers.Col().Len(), "col length")
	assert.Equal(t, 0, peers.Block().Len(), "block length")
	assert.Equal(t, 0, peers.All().Len(), "all length")
}
