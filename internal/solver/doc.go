// Package solver finds solutions to variable-size Sudoku puzzles.
//
// Solve first propagates constraints by identifying naked singles and hidden
// singles. To illustrate, suppose a cell is filled with a value and the value
// is eliminated from its peers' candidates.
//
// Naked singles are any cells that are left with only one candidate. Those
// cells are fill with their only candidates. Hidden singles, on the other hand,
// are cells that hold a candidate value that none of their peers hold. Those
// cells are filled with their unique candidates.
//
// Then, Solve falls back to depth-first backtracking guided by the
// minimum-remaining-values search.
package solver
