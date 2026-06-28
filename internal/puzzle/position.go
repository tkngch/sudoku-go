package puzzle

import "fmt"

type Position struct {
	row uint
	col uint
}

func NewPosition(row, col uint) Position {
	return Position{row: row, col: col}
}

func (p Position) Row() uint { return p.row }

func (p Position) Col() uint { return p.col }

func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.row, p.col)
}
