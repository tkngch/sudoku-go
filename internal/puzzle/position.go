package puzzle

import "fmt"

type Position struct {
	row int
	col int
}

func NewPosition(row, col int) Position {
	return Position{row: row, col: col}
}

func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.row, p.col)
}
