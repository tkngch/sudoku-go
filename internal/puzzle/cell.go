package puzzle

import "fmt"

type Cell struct {
	position   Position
	candidates Candidates
}

func NewCell(position Position, candidates Candidates) Cell {
	return Cell{position: position, candidates: candidates}
}

func (c Cell) Position() Position {
	return c.position
}

func (c Cell) Candidates() Candidates {
	return c.candidates
}

func (c Cell) Replace(newCandidates Candidates) Cell {
	return Cell{position: c.position, candidates: newCandidates}
}

func (c Cell) String() string {
	return fmt.Sprintf("Cell at %v with %v", c.position, c.candidates)
}
