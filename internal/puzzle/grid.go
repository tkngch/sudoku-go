package puzzle

import (
	"errors"
	"fmt"
	"iter"
	"slices"
)

// Grid is a square sudoku grid of cells in row-major order
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

func (g Grid) Clone() Grid {
	g.cells = slices.Clone(g.cells)
	return g
}

// Copy the grid and alter the cell.
func (g Grid) With(position Position, newCandidates Candidates) Grid {
	newGrid := g.Clone()
	newGrid.Set(position, newCandidates)
	return newGrid
}

// Alter the cell in place. Noop if the position is out of range.
func (g *Grid) Set(position Position, newCandidates Candidates) {
	if !g.layout.IsOnGrid(position) {
		return
	}
	index := g.layout.RowMajorIndex(position)
	g.cells[index] = g.cells[index].Replace(newCandidates)
}
