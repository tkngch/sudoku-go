package puzzle

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var ErrInvalidCharacter = errors.New("invalid character")

func Parse(input string) (Grid, error) {
	cellCount := len(input)

	layout, err := NewLayoutFromCellCount(cellCount)
	if err != nil {
		return Grid{}, fmt.Errorf("parse: %w", err)
	}

	minCellValue, maxCellValue := uint8(1), uint8(layout.GridSize())
	cells := make([]Cell, cellCount)
	for i := range cellCount {
		char := input[i]
		position := NewPosition(i/layout.GridSize(), i%layout.GridSize())

		var candidates Candidates
		value, isOk := toUint8(char)
		if isOk && value >= minCellValue && value <= maxCellValue {
			candidates = NewCandidate(value)
		} else if isOk && value == 0 {
			candidates = NewCandidates(maxCellValue)
		} else {
			return Grid{}, fmt.Errorf("parse %q: %w", char, ErrInvalidCharacter)
		}

		cells[i] = NewCell(position, candidates)

	}
	return NewGrid(cells, layout)
}

// Compact representation of grid.
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
	rowsAsString := make([]string, 0, g.layout.GridSize())

	row := make([]string, 0, g.layout.GridSize()*3)
	for cell := range g.Cells() {
		if cell.Position().Col() == 0 && g.layout.IsFirstRowInBlock(cell.Position()) {
			rowsAsString = append(rowsAsString, g.rowSeparator())
		}
		if g.layout.IsFirstColumnInBlock(cell.Position()) {
			row = append(row, "|")
		}
		row = append(row, toValue(cell.Candidates()))

		if cell.Position().Col() == g.layout.GridSize()-1 {
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
	blockCount := g.layout.GridSize() / g.layout.BlockColCount()
	separators := make([]string, blockCount)
	for i := range blockCount {
		separators[i] = strings.Repeat("-", g.layout.BlockColCount()*2+1)
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
