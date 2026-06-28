package puzzle

import (
	"errors"
	"fmt"
	"iter"
	"slices"
)

type Layout struct {
	blockRowCount, blockColCount int
}

var ErrInvalidCellCount = errors.New("invalid cell count")
var cachedPeers = make(map[Layout][][]Position)

func NewLayoutFromCellCount(cellCount int) (Layout, error) {
	switch cellCount {
	case 144:
		return Layout{blockRowCount: 4, blockColCount: 3}, nil
	case 81:
		return Layout{blockRowCount: 3, blockColCount: 3}, nil
	case 36:
		return Layout{blockRowCount: 2, blockColCount: 3}, nil
	case 16:
		return Layout{blockRowCount: 2, blockColCount: 2}, nil

	default:
		err := fmt.Errorf("new layout from cell count [%d]: %w", cellCount, ErrInvalidCellCount)
		return Layout{}, err
	}
}

func (l Layout) BlockColCount() int { return l.blockColCount }

// Return the number of rows or columns in a grid. A grid is a square-shaped, so
// its number of rows equals to its number of columns.
func (l Layout) GridSize() int { return l.blockRowCount * l.blockColCount }

func (l Layout) CellCount() int { return l.GridSize() * l.GridSize() }

func (l Layout) PeerCount() int {
	blockPeerCount := l.blockRowCount*l.blockColCount - 1
	rowPeerCount := l.GridSize() - l.blockRowCount
	colPeerCount := l.GridSize() - l.blockColCount
	return blockPeerCount + rowPeerCount + colPeerCount
}

func (l Layout) PeersOf(position Position) iter.Seq[Position] {
	if !l.IsOnGrid(position) {
		return func(yield func(Position) bool) {}
	}

	peers, isCached := cachedPeers[l]
	if !isCached {
		peers = l.allPeers()
		cachedPeers[l] = peers
	}

	// return an iterator, to prevent the caller from altering the peers.
	return slices.Values(peers[l.RowMajorIndex(position)])
}

func (l Layout) allPeers() [][]Position {
	peers := make([][]Position, l.CellCount())
	peerCount := l.PeerCount()

	for this := range l.allPositions() {
		thisPeers := make([]Position, 0, peerCount)

		for that := range l.allPositions() {
			if l.arePeers(this, that) {
				thisPeers = append(thisPeers, that)
			}
		}

		peers[l.RowMajorIndex(this)] = thisPeers
	}
	return peers
}

func (l Layout) allPositions() iter.Seq[Position] {
	return func(yield func(Position) bool) {
		for row := range l.GridSize() {
			for col := range l.GridSize() {
				if !yield(NewPosition(row, col)) {
					return
				}
			}
		}
	}
}

func (l Layout) arePeers(a, b Position) bool {
	if a == b {
		return false
	}
	if a.Row() == b.Row() {
		return true
	}
	if a.Col() == b.Col() {
		return true
	}
	return l.areInSameBlock(a, b)
}

func (l Layout) areInSameBlock(a, b Position) bool {
	blockA := NewPosition(
		a.Row()/l.blockRowCount,
		a.Col()/l.blockColCount,
	)
	blockB := NewPosition(
		b.Row()/l.blockRowCount,
		b.Col()/l.blockColCount,
	)
	return blockA == blockB
}

func (l Layout) IsOnGrid(position Position) bool {
	return position.Row() >= 0 &&
		position.Row() < l.GridSize() &&
		position.Col() >= 0 &&
		position.Col() < l.GridSize()
}

func (l Layout) IsFirstColumnInBlock(position Position) bool {
	return position.Col() == 0 ||
		!l.areInSameBlock(position, NewPosition(position.Row(), position.Col()-1))
}

func (l Layout) IsFirstRowInBlock(position Position) bool {
	return position.Row() == 0 ||
		!l.areInSameBlock(position, NewPosition(position.Row()-1, position.Col()))
}

func (l Layout) RowMajorIndex(position Position) int {
	return position.Row()*l.GridSize() + position.Col()
}

func (l Layout) String() string {
	return fmt.Sprintf(
		"%d-by-%d grid with %d %d-by-%d blocks",
		l.GridSize(),
		l.GridSize(),
		l.blockColCount*l.blockRowCount,
		l.blockRowCount,
		l.blockColCount,
	)
}
