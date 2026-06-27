package puzzle

import "fmt"

type Layout struct {
	blockRowCount, blockColCount uint8
}

func NewLayout(blockRowCount, blockColCount uint8) Layout {
	return Layout{blockRowCount: blockRowCount, blockColCount: blockColCount}
}

func NewLayoutFromCellCount(cellCount int) (Layout, error) {
	switch cellCount {
	case 144:
		return NewLayout(4, 3), nil
	case 81:
		return NewLayout(3, 3), nil
	case 36:
		return NewLayout(2, 3), nil
	case 16:
		return NewLayout(2, 2), nil

	default:
		return Layout{}, fmt.Errorf("invalid cell count [%d] for layout", cellCount)
	}
}

func (l Layout) BlockRowCount() uint8 { return l.blockRowCount }

func (l Layout) BlockColCount() uint8 { return l.blockColCount }

// Return the number of rows or columns in a grid. A grid is a square-shaped, so
// its number of rows equals to its number of columns.
func (l Layout) GridSize() uint8 { return l.blockRowCount * l.blockColCount }

func (l Layout) String() string {
	return fmt.Sprintf(
		"%d-by-%d grid with %d %d-by-%d blocks",
		l.GridSize(),
		l.GridSize(),
		l.BlockColCount()*l.BlockRowCount(),
		l.BlockRowCount(),
		l.BlockColCount(),
	)
}
