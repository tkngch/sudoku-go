package puzzle_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkngch/sudoku-go/internal/puzzle"
)

func TestNewSingleCandidate(t *testing.T) {
	testCases := []struct {
		input    uint8
		expected puzzle.Candidates
	}{
		{input: 0, expected: 0b0},
		{input: 1, expected: 0b1},
		{input: 2, expected: 0b10},
		{input: 3, expected: 0b100},
		{input: 4, expected: 0b1000},
		{input: 5, expected: 0b10000},
		{input: 6, expected: 0b100000},
		{input: 7, expected: 0b1000000},
		{input: 8, expected: 0b10000000},
		{input: 9, expected: 0b100000000},
		{input: 10, expected: 0b1000000000},
		{input: 11, expected: 0b10000000000},
		{input: 12, expected: 0b100000000000},
		{input: 13, expected: 0b1000000000000},
		{input: 14, expected: 0b10000000000000},
		{input: 15, expected: 0b100000000000000},
		{input: 16, expected: 0b1000000000000000},
		{input: 17, expected: 0b0},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("NewSingleCandidate(%d)", testCase.input),
			func(t *testing.T) {
				assert.Equal(t, testCase.expected, puzzle.NewSingleCandidate(testCase.input))
			},
		)
	}
}

func TestNewCandidatesForRange(t *testing.T) {
	testCases := []struct {
		input    uint8
		expected puzzle.Candidates
	}{
		{input: 0, expected: 0b0},
		{input: 1, expected: 0b1},
		{input: 2, expected: 0b11},
		{input: 3, expected: 0b111},
		{input: 4, expected: 0b1111},
		{input: 5, expected: 0b11111},
		{input: 6, expected: 0b111111},
		{input: 7, expected: 0b1111111},
		{input: 8, expected: 0b11111111},
		{input: 9, expected: 0b111111111},
		{input: 10, expected: 0b1111111111},
		{input: 11, expected: 0b11111111111},
		{input: 12, expected: 0b111111111111},
		{input: 13, expected: 0b1111111111111},
		{input: 14, expected: 0b11111111111111},
		{input: 15, expected: 0b111111111111111},
		{input: 16, expected: 0b1111111111111111},
		{input: 17, expected: 0b1111111111111111}, // 17th value is not added
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("NewCandidatesForRange(%d)", testCase.input),
			func(t *testing.T) {
				assert.Equal(t, testCase.expected, puzzle.NewCandidatesForRange(testCase.input))
			},
		)
	}
}

func TestCandidatesUnion(t *testing.T) {
	testCases := []struct {
		name         string
		initialValue puzzle.Candidates
		adders       []puzzle.Candidates
		expected     puzzle.Candidates
	}{
		{
			name:         "0 and 1",
			initialValue: 0b0,
			adders:       []puzzle.Candidates{0b1},
			expected:     0b1,
		},
		{
			name:         "1 and 2",
			initialValue: 0b01,
			adders:       []puzzle.Candidates{0b10},
			expected:     0b11,
		},
		{
			name:         "1 and 1",
			initialValue: 0b1,
			adders:       []puzzle.Candidates{0b1},
			expected:     0b1,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				value := testCase.initialValue
				for _, add := range testCase.adders {
					value = value.Union(add)
				}
				assert.Equal(t, testCase.expected, value)
			},
		)
	}
}

func TestCandidatesRemove(t *testing.T) {
	testCases := []struct {
		name         string
		initialValue puzzle.Candidates
		substracters []puzzle.Candidates
		expected     puzzle.Candidates
	}{
		{
			name:         "drop 1 from 1,2",
			initialValue: 0b11,
			substracters: []puzzle.Candidates{0b01},
			expected:     0b10,
		},
		{
			name:         "drop 2 from 1,3",
			initialValue: 0b101,
			substracters: []puzzle.Candidates{0b010},
			expected:     0b101,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				value := testCase.initialValue
				for _, sub := range testCase.substracters {
					value = value.Remove(sub)
				}
				assert.Equal(t, testCase.expected, value)
			},
		)
	}
}

func TestCandidatesAll(t *testing.T) {
	testCases := []struct {
		name     string
		input    puzzle.Candidates
		expected []puzzle.Candidates
	}{
		{name: "0", input: 0b0, expected: nil},
		{name: "1", input: 0b1, expected: []puzzle.Candidates{0b1}},
		{name: "1 and 2", input: 0b11, expected: []puzzle.Candidates{0b01, 0b10}},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				actual := slices.Collect(testCase.input.All())
				assert.Equal(t, testCase.expected, actual)
			},
		)
	}
}

func TestCandidatesString(t *testing.T) {
	testCases := []struct {
		name     string
		input    puzzle.Candidates
		expected string
	}{
		{name: "0", input: 0b0, expected: ""},
		{name: "1", input: 0b1, expected: "1"},
		{name: "1 and 2", input: 0b11, expected: "12"},
		{name: "1, 2, 10, and 16", input: 0b1000001000000011, expected: "12ag"},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				assert.Equal(t, testCase.expected, testCase.input.String())
			},
		)
	}
}

func TestCandidatesCount(t *testing.T) {
	testCases := []struct {
		name     string
		input    puzzle.Candidates
		expected int
	}{
		{name: "0", input: 0b0, expected: 0},
		{name: "1", input: 0b1, expected: 1},
		{name: "1 and 2", input: 0b11, expected: 2},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				assert.Equal(t, testCase.expected, testCase.input.Count())
			},
		)
	}
}
