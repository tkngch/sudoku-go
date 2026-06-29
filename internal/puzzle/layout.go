package puzzle

import (
	"errors"
	"fmt"
	"iter"
	"slices"
)

// Layout describes the geometry of a Sudoku grid: its block dimensions, the
// derived grid size, and the precomputed peers of every cell. A Layout is
// immutable once constructed, so a value may be freely copied and shared across
// goroutines.
type Layout struct {
	blockRowCount, blockColCount int
	peers                        [][]Position
}

var ErrInvalidCellCount = errors.New("invalid cell count")

func NewLayoutFromCellCount(cellCount int) (Layout, error) {
	switch cellCount {
	case 144:
		return newLayout(4, 3), nil
	case 81:
		return newLayout(3, 3), nil
	case 36:
		return newLayout(2, 3), nil
	case 16:
		return newLayout(2, 2), nil

	default:
		err := fmt.Errorf("new layout from cell count [%d]: %w", cellCount, ErrInvalidCellCount)
		return Layout{}, err
	}
}

func newLayout(r, c int) Layout {
	l := Layout{blockRowCount: r, blockColCount: c}
	l.peers = l.allPeers()
	return l
}

// Return the number of rows or columns in a grid. A grid is a square-shaped, so
// its number of rows equals to its number of columns.
func (l Layout) GridSize() int { return l.blockRowCount * l.blockColCount }

func (l Layout) CellCount() int { return l.GridSize() * l.GridSize() }

// Iterate over the positions that share a row, column, or block with the given
// position, excluding itself. It yields nothing when the position is off the
// grid.
func (l Layout) PeersOf(position Position) iter.Seq[Position] {
	if !l.IsOnGrid(position) {
		return func(yield func(Position) bool) {}
	}
	// return an iterator, to prevent the caller from altering the peers.
	return slices.Values(l.peers[l.RowMajorIndex(position)])
}

func (l Layout) allPeers() [][]Position {
	blockPeerCount := l.blockRowCount*l.blockColCount - 1
	rowPeerCount := l.GridSize() - l.blockColCount
	colPeerCount := l.GridSize() - l.blockRowCount
	peerCount := blockPeerCount + rowPeerCount + colPeerCount

	peers := make([][]Position, l.CellCount())
	for this := range l.allPositions() {
		thisPeers := make([]Position, 0, peerCount)

		for that := range l.allPositions() {
			if l.arePeers(this, that) {
				thisPeers = append(thisPeers, that)
			}
		}

		peers[l.RowMajorIndex(this)] = thisPeers
	}
	return peers
}

func (l Layout) allPositions() iter.Seq[Position] {
	return func(yield func(Position) bool) {
		for row := range l.GridSize() {
			for col := range l.GridSize() {
				if !yield(NewPosition(row, col)) {
					return
				}
			}
		}
	}
}

func (l Layout) arePeers(a, b Position) bool {
	if a == b {
		return false
	}
	if a.row == b.row {
		return true
	}
	if a.col == b.col {
		return true
	}
	return l.areInSameBlock(a, b)
}

func (l Layout) areInSameBlock(a, b Position) bool {
	blockA := NewPosition(
		a.row/l.blockRowCount,
		a.col/l.blockColCount,
	)
	blockB := NewPosition(
		b.row/l.blockRowCount,
		b.col/l.blockColCount,
	)
	return blockA == blockB
}

func (l Layout) IsOnGrid(position Position) bool {
	return position.row >= 0 &&
		position.row < l.GridSize() &&
		position.col >= 0 &&
		position.col < l.GridSize()
}

func (l Layout) IsFirstColumnInBlock(position Position) bool {
	return position.col == 0 ||
		!l.areInSameBlock(position, NewPosition(position.row, position.col-1))
}

func (l Layout) IsFirstRowInBlock(position Position) bool {
	return position.row == 0 ||
		!l.areInSameBlock(position, NewPosition(position.row-1, position.col))
}

func (l Layout) RowMajorIndex(position Position) int {
	return position.row*l.GridSize() + position.col
}

func (l Layout) String() string {
	return fmt.Sprintf(
		"%d-by-%d grid with %d %d-by-%d blocks",
		l.GridSize(),
		l.GridSize(),
		l.blockColCount*l.blockRowCount,
		l.blockRowCount,
		l.blockColCount,
	)
}
