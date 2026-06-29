package solver

import (
	"errors"

	"github.com/tkngch/sudoku-go/internal/puzzle"
)

var ErrSolutionNotFound = errors.New("solution not found")

// Solve returns a solved copy of the grid, or ErrSolutionNotFound if the grid
// has no solution. The input grid is not modified.
func Solve(grid *puzzle.Grid) (*puzzle.Grid, error) {
	grid = grid.Clone() // do not mutate the input argument
	dropAllInvalidCandidates(grid)

	if isSolved(grid) {
		return grid, nil
	}

	return searchSolution(grid)
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

// Mutate the grid to drop invalid candidates for all cells.
func dropAllInvalidCandidates(grid *puzzle.Grid) {
	for {
		isAnyValueDropped := false

		for cell := range grid.Cells() {
			isDropped := dropInvalidCandidates(grid, cell)
			isAnyValueDropped = isAnyValueDropped || isDropped
		}

		if !isAnyValueDropped {
			return
		}
	}
}

// Mutate the grid to drop invalid candidates for the cell.
func dropInvalidCandidates(grid *puzzle.Grid, cell puzzle.Cell) bool {
	invalidCandidates := puzzle.NewSingleCandidate(0)

	for peer := range grid.PeersOf(cell.Position()) {
		if peer.Candidates().Count() == 1 {
			invalidCandidates = invalidCandidates.Union(peer.Candidates())
		}
	}

	newCandidates := cell.Candidates().Remove(invalidCandidates)
	if cell.Candidates() != newCandidates {
		grid.Set(cell.Position(), newCandidates)

		return true
	}

	return false
}

func searchSolution(grid *puzzle.Grid) (*puzzle.Grid, error) {
	// 1. Find the cell that has more than one candidates.
	cell, isFound := findUnfilledCell(grid)
	if !isFound && isSolved(grid) {
		return grid, nil
	} else if !isFound {
		return grid, ErrSolutionNotFound
	}

	// 2. Pick one of its candidates and replace the cell
	for value := range cell.Candidates().All() {
		newGrid := grid.Clone()
		newGrid.Set(cell.Position(), value)

		// 3. Update the peers and the others
		dropAllInvalidCandidates(newGrid)

		solved, err := searchSolution(newGrid)
		if err == nil {
			return solved, nil
		} else if !errors.Is(err, ErrSolutionNotFound) {
			return solved, err
		}
	}

	return grid, ErrSolutionNotFound
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
