package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tkngch/sudoku-go/internal/puzzle"
	"github.com/tkngch/sudoku-go/internal/solver"
)

func main() {
	flag.Usage = func() {
		_, err := fmt.Fprintf(
			flag.CommandLine.Output(),
			"usage: %s <sudoku_puzzle>\n", os.Args[0],
		)
		if err != nil {
			os.Exit(1)
		}

		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	sudoku, err := puzzle.Parse(flag.Arg(0))
	if err != nil {
		os.Exit(3)
	}

	_, err = fmt.Fprintf(os.Stderr, "Sudoku\n%s\n", sudoku.Render())
	if err != nil {
		os.Exit(4)
	}

	solution, err := solver.Solve(sudoku)
	if err != nil {
		os.Exit(5)
	}

	_, err = fmt.Fprintf(os.Stderr, "Solution\n%s\n", solution.Render())
	if err != nil {
		os.Exit(6)
	}

	_, err = fmt.Fprintf(os.Stdout, "%s\n", solution.String())
	if err != nil {
		os.Exit(7)
	}
}
