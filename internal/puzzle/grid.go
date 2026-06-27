package puzzle

import (
	"iter"
	"slices"
)

// Grid is a square sudoku grid of cells in row-major order
type Grid struct {
	cells      []Cell
	peers      Peers
	gridSize   GridSize
	regionSize RegionSize
}

func NewGrid(cells []Cell, gridSize GridSize, regionSize RegionSize) Grid {
	peers := NewPeers(gridSize, regionSize)
	return Grid{cells: cells, peers: peers, gridSize: gridSize, regionSize: regionSize}
}

func (g Grid) GetPeers(position Position) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		for _, peerPosition := range g.peers.Of(position) {
			peer := g.cells[g.getRowMajorIndex(peerPosition)]
			if !yield(peer) {
				return
			}
		}
	}
}

// Return the cells in row-major order.
func (g Grid) Cells() []Cell {
	return g.cells[:]
}

func (g Grid) Clone() Grid {
	g.cells = slices.Clone(g.cells)
	return g
}

// Copy the grid and alter the cell
func (g Grid) With(position Position, newCandidates Candidates) Grid {
	newGrid := g.Clone()
	newGrid.Set(position, newCandidates)
	return newGrid
}

// Alter the cell in place
func (g Grid) Set(position Position, newCandidates Candidates) {
	index := g.getRowMajorIndex(position)
	g.cells[index] = g.cells[index].Replace(newCandidates)
}

func (g Grid) getRowMajorIndex(position Position) int {
	return int(position.Row())*int(g.gridSize.ColCount()) + int(position.Col())
}
