package solver

import (
	"context"
	"errors"
	"fmt"

	"github.com/tkngch/sudoku-go/internal/puzzle"
)

// solver solves a single Sudoku grid. It holds scratch buffers that the hot
// propagation loop reuses across the whole search, so each buffer is grown at
// most once per solve instead of once per step.
//
// Solve creates a fresh solver per call, so a solver is never shared across
// goroutines and needs no synchronization.
type solver struct {
	changed       []puzzle.Cell     // removeInvalidCandidatesFromPeers output
	hiddenSingles []puzzle.Cell     // revealHiddenSingles output
	positions     []puzzle.Position // revealHiddenSingles per-unit scratch
	worklist      []puzzle.Cell     // removeInvalidCandidates BFS queue
}

var (
	// ErrInvalidGrid is returned by Solve when the grid is nil or contains a
	// cell with no candidate values.
	ErrInvalidGrid = errors.New("invalid grid")

	// ErrSolutionNotFound is returned by Solve when the grid has no solution.
	ErrSolutionNotFound = errors.New("solution not found")
)

// Solve returns a solved copy of the grid, or ErrSolutionNotFound if the grid
// has no solution, or ErrInvalidGrid if the grid is nil. The input grid is not
// modified.
//
// ctx must not be nil. If the context is done before a solution is found, Solve
// stops and returns a wrapped ctx.Err() (context.Canceled or
// context.DeadlineExceeded).
//
// When a puzzle admits more than one solution, Solve returns one of them and
// does not detect or report non-uniqueness.
func Solve(ctx context.Context, grid *puzzle.Grid) (*puzzle.Grid, error) {
	if grid == nil {
		return nil, ErrInvalidGrid
	}

	err := ctx.Err()
	if err != nil {
		return nil, fmt.Errorf("solve: %w", err)
	}

	cellCount := grid.CellCount()
	sudokuSolver := &solver{
		changed:       make([]puzzle.Cell, 0, cellCount),
		hiddenSingles: make([]puzzle.Cell, 0, cellCount),
		positions:     make([]puzzle.Position, 0, cellCount),
		worklist:      make([]puzzle.Cell, 0, cellCount),
	}

	return sudokuSolver.solve(ctx, grid)
}

func (s *solver) solve(ctx context.Context, grid *puzzle.Grid) (*puzzle.Grid, error) {
	control := puzzle.NewGridControl(grid.Clone())

	knownCells := make([]puzzle.Cell, 0)

	for cell := range control.Cells() {
		switch cell.Candidates().Count() {
		case 0:
			// A cell with no candidates cannot hold any value, so the grid must
			// be malformed not just unsolvable. Outputs from Parse never reach
			// here, but a Grid built directly via NewGrid can.
			return nil, ErrInvalidGrid

		case 1:
			knownCells = append(knownCells, cell)
		}
	}

	ok := s.removeInvalidCandidates(control, knownCells...)
	if !ok {
		return nil, ErrSolutionNotFound
	}

	return s.searchSolution(ctx, control)
}

// removeInvalidCandidates propagates the values of the revealed cells.
// removeInvalidCandidates returns false when the grid becomes unsolvable:
// either a peer is left with no candidates, or no peer can hold an eliminated
// value.
func (s *solver) removeInvalidCandidates(
	grid *puzzle.GridControl,
	newlyRevealedCells ...puzzle.Cell,
) bool {
	// Reuse the scratch buffer for performance: we don't want to allocate a new
	// slice here.
	s.worklist = append(s.worklist[:0], newlyRevealedCells...)

	// note: len(s.worklist) is re-evaluated every iteration, so if we append
	// to s.worklist within the loop, iteration reaches the newly appended
	// items.
	for idx := 0; idx < len(s.worklist); idx++ {
		revealed := s.worklist[idx]

		s.removeInvalidCandidatesFromPeers(grid, revealed)

		for _, changed := range s.changed {
			switch changed.Candidates().Count() {
			case 0:
				return false
			case 1:
				s.worklist = append(s.worklist, changed)
			}

			ok := s.revealHiddenSingles(grid, changed.Position(), revealed.Candidates())
			if !ok {
				return false
			}

			s.worklist = append(s.worklist, s.hiddenSingles...)
		}
	}

	return true
}

// removeInvalidCandidatesFromPeers removes revealed's value from revealed's
// peers, recording the peers whose candidates changed in s.changed.
func (s *solver) removeInvalidCandidatesFromPeers(grid *puzzle.GridControl, revealed puzzle.Cell) {
	// Reuse the scratch buffer for performance: we don't want to allocate a new
	// slice here.
	s.changed = s.changed[:0]

	peers := grid.AllPeersOf(revealed.Position())
	for idx := range peers.Len() {
		position := peers.At(idx)
		candidates := grid.CandidatesAt(position)

		reduced := candidates.Remove(revealed.Candidates())
		if reduced == candidates {
			continue
		}

		grid.Update(position, candidates, reduced)
		s.changed = append(s.changed, puzzle.NewCell(position, reduced))
	}
}

// After a candidate value is eliminated from the position, this eliminated
// candidate value should be filled in on one of its peers. If there is only one
// cell in the peers that can take the eliminated candidate value, fill that
// cell with it.
func (s *solver) revealHiddenSingles(
	grid *puzzle.GridControl,
	position puzzle.Position,
	eliminatedCandidates puzzle.Candidates,
) bool {
	s.hiddenSingles = s.hiddenSingles[:0]

	if eliminatedCandidates.Count() != 1 {
		// Unreachable in practice: callers only pass the candidate value of
		// revealed cell, which has only one candidate value. This defensive
		// guard is here to highlight the assumption that the eliminated
		// candidates only have one value.
		return true
	}

	for _, peers := range grid.EachPeersOf(position) {
		s.positions = s.positions[:0]

		for idx := range peers.Len() {
			position := peers.At(idx)
			if grid.CandidatesAt(position).Contains(eliminatedCandidates) {
				s.positions = append(s.positions, position)
				if len(s.positions) > 1 {
					break
				}
			}
		}

		switch len(s.positions) {
		case 0:
			// None of the peers can take the eliminated value, so the value
			// should not have been eliminated.
			return false

		case 1:
			// Skip the cell which has only the eliminated value as its candidate values.
			current := grid.CandidatesAt(s.positions[0])
			if current != eliminatedCandidates {
				grid.Update(s.positions[0], current, eliminatedCandidates)
				s.hiddenSingles = append(
					s.hiddenSingles,
					puzzle.NewCell(s.positions[0], eliminatedCandidates),
				)
			}

		default:
		}
	}

	return true
}

func (s *solver) searchSolution(
	ctx context.Context,
	grid *puzzle.GridControl,
) (*puzzle.Grid, error) {
	err := ctx.Err()
	if err != nil {
		return nil, fmt.Errorf("search solution: %w", err)
	}

	cell, isFound := unfilledCellWithFewestCandidates(grid.Grid)
	if !isFound {
		if isSolved(grid.Grid) {
			return grid.Grid, nil
		}

		return nil, ErrSolutionNotFound
	}

	for value := range cell.Candidates().All() {
		grid.Checkpoint()
		grid.Update(cell.Position(), cell.Candidates(), value)

		ok := s.removeInvalidCandidates(grid, puzzle.NewCell(cell.Position(), value))

		if ok {
			solution, err := s.searchSolution(ctx, grid)
			if err == nil {
				grid.CommitUpdates()

				return solution, nil
			}

			if !errors.Is(err, ErrSolutionNotFound) {
				// Any error other than ErrSolutionNotFound is unrecoverable, so
				// propagate it instead of moving on to the next candidate. An
				// example of unrecoverable error is context cancellation: the
				// timeout was reached or the search was interrupted.
				return nil, err
			}
		}

		grid.RestoreCheckpoint()
	}

	return nil, ErrSolutionNotFound
}

// unfilledCellWithFewestCandidates finds the cell that has the smallest number
// of candidates among the cells which have more than one candidates.
func unfilledCellWithFewestCandidates(grid *puzzle.Grid) (puzzle.Cell, bool) {
	isFound := false

	// minBranchingCandidates is the fewest candidates an unfilled cell can
	// have, making such a cell an immediate minimum-remaining-values pick.
	const minBranchingCandidates = 2

	var foundCell puzzle.Cell

	for cell := range grid.Cells() {
		count := cell.Candidates().Count()
		switch count {
		case 0:
			// Unreachable in practice: callers only pass grids that survived
			// removeInvalidCandidates, which rejects any grid with an empty cell.
			return foundCell, false
		case 1:
			continue
		case minBranchingCandidates:
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

		peers := grid.AllPeersOf(cell.Position())
		for idx := range peers.Len() {
			if cell.Candidates() == grid.CandidatesAt(peers.At(idx)) {
				return false
			}
		}
	}

	return true
}
