package puzzle

import (
	"iter"
	"math/bits"
	"strconv"
	"strings"
)

// Candidate values of a cell, where value x is represented with the x-th lowest bit.
// For example
// - 0001 indicates that the cell takes the value of 1;
// - 0010 indicates that the cell takes the value of 2;
// - 0011 indicates that the cell takes the value of 1 or 2; and
// - 0111 indicates that the cell takes the value of 1 or 2 or 3.
type Candidates uint16

// Given that Candidates has 16 bit width, the max value it can represent is 16.
const maxCandidateValue = 16

func NewCandidates(valueCount uint8) Candidates {
	d := Candidates(0)
	for value := uint8(1); value <= valueCount; value++ {
		d = d | NewCandidate(value)
	}
	return d
}

func NewCandidate(value uint8) Candidates {
	return 1 << (value - 1)
}

func (c Candidates) Add(other Candidates) Candidates {
	return c | other
}

func (c Candidates) Drop(other Candidates) Candidates {
	return c &^ other
}

func (c Candidates) Values() iter.Seq[Candidates] {
	return func(yield func(Candidates) bool) {
		for num := uint8(1); num <= maxCandidateValue; num++ {
			value := NewCandidate(num)
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

	for value := range c.Values() {
		i := bits.TrailingZeros16(uint16(value)) + 1
		vals = append(vals, strconv.Itoa(i))
	}

	return strings.Join(vals, ",")
}
