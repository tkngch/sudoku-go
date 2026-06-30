package puzzle

import (
	"iter"
	"slices"
)

type Peers struct {
	row      []Position
	column   []Position
	block    []Position
	allPeers []Position
}

func NewEmptyPeers() Peers {
	return NewPeers([]Position{}, []Position{}, []Position{})
}

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

func (p Peers) Row() iter.Seq[Position] {
	return slices.Values(p.row)
}

func (p Peers) Col() iter.Seq[Position] {
	return slices.Values(p.column)
}

func (p Peers) Block() iter.Seq[Position] {
	return slices.Values(p.block)
}
