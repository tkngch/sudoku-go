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

	ok := removeInvalidCandidates(grid, knownCells)
	if !ok {
		return nil, ErrSolutionNotFound
	}

	return searchSolution(grid)
}

// removeInvalidCandidates propagates the values of the revealed cells to their
// peers. It returns false when a peer left with no candidates.
func removeInvalidCandidates(grid *puzzle.Grid, newlyRevealedCells []puzzle.Cell) bool {
	for len(newlyRevealedCells) > 0 {
		revealed := newlyRevealedCells[0]
		newlyRevealedCells = newlyRevealedCells[1:]

		changedCells := removeInvalidCandidatesFromPeers(grid, revealed)
		for _, cell := range changedCells {
			switch cell.Candidates().Count() {
			case 0:
				return false
			case 1:
				newlyRevealedCells = append(newlyRevealedCells, cell)
			}

			hiddenSingles, ok := revealHiddenSingles(grid, cell.Position(), revealed.Candidates())
			if !ok {
				return false
			}

			newlyRevealedCells = append(newlyRevealedCells, hiddenSingles...)
		}
	}

	return true
}

// removeInvalidCandidatesFromPeers removes revealed's value from revealed's
// peers. It returns the peers whose candidate values have changed. ended up
// with a single candidate (naked singles) and the peers whose candidates
// changed.
func removeInvalidCandidatesFromPeers(grid *puzzle.Grid, revealed puzzle.Cell) []puzzle.Cell {
	changed := make([]puzzle.Cell, 0)

	for peer := range grid.PeersOf(revealed.Position()).All() {
		reduced := peer.Candidates().Remove(revealed.Candidates())
		if reduced == peer.Candidates() {
			continue
		}

		grid.Set(peer.Position(), reduced)
		changed = append(changed, puzzle.NewCell(peer.Position(), reduced))
	}

	return changed
}

// After a candidate value is eliminated from the position, this eliminated
// candidate value should be filled in on one of its peers. If there is only one
// cell in the peers that can take the eliminated candidate value, fill that
// cell with it.
func revealHiddenSingles(
	grid *puzzle.Grid,
	position puzzle.Position,
	eliminatedCandidates puzzle.Candidates,
) ([]puzzle.Cell, bool) {
	hiddenSingles := make([]puzzle.Cell, 0)
	if eliminatedCandidates.Count() != 1 {
		return hiddenSingles, true
	}

	cellsWithEliminatedCandidates := make([]puzzle.Cell, 0)
	for _, peers := range grid.PeersOf(position).Each() {
		cellsWithEliminatedCandidates = cellsWithEliminatedCandidates[:0]

		for peer := range peers {
			if peer.Candidates().Contains(eliminatedCandidates) {
				cellsWithEliminatedCandidates = append(cellsWithEliminatedCandidates, peer)
			}
		}

		switch len(cellsWithEliminatedCandidates) {
		case 0:
			// None of the peers can take the eliminate value, so the value
			// should not have been eliminated.
			return nil, false

		case 1:
			grid.Set(cellsWithEliminatedCandidates[0].Position(), eliminatedCandidates)
			hiddenSingles = append(
				hiddenSingles,
				puzzle.NewCell(cellsWithEliminatedCandidates[0].Position(), eliminatedCandidates),
			)

		default:
		}
	}

	return hiddenSingles, true
}

func searchSolution(grid *puzzle.Grid) (*puzzle.Grid, error) {
	cell, isFound := unfilledCellWithFewestCandidates(grid)
	if !isFound {
		if isSolved(grid) {
			return grid, nil
		}

		return nil, ErrSolutionNotFound
	}

	for value := range cell.Candidates().All() {
		newGrid := grid.Clone()
		newGrid.Set(cell.Position(), value)

		ok := removeInvalidCandidates(
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

// unfilledCellWithFewestCandidates finds the cell that has the smallest number
// of candidates among the cells which has more than one candidates.
func unfilledCellWithFewestCandidates(grid *puzzle.Grid) (puzzle.Cell, bool) {
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

func isSolved(grid *puzzle.Grid) bool {
	for cell := range grid.Cells() {
		if cell.Candidates().Count() != 1 {
			return false
		}

		for peer := range grid.PeersOf(cell.Position()).All() {
			if cell.Candidates() == peer.Candidates() {
				return false
			}
		}
	}

	return true
}
