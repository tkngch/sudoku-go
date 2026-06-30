package puzzle

import (
	"iter"
	"slices"
)

type Peers[T comparable] struct {
	Row    iter.Seq[T]
	Column iter.Seq[T]
	Block  iter.Seq[T]
}

func NewEmptyPeers[T comparable]() Peers[T] {
	return Peers[T]{
		Row:    func(yield func(T) bool) {},
		Column: func(yield func(T) bool) {},
		Block:  func(yield func(T) bool) {},
	}
}

func NewPeers[T comparable](rowPeers, colPeers, blockPeers []T) Peers[T] {
	return Peers[T]{
		Row:    slices.Values(rowPeers),
		Column: slices.Values(colPeers),
		Block:  slices.Values(blockPeers),
	}
}

// All returns a single iterator over the peers. Duplicates are removed from the
// three peers (row, column, and block).
func (p Peers[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		yielded := make(map[T]bool)

		for _, seq := range p.Each() {
			for item := range seq {
				if _, isFound := yielded[item]; isFound {
					continue
				}

				if !yield(item) {
					return
				}

				yielded[item] = true
			}
		}
	}
}

// Each returns an array of iterators. The first provides the peers by row, the
// second provides the peers by column, and the third provides the peers by
// block.
func (p Peers[T]) Each() [3]iter.Seq[T] {
	return [3]iter.Seq[T]{p.Row, p.Column, p.Block}
}
