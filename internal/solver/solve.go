package solver

import (
	"errors"

	"github.com/tkngch/sudoku-go/internal/puzzle"
)

var ErrSolutionNotFound = errors.New("solution not found")

// Solve returns a solved copy of the grid, or ErrSolutionNotFound if the grid
// has no solution. The input grid is not modified.
func Solve(grid *puzzle.Grid) (*puzzle.Grid, error) {
	grid = grid.Clone()

	knownCells := make([]puzzle.Cell, 0)

	for cell := range grid.Cells() {
		if cell.Candidates().Count() == 1 {
			knownCells = append(knownCells, cell)
		}
	}

	ok := eliminateInvalidCandidates(grid, knownCells)
	if !ok {
		return nil, ErrSolutionNotFound
	}

	if isSolved(grid) {
		return grid, nil
	}

	solution, err := searchSolution(grid)
	if err != nil {
		return nil, ErrSolutionNotFound
	}

	return solution, nil
}

// eliminateInvalidCandidates propagates the values of the revealed cells to their
// peers. It returns false when a peer left with no candidates.
func eliminateInvalidCandidates(grid *puzzle.Grid, newlyRevealedCells []puzzle.Cell) bool {
	for len(newlyRevealedCells) > 0 {
		revealed := newlyRevealedCells[0]
		newlyRevealedCells = newlyRevealedCells[1:]

		for peer := range grid.PeersOf(revealed.Position()) {
			reduced := peer.Candidates().Remove(revealed.Candidates())
			if reduced == peer.Candidates() {
				continue
			}

			if reduced.Count() == 0 {
				return false // contradiction — prune this branch
			}

			grid.Set(peer.Position(), reduced)

			if reduced.Count() == 1 {
				newlyRevealedCells = append(
					newlyRevealedCells,
					puzzle.NewCell(peer.Position(), reduced),
				)
			}
		}
	}

	return true
}

func isSolved(grid *puzzle.Grid) bool {
	for cell := range grid.Cells() {
		if cell.Candidates().Count() != 1 {
			return false
		}

		for peer := range grid.PeersOf(cell.Position()) {
			if cell.Candidates() == peer.Candidates() {
				return false
			}
		}
	}

	return true
}

func searchSolution(grid *puzzle.Grid) (*puzzle.Grid, error) {
	cell, isFound := findUnfilledCell(grid)
	if !isFound {
		if isSolved(grid) {
			return grid, nil
		}

		return nil, ErrSolutionNotFound
	}

	for value := range cell.Candidates().All() {
		newGrid := grid.Clone()
		newGrid.Set(cell.Position(), value)

		ok := eliminateInvalidCandidates(
			newGrid,
			[]puzzle.Cell{puzzle.NewCell(cell.Position(), value)},
		)
		if ok {
			solution, err := searchSolution(newGrid)
			if err == nil {
				return solution, nil
			}
		}
	}

	return nil, ErrSolutionNotFound
}

// Find the cell that has the smallest number of candidates among the cells
// which has more than one candidates.
func findUnfilledCell(grid *puzzle.Grid) (puzzle.Cell, bool) {
	isFound := false

	var foundCell puzzle.Cell

	for cell := range grid.Cells() {
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
