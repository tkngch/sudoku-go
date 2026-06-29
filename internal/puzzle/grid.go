package puzzle

import (
	"errors"
	"fmt"
	"iter"
	"slices"
)

// Grid is a square sudoku grid of cells in row-major order.
//
// Grid's cells are stored in a slice, so copying a Grid (by assignment,
// function argument, or return value) produces a shallow copy that shares the
// same underlying cells. Mutating such a copy with Set mutates every other
// copy. Use Clone for an independent copy, or With for a copy-on-write update.
type Grid struct {
	cells  []Cell
	layout Layout
}

var ErrInvalidCells = errors.New("invalid cells")

func NewGrid(cells []Cell, layout Layout) (Grid, error) {
	if len(cells) != layout.CellCount() {
		err := fmt.Errorf(
			"expected %d cells, got %d: %w",
			layout.CellCount(), len(cells), ErrInvalidCells,
		)
		return Grid{}, err
	}

	for i, cell := range cells {
		expected := NewPosition(i/layout.GridSize(), i%layout.GridSize())
		if cell.Position() != expected {
			err := fmt.Errorf(
				"cells are not row-major ordered (cell %d has position %v when %v is expected): %w",
				i, cell.Position(), expected, ErrInvalidCells,
			)
			return Grid{}, err
		}
	}

	return Grid{cells: cells, layout: layout}, nil
}

func (g Grid) PeersOf(position Position) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		for peerPosition := range g.layout.PeersOf(position) {
			peer := g.cells[g.layout.RowMajorIndex(peerPosition)]
			if !yield(peer) {
				return
			}
		}
	}
}

// Iterate over the cells in row-major order.
func (g Grid) Cells() iter.Seq[Cell] {
	return slices.Values(g.cells)
}

// Clone returns an independent copy of the grid whose cells can be mutated with
// Set without affecting the original.
func (g Grid) Clone() Grid {
	g.cells = slices.Clone(g.cells)
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
	g.cells[index] = g.cells[index].Replace(newCandidates)
}
