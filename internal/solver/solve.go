package solver

import (
	"container/list"
	"errors"

	"github.com/tkngch/sudoku-go/internal/puzzle"
)

type sudokuSolver struct {
	grid               *puzzle.Grid
	newlyRevealedCells *list.List
}

var ErrSolutionNotFound = errors.New("solution not found")

// Solve returns a solved copy of the grid, or ErrSolutionNotFound if the grid
// has no solution. The input grid is not modified.
func Solve(grid *puzzle.Grid) (*puzzle.Grid, error) {
	solver := sudokuSolver{
		grid:               grid.Clone(),
		newlyRevealedCells: list.New(),
	}

	for cell := range solver.grid.Cells() {
		if cell.Candidates().Count() == 1 {
			solver.newlyRevealedCells.PushBack(cell)
		}
	}

	ok := solver.eliminateInvalidCandidates()
	if !ok {
		return nil, ErrSolutionNotFound
	}

	if solver.isSolved() {
		return solver.grid, nil
	}

	err := solver.searchSolution()
	if err != nil {
		return nil, ErrSolutionNotFound
	}

	return solver.grid, nil
}

// eliminateInvalidCandidates propagates the values of the revealed cells to their
// peers. It returns false when a peer left with no candidates.
func (s sudokuSolver) eliminateInvalidCandidates() bool {
	for revealed := s.newlyRevealedCells.Front(); revealed != nil; revealed = revealed.Next() {
		cell, isCell := revealed.Value.(puzzle.Cell)
		if !isCell {
			continue
		}

		for peer := range s.grid.PeersOf(cell.Position()) {
			reduced := peer.Candidates().Remove(cell.Candidates())
			if reduced == peer.Candidates() {
				continue
			}

			if reduced.Count() == 0 {
				s.newlyRevealedCells = s.newlyRevealedCells.Init()

				return false // contradiction — prune this branch
			}

			s.grid.Set(peer.Position(), reduced)

			if reduced.Count() == 1 {
				s.newlyRevealedCells.PushBack(puzzle.NewCell(peer.Position(), reduced))
			}
		}
	}

	return true
}

func (s sudokuSolver) isSolved() bool {
	for cell := range s.grid.Cells() {
		if cell.Candidates().Count() != 1 {
			return false
		}

		for peer := range s.grid.PeersOf(cell.Position()) {
			if cell.Candidates() == peer.Candidates() {
				return false
			}
		}
	}

	return true
}

func (s sudokuSolver) searchSolution() error {
	cell, isFound := s.findUnfilledCell()
	if !isFound {
		if s.isSolved() {
			return nil
		}

		return ErrSolutionNotFound
	}

	for value := range cell.Candidates().All() {
		current := s.grid

		s.grid = current.Clone()
		s.grid.Set(cell.Position(), value)

		ok := s.eliminateInvalidCandidates()
		if ok {
			err := s.searchSolution()
			if err == nil {
				return nil
			}
		}

		s.grid = current
	}

	return ErrSolutionNotFound
}

// Find the cell that has the smallest number of candidates among the cells
// which has more than one candidates.
func (s sudokuSolver) findUnfilledCell() (puzzle.Cell, bool) {
	isFound := false

	var foundCell puzzle.Cell

	for cell := range s.grid.Cells() {
		count := cell.Candidates().Count()
		switch count {
		case 0: // search went down the path without solution. fail early.
			return foundCell, false
		case 1:
			continue
		case 2:
			return cell, true
		default:
			if !isFound || count < foundCell.Candidates().Count() {
				isFound = true
				foundCell = cell
			}
		}
	}

	return foundCell, isFound
}
