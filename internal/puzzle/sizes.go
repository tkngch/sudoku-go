package puzzle

import "fmt"

type size struct {
	rowCount uint8
	colCount uint8
}

func newSize(rowCount, colCount uint8) size {
	return size{rowCount: rowCount, colCount: colCount}
}

func (s size) RowCount() uint8 {
	return s.rowCount
}

func (s size) ColCount() uint8 {
	return s.colCount
}

func (s size) String() string {
	return fmt.Sprintf("%d by %d", s.rowCount, s.colCount)
}

type GridSize struct{ size }

// Grid is expected to be a square, not a rectangle.
// So its row-count is expected to be the same as its column-count.
func NewGridSize(count uint8) GridSize {
	return GridSize{newSize(count, count)}
}

func NewGridSizeFromCellCount(cellCount int) (GridSize, error) {
	switch cellCount {
	case 144:
		return NewGridSize(12), nil
	case 81:
		return NewGridSize(9), nil
	case 36:
		return NewGridSize(6), nil
	case 16:
		return NewGridSize(4), nil

	default:
		return GridSize{}, fmt.Errorf("invalid cell count [%d] for grid size", cellCount)
	}
}

type RegionSize struct{ size }

func NewRegionSize(rowCount, colCount uint8) RegionSize {
	return RegionSize{newSize(rowCount, colCount)}
}

func NewRegionSizeFromCellCount(cellCount int) (RegionSize, error) {
	switch cellCount {
	case 144:
		return NewRegionSize(4, 3), nil
	case 81:
		return NewRegionSize(3, 3), nil
	case 36:
		return NewRegionSize(2, 3), nil
	case 16:
		return NewRegionSize(2, 2), nil

	default:
		return RegionSize{}, fmt.Errorf("invalid cell count [%d] for region size", cellCount)
	}

}
