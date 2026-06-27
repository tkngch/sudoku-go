package puzzle

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Parse(input string) (Grid, error) {
	cellCount := len(input)

	gridSize, err := NewGridSizeFromCellCount(cellCount)
	if err != nil {
		return Grid{}, fmt.Errorf("parse: %w", err)
	}
	regionSize, err := NewRegionSizeFromCellCount(cellCount)
	if err != nil {
		return Grid{}, fmt.Errorf("parse: %w", err)
	}

	minCellValue, maxCellValue := uint8(1), gridSize.rowCount
	cells := make([]Cell, cellCount)
	for i := range cellCount {
		char := input[i]
		position := NewPosition(uint8(i)/gridSize.ColCount(), uint8(i)%gridSize.ColCount())

		var candidates Candidates
		value := toUint8(char)
		if value >= minCellValue && value <= maxCellValue {
			candidates = NewCandidate(value)
		} else {
			candidates = NewCandidates(maxCellValue)
		}

		cells[i] = NewCell(position, candidates)

	}
	return NewGrid(cells, gridSize, regionSize), nil
}

// Compact representation of grid.
func (g Grid) String() string {
	chars := make([]string, len(g.Cells()))
	for i, cell := range g.Cells() {
		chars[i] = toString(cell.Candidates())
	}
	return strings.Join(chars, "")
}

// Multiline, pretty printing of Grid.
func (g Grid) Render() string {
	rowsAsString := make([]string, 0, g.gridSize.RowCount())

	row := make([]string, 0, g.gridSize.colCount*3)
	for _, cell := range g.cells {
		if cell.position.col == 0 && cell.position.row%g.regionSize.rowCount == 0 {
			rowsAsString = append(rowsAsString, g.rowSeparator())
		}
		if cell.position.col%g.regionSize.colCount == 0 {
			row = append(row, "|")
		}
		row = append(row, toString(cell.Candidates()))

		if cell.position.col == g.gridSize.colCount-1 {
			row = append(row, "|")
			rowsAsString = append(rowsAsString, strings.Join(row, " "))
			row = row[:0]
		}
	}

	rowsAsString = append(rowsAsString, g.rowSeparator())
	rowsAsString = append(rowsAsString, fmt.Sprintf("Grid Size: %v", g.gridSize))
	rowsAsString = append(rowsAsString, fmt.Sprintf("Region Size: %v", g.regionSize))

	return strings.Join(rowsAsString, "\n")
}

func (g Grid) rowSeparator() string {
	regionCount := g.gridSize.ColCount() / g.regionSize.ColCount()
	separators := make([]string, regionCount)
	for i := range regionCount {
		separators[i] = strings.Repeat("-", int(g.regionSize.ColCount())*2+1)
	}
	return "+" + strings.Join(separators, "+") + "+"
}

func toUint8(char byte) uint8 {
	switch {
	case '0' <= char && char <= '9':
		return uint8(char - '0')

	case 'a' <= char && char <= 'f':
		return uint8(char-'a') + 10

	case 'A' <= char && char <= 'F':
		return uint8(char-'A') + 10
	}
	return math.MaxUint8
}

func toString(x Candidates) string {
	if x.Count() != 1 {
		return "."
	}

	asInt, err := strconv.Atoi(x.String())
	if err != nil {
		return "."
	}
	return string(fromInt(asInt))
}

func fromInt(x int) byte {
	if 1 <= x && x <= 9 {
		return byte('0' + x)
	}
	if 10 <= x && x <= 16 {
		return byte('a' + x - 10)
	}
	return '.'
}
