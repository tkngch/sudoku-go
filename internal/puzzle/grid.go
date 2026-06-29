package puzzle

import (
	"errors"
	"fmt"
	"iter"
	"slices"
)

// Grid is a square sudoku grid of cells in row-major order.
//
// Grid's cell-candidates are stored in a slice, so copying a Grid (by
// assignment, function argument, or return value) produces a shallow copy that
// shares the same underlying cells. Mutating such a copy with Set mutates every
// other copy. Use Clone for an independent copy, or With for a copy-on-write
// update.
type Grid struct {
	cellCandidates []Candidates
	layout         Layout
}

var ErrInvalidCells = errors.New("invalid cells")

// NewGrid returns a Grid holding cells laid out by layout. It returns
// ErrInvalidCells when len(cells) != layout.CellCount().
func NewGrid(cells []Candidates, layout Layout) (Grid, error) {
	if len(cells) != layout.CellCount() {
		err := fmt.Errorf(
			"expected %d cells, got %d: %w",
			layout.CellCount(), len(cells), ErrInvalidCells,
		)
		return Grid{}, err
	}

	return Grid{cellCandidates: slices.Clone(cells), layout: layout}, nil
}

// PeersOf returns an iterator over the cells that share a row, column, or block
// with position, excluding the cell at position itself.
func (g Grid) PeersOf(position Position) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		for peerPosition := range g.layout.PeersOf(position) {
			peer := g.cellCandidates[g.layout.RowMajorIndex(peerPosition)]
			if !yield(NewCell(peerPosition, peer)) {
				return
			}
		}
	}
}

func (g Grid) Cells() iter.Seq[Cell] {
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
func (g Grid) Clone() Grid {
	g.cellCandidates = slices.Clone(g.cellCandidates)
	return g
}

// With returns an independent copy of the grid with the cell at position set to
// newCandidates, leaving the original unchanged. It is the copy-on-write
// counterpart to Set. If position is out of range, the returned grid is an
// unmodified copy.
func (g Grid) With(position Position, newCandidates Candidates) Grid {
	newGrid := g.Clone()
	newGrid.Set(position, newCandidates)
	return newGrid
}

// Set alters the cell at position in place; it is a noop if the position is out
// of range. Because copies of a Grid share their cells, Set also mutates every
// shallow copy made since the last Clone. Call Clone first, or use With, to
// keep the original intact.
func (g *Grid) Set(position Position, newCandidates Candidates) {
	if !g.layout.IsOnGrid(position) {
		return
	}
	index := g.layout.RowMajorIndex(position)
	g.cellCandidates[index] = newCandidates
}
