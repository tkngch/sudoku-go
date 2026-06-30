package puzzle

import (
	"iter"
	"slices"
)

type Peers[T comparable] struct {
	row      []T
	column   []T
	block    []T
	allPeers []T
}

func NewEmptyPeers[T comparable]() Peers[T] {
	return NewPeers([]T{}, []T{}, []T{})
}

func NewPeers[T comparable](rowPeers, colPeers, blockPeers []T) Peers[T] {
	allPeers := make([]T, 0, len(rowPeers)+len(colPeers)+len(blockPeers))

	included := make(map[T]bool)
	for _, item := range slices.Concat(rowPeers, colPeers, blockPeers) {
		if _, isIncluded := included[item]; isIncluded {
			continue
		}

		allPeers = append(allPeers, item)
		included[item] = true
	}

	return Peers[T]{
		row:      slices.Clone(rowPeers),
		column:   slices.Clone(colPeers),
		block:    slices.Clone(blockPeers),
		allPeers: allPeers,
	}
}

// All returns a single iterator over the peers. Duplicates are removed from the
// three peers (row, column, and block).
func (p Peers[T]) All() iter.Seq[T] {
	return slices.Values(p.allPeers)
}

// Each returns an array of iterators. The first provides the peers by row, the
// second provides the peers by column, and the third provides the peers by
// block.
func (p Peers[T]) Each() [3]iter.Seq[T] {
	return [3]iter.Seq[T]{p.Row(), p.Col(), p.Block()}
}

func (p Peers[T]) Row() iter.Seq[T] {
	return slices.Values(p.row)
}

func (p Peers[T]) Col() iter.Seq[T] {
	return slices.Values(p.column)
}

func (p Peers[T]) Block() iter.Seq[T] {
	return slices.Values(p.block)
}
