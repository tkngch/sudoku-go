package puzzle

import (
	"slices"
)

// Peers holds the positions that share a row, column, or block with a given
// cell.
type Peers struct {
	row      PositionList
	column   PositionList
	block    PositionList
	allPeers PositionList
}

// PositionList is a read-only view over a precomputed peer slice: length and
// indexed access, with no exported way to reach or mutate the backing array.
type PositionList struct{ s []Position }

// NewPositionList clones positions and returns PositionList. So the returned
// list would not change when the input slice changes.
func NewPositionList(positions []Position) PositionList {
	return PositionList{slices.Clone(positions)}
}

// Len returns the number of elements in positionList.
func (l PositionList) Len() int { return len(l.s) }

// At returns Position at idx.
func (l PositionList) At(idx int) Position { return l.s[idx] }

// NewEmptyPeers returns a Peers with no members.
func NewEmptyPeers() Peers {
	return NewPeers([]Position{}, []Position{}, []Position{})
}

// NewPeers returns a Peers built from the given row, column, and block peer
// positions.
func NewPeers(rowPeers, colPeers, blockPeers []Position) Peers {
	allPeers := make([]Position, 0, len(rowPeers)+len(colPeers)+len(blockPeers))

	included := make(map[Position]bool)
	for _, item := range slices.Concat(rowPeers, colPeers, blockPeers) {
		if _, isIncluded := included[item]; isIncluded {
			continue
		}

		allPeers = append(allPeers, item)
		included[item] = true
	}

	return Peers{
		row:      NewPositionList(rowPeers),
		column:   NewPositionList(colPeers),
		block:    NewPositionList(blockPeers),
		allPeers: NewPositionList(allPeers),
	}
}

// All returns a positionList with the peers. Duplicates are removed from the
// three peers (row, column, and block).
func (p Peers) All() PositionList {
	return p.allPeers
}

// Row returns a positionList with the row peers.
func (p Peers) Row() PositionList {
	return p.row
}

// Col returns a positionList with the column peers.
func (p Peers) Col() PositionList {
	return p.column
}

// Block returns a positionList with the block peers.
func (p Peers) Block() PositionList {
	return p.block
}
