package puzzle

import (
	"errors"
	"fmt"
	"iter"
	"slices"
)

// Grid is a square Sudoku grid of cells in row-major order.
//
// Grid is a mutable reference type and is used through a pointer: copying a
// *Grid aliases the same underlying cells, so Set mutates every alias. Use
// Clone for an independent copy.
type Grid struct {
	cellCandidates []Candidates
	layout         Layout
}

// ErrInvalidCells is returned by NewGrid when the number of cells is
// inconsistent with the layout.
var ErrInvalidCells = errors.New("invalid cells")

// NewGrid returns a Grid holding cells laid out by layout. It returns
// ErrInvalidCells when len(cells) != layout.cellCount().
func NewGrid(cells []Candidates, layout Layout) (*Grid, error) {
	if len(cells) != layout.cellCount() {
		err := fmt.Errorf(
			"expected %d cells, got %d: %w",
			layout.cellCount(), len(cells), ErrInvalidCells,
		)

		return nil, err
	}

	grid := Grid{cellCandidates: slices.Clone(cells), layout: layout}

	return &grid, nil
}

// EachPeersOf returns three separate list of positions over the cells that
// share, respectively, the row, the column, and the block of the provided
// position (in that order), each excluding the cell at position itself.
func (g *Grid) EachPeersOf(position Position) [3]PositionList {
	peers := g.layout.PeersOf(position)

	return [3]PositionList{peers.Row(), peers.Col(), peers.Block()}
}

// AllPeersOf returns a list of positions over the distinct cells that share the
// row, column or block. It excludes the cell at the provided position, because
// a cell cannot be a peer of itself.
func (g *Grid) AllPeersOf(position Position) PositionList {
	return g.layout.PeersOf(position).All()
}

// CandidatesAt returns the candidate values at the position.
func (g *Grid) CandidatesAt(p Position) Candidates {
	return g.cellCandidates[g.layout.RowMajorIndex(p)]
}

// Cells returns an iterator over every cell of the grid in row-major order.
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
