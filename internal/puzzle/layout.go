package puzzle

import (
	"errors"
	"fmt"
)

type Layout struct {
	blockRowCount, blockColCount int
}

var ErrInvalidCellCount = errors.New("invalid cell count")

func NewLayoutFromCellCount(cellCount int) (Layout, error) {
	switch cellCount {
	case 144:
		return Layout{blockRowCount: 4, blockColCount: 3}, nil
	case 81:
		return Layout{blockRowCount: 3, blockColCount: 3}, nil
	case 36:
		return Layout{blockRowCount: 2, blockColCount: 3}, nil
	case 16:
		return Layout{blockRowCount: 2, blockColCount: 2}, nil

	default:
		err := fmt.Errorf("new layout from cell count [%d]: %w", cellCount, ErrInvalidCellCount)
		return Layout{}, err
	}
}

func (l Layout) BlockColCount() int { return l.blockColCount }

// Return the number of rows or columns in a grid. A grid is a square-shaped, so
// its number of rows equals to its number of columns.
func (l Layout) GridSize() int { return l.blockRowCount * l.blockColCount }

func (l Layout) CellCount() int { return l.GridSize() * l.GridSize() }

func (l Layout) PeerCount() int {
	blockPeerCount := l.blockRowCount*l.blockColCount - 1
	rowPeerCount := l.GridSize() - l.blockRowCount
	colPeerCount := l.GridSize() - l.blockColCount
	return blockPeerCount + rowPeerCount + colPeerCount
}

func (l Layout) AreInSameBlock(this, that Position) bool {
	thisBlock := NewPosition(
		this.Row()/l.blockRowCount,
		this.Col()/l.blockColCount,
	)
	thatBlock := NewPosition(
		that.Row()/l.blockRowCount,
		that.Col()/l.blockColCount,
	)
	return thisBlock == thatBlock
}

func (l Layout) IsOnGrid(position Position) bool {
	return position.Row() >= 0 &&
		position.Row() < l.GridSize() &&
		position.Col() >= 0 &&
		position.Col() < l.GridSize()
}

func (l Layout) IsFirstColumnInBlock(position Position) bool {
	return position.Col() == 0 ||
		!l.AreInSameBlock(position, NewPosition(position.Row(), position.Col()-1))
}

func (l Layout) IsFirstRowInBlock(position Position) bool {
	return position.Row() == 0 ||
		!l.AreInSameBlock(position, NewPosition(position.Row()-1, position.Col()))
}

func (l Layout) RowMajorIndex(position Position) int {
	return int(position.Row())*int(l.GridSize()) + int(position.Col())
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
