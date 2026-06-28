package puzzle

import (
	"iter"
	"slices"
)

type Peers struct {
	peers map[Position][]Position
}

func NewPeers(layout Layout) Peers {
	peers := make(map[Position][]Position, int(layout.GridSize())*int(layout.GridSize()))

	peerCount := layout.PeerCount()

	for this := range positions(layout.GridSize()) {
		thisPeers := make([]Position, 0, peerCount)

		for that := range positions(layout.GridSize()) {
			if arePeers(this, that, layout) {
				thisPeers = append(thisPeers, that)
			}
		}

		peers[this] = thisPeers
	}
	return Peers{peers: peers}
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

func positions(gridSize uint) iter.Seq[Position] {
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
	// return an iterator, to prevent the caller from altering the peers.
	return slices.Values(p.peers[pos])
}
