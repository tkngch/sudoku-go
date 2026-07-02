package solver_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tkngch/sudoku-go/internal/puzzle"
	"github.com/tkngch/sudoku-go/internal/solver"
)

func TestSolve(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		input         *string
		expected      *string
		expectedError error
	}{
		{
			name:          "4x4 without backtracking",
			input:         new(".234" + "3.12" + "43.1" + "214."),
			expected:      new("1234" + "3412" + "4321" + "2143"),
			expectedError: nil,
		},
		{
			// 6x6 uses rectangular 2x3 blocks (2 rows, 3 columns), a block
			// shape the 4x4 and 9x9 cases never exercise. One cell is blanked
			// per row, so each blank is resolved by propagating the provided
			// values.
			name:          "6x6 without backtracking",
			input:         new(".23456" + "4.6123" + "23.564" + "564.31" + "3126.5" + "64531."),
			expected:      new("123456" + "456123" + "231564" + "564231" + "312645" + "645312"),
			expectedError: nil,
		},
		{
			name: "9x9 without backtracking",
			input: new("812753649" + "940080170" + "005491203" + "054237090" +
				"369845020" + "207069504" + "521970360" + "438526917" + "796318452"),
			expected: new("812753649" + "943682175" + "675491283" + "154237896" +
				"369845721" + "287169534" + "521974368" + "438526917" + "796318452"),
			expectedError: nil,
		},
		{
			name: "9x9 with backtracking",
			input: new("800000000" + "003600000" + "070090200" + "050007000" +
				"000045700" + "000100030" + "001000068" + "008500010" + "090000400"),
			expected: new("812753649" + "943682175" + "675491283" + "154237896" +
				"369845721" + "287169534" + "521974368" + "438526917" + "796318452"),
			expectedError: nil,
		},
		{
			name:          "4x4 with no solution",
			input:         new("11.." + "...." + "...." + "...."),
			expected:      nil,
			expectedError: solver.ErrSolutionNotFound,
		},
		{
			name:          "4x4 already solved",
			input:         new("1234" + "3412" + "4321" + "2143"),
			expected:      new("1234" + "3412" + "4321" + "2143"),
			expectedError: nil,
		},
		{
			// The "9x9 with backtracking" puzzle with (0,1) pinned to 2,
			// whereas its unique solution needs 1 there. Solve must exhaust the
			// search and fail, unlike "4x4 with no solution" which is rejected
			// before the search begins.
			name: "9x9 unsolvable, fails during search",
			input: new("820000000" + "003600000" + "070090200" + "050007000" +
				"000045700" + "000100030" + "001000068" + "008500010" + "090000400"),
			expected:      nil,
			expectedError: solver.ErrSolutionNotFound,
		},
		{
			name:          "nil grid",
			input:         nil,
			expected:      nil,
			expectedError: solver.ErrInvalidGrid,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				var (
					grid     *puzzle.Grid
					err      error
					expected *puzzle.Grid
				)

				if testCase.input == nil {
					grid = nil
				} else {
					grid, err = puzzle.Parse(*testCase.input)
					require.NoErrorf(t, err, "could not parse [%v]", testCase.input)
				}

				actual, err := solver.Solve(context.Background(), grid)
				if testCase.expectedError != nil {
					require.ErrorIs(t, err, testCase.expectedError)
				} else {
					require.NoError(t, err, "could not find a solution")
				}

				if testCase.expected == nil {
					expected = nil
				} else {
					expected, err = puzzle.Parse(*testCase.expected)
					require.NoErrorf(t, err, "could not parse [%v]", testCase.expected)
				}

				assert.Equalf(
					t,
					expected,
					actual,
					"expected\n%s\nactual\n%s",
					expected.Render(),
					actual.Render(),
				)
			},
		)
	}
}

// cancelingContext returns nil for its first n Err() calls, then
// context.Canceled.
type cancelingContext struct {
	calls, canceledAfter int
}

func newCancelingContext(canceledAfter int) *cancelingContext {
	return &cancelingContext{
		calls:         0,
		canceledAfter: canceledAfter,
	}
}

func (c *cancelingContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (c *cancelingContext) Done() <-chan struct{}       { return nil }
func (c *cancelingContext) Value(any) any               { return nil }

func (c *cancelingContext) Err() error {
	c.calls++
	if c.calls <= c.canceledAfter {
		return nil
	}

	return context.Canceled
}

func TestSolveContext(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		ctx         func() (context.Context, context.CancelFunc)
		expectedErr error
	}{
		{
			name: "timed-out",
			ctx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithDeadline(context.Background(), time.Unix(0, 0))

				return ctx, cancel
			},
			expectedErr: context.DeadlineExceeded,
		},
		{
			name: "already canceled",
			ctx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				return ctx, nil
			},
			expectedErr: context.Canceled,
		},
		{
			// Cancellation should be honored inside the search, not only at
			// Solve's entry guard.
			name:        "canceled early",
			ctx:         func() (context.Context, context.CancelFunc) { return newCancelingContext(1), nil },
			expectedErr: context.Canceled,
		},
		{
			// Cancellation should be honored inside the search, even deep in
			// several search-recursion.
			name:        "canceled deep in the search",
			ctx:         func() (context.Context, context.CancelFunc) { return newCancelingContext(4), nil },
			expectedErr: context.Canceled,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				grid, err := puzzle.Parse(strings.Repeat(".", 16))
				require.NoError(t, err)

				ctx, cancel := testCase.ctx()
				if cancel != nil {
					defer cancel()
				}

				solution, err := solver.Solve(ctx, grid)
				require.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, solution)
			},
		)
	}
}
