package puzzle

import (
	"iter"
)

type Peers struct {
	peers map[Position][]Position
}

func NewPeers(grid GridSize, region RegionSize) Peers {
	peers := make(map[Position][]Position, grid.RowCount()*grid.ColCount())

	regionPeerCount := region.RowCount()*region.ColCount() - 1
	rowPeerCount := grid.RowCount() - region.RowCount()
	colPeerCount := grid.ColCount() - region.ColCount()
	peerCount := regionPeerCount + rowPeerCount + colPeerCount

	for this := range positions(grid) {
		thisPeers := make([]Position, 0, peerCount)

		for that := range positions(grid) {
			if arePeers(this, that, region) {
				thisPeers = append(thisPeers, that)
			}
		}

		peers[this] = thisPeers
	}
	return Peers{peers: peers}
}

func arePeers(this, that Position, region RegionSize) bool {
	if this == that {
		return false
	}

	if this.Row() == that.Row() {
		return true
	}

	if this.Col() == that.Col() {
		return true
	}

	thisRegion := NewPosition(this.Row()/region.RowCount(), this.Col()/region.ColCount())
	thatRegion := NewPosition(that.Row()/region.RowCount(), that.Col()/region.ColCount())
	return thisRegion == thatRegion
}

func positions(grid GridSize) iter.Seq[Position] {
	return func(yield func(Position) bool) {
		for row := range grid.RowCount() {
			for col := range grid.ColCount() {
				if !yield(NewPosition(row, col)) {
					return
				}
			}
		}
	}
}

func (p Peers) Of(pos Position) []Position {
	return p.peers[pos]
}
