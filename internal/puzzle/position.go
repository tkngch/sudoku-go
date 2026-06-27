package puzzle

import "fmt"

type Position struct {
	row uint8
	col uint8
}

func NewPosition(row, col uint8) Position {
	return Position{row: row, col: col}
}

func (p Position) Row() uint8 { return p.row }

func (p Position) Col() uint8 { return p.col }

func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.row, p.col)
}
