package puzzle

import "slices"

// GridControl is a grid with cached updates.
type GridControl struct {
	*Grid

	updates []gridUpdate
}

type gridUpdate struct {
	position   Position
	oldValue   Candidates
	newValue   Candidates
	checkpoint bool
}

// NewGridControl returns GridControl that is ready to use.
func NewGridControl(grid *Grid) *GridControl {
	return &GridControl{
		Grid:    grid,
		updates: make([]gridUpdate, 0, grid.CellCount()),
	}
}

// Update sets the grid cell's value to the new candidate. The old candidate
// values are cached until CommitUpdates or RollbackUpdates iscalled.
func (g *GridControl) Update(
	position Position,
	currentCandidates Candidates,
	newCandidates Candidates,
) {
	update := gridUpdate{
		position:   position,
		oldValue:   currentCandidates,
		newValue:   newCandidates,
		checkpoint: false,
	}
	g.Grid.Update(position, newCandidates)
	g.updates = append(g.updates, update)
}

// Checkpoint adds a checkpoint to the last update.
func (g *GridControl) Checkpoint() {
	if len(g.updates) > 0 {
		g.updates[len(g.updates)-1].checkpoint = true
	}
}

// CommitUpdates clears the cached updates, such that the updates cannot be rolled back later.
func (g *GridControl) CommitUpdates() {
	g.updates = g.updates[:0]
}

// RestoreCheckpoint roll backs the updates to the last checkpoint.
func (g *GridControl) RestoreCheckpoint() {
	truncateTo := 0

	for idx, update := range slices.Backward(g.updates) {
		if update.checkpoint {
			g.updates[idx].checkpoint = false
			truncateTo = idx + 1

			break
		}

		g.Grid.Update(update.position, update.oldValue)
	}

	g.updates = g.updates[:truncateTo]
}
