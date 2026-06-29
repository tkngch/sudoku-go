package puzzle

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var ErrInvalidCharacter = errors.New("invalid character")

// Parse reads a puzzle from its compact form: one character per cell in
// row-major order, whose length selects the layout (see
// NewLayoutFromCellCount). '0' or '.' is an empty cell (all candidates);
// '1'-'9' and 'a'-'g'/'A'-'G' (values 10-16) are givens. It returns
// ErrInvalidCellCount or ErrInvalidCharacter for malformed input.
func Parse(input string) (Grid, error) {
	cellCount := len(input)

	layout, err := NewLayoutFromCellCount(cellCount)
	if err != nil {
		return Grid{}, fmt.Errorf("parse: %w", err)
	}

	minCellValue, maxCellValue := uint8(1), uint8(layout.GridSize())
	cells := make([]Candidates, cellCount)
	for i := range cellCount {
		char := input[i]

		value, ok := toUint8(char)
		if ok && value >= minCellValue && value <= maxCellValue {
			cells[i] = NewSingleCandidate(value)
		} else if ok && value == 0 {
			cells[i] = NewCandidatesForRange(maxCellValue)
		} else {
			return Grid{}, fmt.Errorf("parse %q: %w", char, ErrInvalidCharacter)
		}
	}
	return NewGrid(cells, layout)
}

// Return the compact, single-line form of the grid: one character per cell in
// row-major order, inverse to Parse for solved or given cells.
//
// It is lossy: a cell with more than one candidate is written as '.', so String
// preserves only cells with single candidates and is not a serialization of
// unsolved puzzle.
func (g Grid) String() string {
	cells := slices.Collect(g.Cells())

	var b strings.Builder
	b.Grow(len(cells))
	for _, cell := range cells {
		b.WriteString(toValue(cell.Candidates()))
	}
	return b.String()
}

// Multiline, pretty printing of Grid.
func (g Grid) Render() string {
	if len(g.cellCandidates) == 0 {
		return ""
	}

	rowsAsString := make([]string, 0, g.layout.GridSize())

	row := make([]string, 0, g.layout.GridSize()*3)
	for cell := range g.Cells() {
		if cell.Position().col == 0 && g.layout.IsFirstRowInBlock(cell.Position()) {
			rowsAsString = append(rowsAsString, g.rowSeparator())
		}
		if g.layout.IsFirstColumnInBlock(cell.Position()) {
			row = append(row, "|")
		}
		row = append(row, toValue(cell.Candidates()))

		if cell.Position().col == g.layout.GridSize()-1 {
			row = append(row, "|")
			rowsAsString = append(rowsAsString, strings.Join(row, " "))
			row = row[:0]
		}
	}

	rowsAsString = append(rowsAsString, g.rowSeparator())
	rowsAsString = append(rowsAsString, g.layout.String())

	return strings.Join(rowsAsString, "\n")
}

func (g Grid) rowSeparator() string {
	blockCount := g.layout.GridSize() / g.layout.blockColCount
	separators := make([]string, blockCount)
	for i := range blockCount {
		separators[i] = strings.Repeat("-", g.layout.blockColCount*2+1)
	}
	return "+" + strings.Join(separators, "+") + "+"
}

func toUint8(char byte) (uint8, bool) {
	switch {
	case char == '0' || char == '.':
		return 0, true

	case '0' <= char && char <= '9':
		return uint8(char - '0'), true

	case 'a' <= char && char <= 'g':
		return uint8(char-'a') + 10, true

	case 'A' <= char && char <= 'G':
		return uint8(char-'A') + 10, true
	}
	return 0, false
}

func toValue(x Candidates) string {
	if x.Count() != 1 {
		return "."
	}
	return x.String()
}
