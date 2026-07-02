package puzzle

import "fmt"

// Cell holds a position within a grid and the candidate values it can take.
type Cell struct {
	position   Position
	candidates Candidates
}

// NewCell returns an instance of Cell with the provided values.
func NewCell(position Position, candidates Candidates) Cell {
	return Cell{position: position, candidates: candidates}
}

// Position returns the cell's position.
func (c Cell) Position() Position {
	return c.position
}

// Candidates returns the candidate values that the cell may take.
func (c Cell) Candidates() Candidates {
	return c.candidates
}

func (c Cell) String() string {
	return fmt.Sprintf("Cell at %v with %v", c.position, c.candidates)
}
