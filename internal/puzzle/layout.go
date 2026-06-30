package puzzle

import (
	"errors"
	"fmt"
	"iter"
)

// Layout describes the geometry of a Sudoku grid: its block dimensions, the
// derived grid size, and the precomputed peers of every cell. A Layout is
// immutable once constructed, so a value may be freely copied and shared across
// goroutines.
type Layout struct {
	blockRowCount, blockColCount int
	peers                        []Peers
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

// GridSize returns the number of rows or columns in a grid. A grid is
// square-shaped, so its number of rows equals to its number of columns.
func (l Layout) GridSize() int { return l.blockRowCount * l.blockColCount }

func (l Layout) CellCount() int { return l.GridSize() * l.GridSize() }

// PeersOf returns iterators over the precomputed peers. A peer shares the row,
// the column or the group with the provided position. When the provided
// position is off the grid, PeersOf returns empty iterators.
func (l Layout) PeersOf(position Position) Peers {
	if !l.IsOnGrid(position) {
		return NewEmptyPeers()
	}

	return l.peers[l.RowMajorIndex(position)]
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

func (l Layout) allPeers() []Peers {
	rowPeerCount := l.GridSize() - 1
	colPeerCount := l.GridSize() - 1
	blockPeerCount := l.blockRowCount*l.blockColCount - 1

	allPeers := make([]Peers, l.CellCount())

	for this := range l.allPositions() {
		rowPeers := make([]Position, 0, rowPeerCount)
		colPeers := make([]Position, 0, colPeerCount)
		blockPeers := make([]Position, 0, blockPeerCount)

		for that := range l.allPositions() {
			if this == that {
				continue
			}

			if this.row == that.row {
				rowPeers = append(rowPeers, that)
			}

			if this.col == that.col {
				colPeers = append(colPeers, that)
			}

			if l.areInSameBlock(this, that) {
				blockPeers = append(blockPeers, that)
			}
		}

		allPeers[l.RowMajorIndex(this)] = NewPeers(rowPeers, colPeers, blockPeers)
	}

	return allPeers
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

func (l Layout) areInSameBlock(positionA, positionB Position) bool {
	blockA := NewPosition(
		positionA.row/l.blockRowCount,
		positionA.col/l.blockColCount,
	)
	blockB := NewPosition(
		positionB.row/l.blockRowCount,
		positionB.col/l.blockColCount,
	)

	return blockA == blockB
}
