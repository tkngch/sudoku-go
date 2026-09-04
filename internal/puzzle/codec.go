package puzzle

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	// ErrInvalidCharacter is returned by Parse when the input contains a
	// character that no cell accepts.
	ErrInvalidCharacter = errors.New("invalid character")

	// ErrValueOutOfRange is returned by Parse when the input contains a valid
	// digit or letter whose value is larger than the grid size, such as '9' in
	// a 4x4 puzzle.
	ErrValueOutOfRange = errors.New("value out of range")
)

// Parse reads a puzzle written as one character per cell, in row-major order.
// Parse ignores whitespace, so a puzzle may span one line or several lines, for
// example as a pasted grid. The number of characters that remain selects the
// layout (see NewLayoutForCellCount).
//
// '0' or '.' is an empty cell (all candidates); '1'-'9' and 'a'-'g'/'A'-'G'
// (values 10-16) are givens. Parse returns ErrInvalidCellCount,
// ErrInvalidCharacter, or ErrValueOutOfRange for malformed input.
//
// An invalid character takes precedence over an invalid cell count. A stray
// multi-byte character, such as a BOM that a paste carries, therefore reports
// ErrInvalidCharacter and not a misleading ErrInvalidCellCount.
func Parse(input string) (*Grid, error) {
	compact := strings.Join(strings.Fields(input), "")

	values := make([]int, 0, len(compact))
	for _, char := range compact {
		value, ok := toInt(char)
		if !ok {
			return nil, fmt.Errorf("parse %q: %w", char, ErrInvalidCharacter)
		}

		values = append(values, value)
	}

	layout, err := NewLayoutForCellCount(len(values))
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	maxCellValue := layout.GridSize()
	cells := make([]Candidates, 0, len(compact))

	for _, value := range values {
		switch {
		case value == 0:
			cells = append(cells, NewCandidatesForRange(maxCellValue))
		case value <= maxCellValue:
			cells = append(cells, NewSingleCandidate(value))
		default:
			return nil, fmt.Errorf("parse %d: %w", value, ErrValueOutOfRange)
		}
	}

	return NewGrid(cells, layout)
}

// String returns the compact, single-line form of the grid: one character per
// cell in row-major order, inverse to Parse for solved or given cells.
//
// It is lossy: a cell with more than one candidate is written as '.', so String
// preserves only cells with single candidates and is not a serialization of
// unsolved puzzle.
func (g *Grid) String() string {
	if g == nil {
		return ""
	}

	cells := slices.Collect(g.Cells())

	var builder strings.Builder

	builder.Grow(len(cells))

	for _, cell := range cells {
		builder.WriteString(toValue(cell.Candidates()))
	}

	return builder.String()
}

// Render returns a multiline, pretty-printed rendering of the grid.
func (g *Grid) Render() string {
	if g == nil || len(g.cellCandidates) == 0 {
		return ""
	}

	// maxRowElementsPerCell over-estimates the strings a rendered cell adds to
	// a row (its value plus block separators), for slice pre-allocation.
	const maxRowElementsPerCell = 3

	rowsAsString := make([]string, 0, g.layout.GridSize())

	row := make([]string, 0, g.layout.GridSize()*maxRowElementsPerCell)
	for cell := range g.Cells() {
		if cell.Position().col == 0 && g.layout.isFirstRowInBlock(cell.Position()) {
			rowsAsString = append(rowsAsString, g.rowSeparator())
		}

		if g.layout.isFirstColumnInBlock(cell.Position()) {
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

	return strings.Join(rowsAsString, "\n")
}

func (g *Grid) rowSeparator() string {
	const dashesPerColumn = 2

	blockCount := g.layout.GridSize() / g.layout.blockColCount

	separators := make([]string, blockCount)
	for i := range blockCount {
		separators[i] = strings.Repeat("-", g.layout.blockColCount*dashesPerColumn+1)
	}

	return "+" + strings.Join(separators, "+") + "+"
}

func toInt(char rune) (int, bool) {
	switch {
	case char == '0' || char == '.':
		return 0, true

	case '0' <= char && char <= '9':
		return int(char - '0'), true

	case 'a' <= char && char <= 'g':
		return int(char - 'a' + firstLetterValue), true

	case 'A' <= char && char <= 'G':
		return int(char - 'A' + firstLetterValue), true
	}

	return 0, false
}

func toValue(x Candidates) string {
	if x.Count() != 1 {
		return "."
	}

	return x.String()
}
