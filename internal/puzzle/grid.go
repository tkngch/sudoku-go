package puzzle

import (
	"errors"
	"fmt"
	"iter"
	"slices"
)

// Grid is a square sudoku grid of cells in row-major order.
//
// Grid is a mutable reference type and is used through a pointer: copying a
// *Grid aliases the same underlying cells, so Set mutates every alias. Use
// Clone for an independent copy.
type Grid struct {
	cellCandidates []Candidates
	layout         Layout
}

var ErrInvalidCells = errors.New("invalid cells")

// NewGrid returns a Grid holding cells laid out by layout. It returns
// ErrInvalidCells when len(cells) != layout.CellCount().
func NewGrid(cells []Candidates, layout Layout) (*Grid, error) {
	if len(cells) != layout.CellCount() {
		err := fmt.Errorf(
			"expected %d cells, got %d: %w",
			layout.CellCount(), len(cells), ErrInvalidCells,
		)

		return nil, err
	}

	grid := Grid{cellCandidates: slices.Clone(cells), layout: layout}

	return &grid, nil
}

// EachPeersOf returns an array of iterators over the cells that share a row,
// column, or block with the provided position.
func (g *Grid) EachPeersOf(position Position) [3]iter.Seq[Cell] {
	peers := g.layout.PeersOf(position)

	return [3]iter.Seq[Cell]{
		g.cellsAt(peers.Row()),
		g.cellsAt(peers.Col()),
		g.cellsAt(peers.Block()),
	}
}

func (g *Grid) AllPeersOf(position Position) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		for peer := range g.layout.PeersOf(position).All() {
			if !yield(g.cellAt(peer)) {
				return
			}
		}
	}
}

func (g *Grid) Cells() iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		for i, candidates := range g.cellCandidates {
			position := NewPosition(i/g.layout.GridSize(), i%g.layout.GridSize())
			if !yield(NewCell(position, candidates)) {
				return
			}
		}
	}
}

// Clone returns an independent copy of the grid whose cells can be mutated with
// Set without affecting the original.
func (g *Grid) Clone() *Grid {
	return &Grid{
		cellCandidates: slices.Clone(g.cellCandidates),
		layout:         g.layout,
	}
}

// Set alters the cell at position in place; it is a noop if the position is out
// of range. Because copies of a *Grid alias the same cells, Set also mutates
// every alias. Call Clone first to keep the original intact.
func (g *Grid) Set(position Position, newCandidates Candidates) {
	if !g.layout.IsOnGrid(position) {
		return
	}

	index := g.layout.RowMajorIndex(position)
	g.cellCandidates[index] = newCandidates
}

func (g *Grid) cellAt(position Position) Cell {
	candidates := g.cellCandidates[g.layout.RowMajorIndex(position)]

	return NewCell(position, candidates)
}

func (g *Grid) cellsAt(positions iter.Seq[Position]) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		for position := range positions {
			if !yield(g.cellAt(position)) {
				return
			}
		}
	}
}
