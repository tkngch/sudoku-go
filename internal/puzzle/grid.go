package puzzle

import (
	"fmt"
	"iter"
	"slices"
)

// Grid is a square sudoku grid of cells in row-major order
type Grid struct {
	cells  []Cell
	peers  Peers
	layout Layout
}

func NewGrid(cells []Cell, layout Layout) (Grid, error) {
	if len(cells) != layout.CellCount() {
		err := fmt.Errorf("expected %d cells, got %d", layout.CellCount(), len(cells))
		return Grid{}, err
	}

	peers := NewPeers(layout)
	return Grid{cells: cells, peers: peers, layout: layout}, nil
}

func (g Grid) PeersOf(position Position) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		for peerPosition := range g.peers.Of(position) {
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
