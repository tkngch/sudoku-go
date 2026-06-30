package puzzle

import (
	"iter"
	"math/bits"
	"strings"
)

// Candidates holds values of a cell, where value x is represented with the x-th
// lowest bit. For example,
// - 0b001 indicates that the cell takes the value of 1;
// - 0b010 indicates that the cell takes the value of 2;
// - 0b011 indicates that the cell takes the value of 1 or 2; and
// - 0b111 indicates that the cell takes the value of 1 or 2 or 3.
type Candidates uint16

// Given that Candidates has 16 bit width, the max value it can represent is 16.
const maxCandidateValue = 16

// NewCandidatesForRange makes a new candidate that represents the ranged
// values: from 1 to the provided max-value or 16, whichever is smaller.
func NewCandidatesForRange(maxValue int) Candidates {
	var candidates Candidates

	for value := 1; value <= maxValue && value <= maxCandidateValue; value++ {
		candidates |= NewSingleCandidate(value)
	}

	return candidates
}

// NewSingleCandidate makes a new candidate that represents the provided value.
// Only values from 1 to 16 are supported. NewSingleCandidate(17), for example,
// returns an empty candidates.
func NewSingleCandidate(value int) Candidates {
	if value < 1 || value > maxCandidateValue {
		return 0
	}

	return 1 << (value - 1)
}

func (c Candidates) Union(other Candidates) Candidates {
	return c | other
}

func (c Candidates) Remove(other Candidates) Candidates {
	return c &^ other
}

// All iterates over each candidate as a single-bit Candidates, in ascending
// value order. For example, the set {1, 3} (represented as 0b101) yields 0b001
// then 0b100.
func (c Candidates) All() iter.Seq[Candidates] {
	return func(yield func(Candidates) bool) {
		for num := 1; num <= maxCandidateValue; num++ {
			value := NewSingleCandidate(num)
			if c&value != 0 && !yield(value) {
				return
			}
		}
	}
}

func (c Candidates) Count() int {
	return bits.OnesCount16(uint16(c))
}

func (c Candidates) String() string {
	vals := make([]string, 0, c.Count())

	for value := range c.All() {
		vals = append(vals, string(value.char()))
	}

	return strings.Join(vals, "")
}

// char returns the display char for the first candidate value. 1..9 are mapped
// to '1'..'9', and 10..16 are mapped 'a'..'g'.
func (c Candidates) char() byte {
	v := bits.TrailingZeros16(uint16(c)) + 1
	switch {
	case 1 <= v && v <= 9:
		return byte('0' + v)
	case 10 <= v && v <= 16:
		return byte('a' + v - 10)
	default:
		return '.'
	}
}
