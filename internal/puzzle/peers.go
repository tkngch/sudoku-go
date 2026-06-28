package puzzle

import (
	"iter"
	"slices"
)

type Peers struct {
	peers  [][]Position
	layout Layout
}

func NewPeers(layout Layout) Peers {
	peers := make([][]Position, layout.CellCount())
	peerCount := layout.PeerCount()

	for this := range positions(layout.GridSize()) {
		thisPeers := make([]Position, 0, peerCount)

		for that := range positions(layout.GridSize()) {
			if arePeers(this, that, layout) {
				thisPeers = append(thisPeers, that)
			}
		}

		peers[layout.RowMajorIndex(this)] = thisPeers
	}
	return Peers{peers: peers, layout: layout}
}

func arePeers(this, that Position, layout Layout) bool {
	if this == that {
		return false
	}

	if this.Row() == that.Row() {
		return true
	}

	if this.Col() == that.Col() {
		return true
	}

	return layout.AreInSameBlock(this, that)
}

func positions(gridSize int) iter.Seq[Position] {
	return func(yield func(Position) bool) {
		for row := range gridSize {
			for col := range gridSize {
				if !yield(NewPosition(row, col)) {
					return
				}
			}
		}
	}
}

// Return an iterator over the peers: the other cells sharing its row, column,
// or block.
func (p Peers) Of(pos Position) iter.Seq[Position] {
	if p.layout.IsOnGrid(pos) {
		// return an iterator, to prevent the caller from altering the peers.
		return slices.Values(p.peers[p.layout.RowMajorIndex(pos)])
	}
	return func(yield func(Position) bool) {}
}
