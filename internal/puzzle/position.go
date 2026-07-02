package puzzle

import "fmt"

// Position is a zero-based (row, column) coordinate on the grid.
type Position struct {
	row int
	col int
}

// NewPosition returns the Position at row and col.
func NewPosition(row, col int) Position {
	return Position{row: row, col: col}
}

func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.row, p.col)
}
