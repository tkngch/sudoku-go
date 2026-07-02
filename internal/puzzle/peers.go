package puzzle

import (
	"iter"
	"slices"
)

// Peers holds the positions that share a row, column, or block with a given
// cell.
type Peers struct {
	row      []Position
	column   []Position
	block    []Position
	allPeers []Position
}

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
		row:      slices.Clone(rowPeers),
		column:   slices.Clone(colPeers),
		block:    slices.Clone(blockPeers),
		allPeers: allPeers,
	}
}

// All returns a single iterator over the peers. Duplicates are removed from the
// three peers (row, column, and block).
func (p Peers) All() iter.Seq[Position] {
	return slices.Values(p.allPeers)
}

// Row returns an iterator over the row peers.
func (p Peers) Row() iter.Seq[Position] {
	return slices.Values(p.row)
}

// Col returns an iterator over the column peers.
func (p Peers) Col() iter.Seq[Position] {
	return slices.Values(p.column)
}

// Block returns an iterator over the block peers.
func (p Peers) Block() iter.Seq[Position] {
	return slices.Values(p.block)
}
