# sudoku-go

A solver for variable-size Sudoku puzzles — 4×4, 6×6, 9×9, and 12×12.

## How to use

```sh
make build                 # produces ./build/sudoku
./build/sudoku <puzzle>      # solve a puzzle passed as an argument
echo "<puzzle>" | ./build/sudoku   # or pipe it in on stdin
```

Without building, you can also use:

```sh
go run ./cmd/sudoku <puzzle>
```

## Input format

A puzzle is a string of one character per cell, in row-major order. Its length
determines the layout:

| Length | Grid Size |
| -----: | --------: |
|     16 | 4×4       |
|     36 | 6×6       |
|     81 | 9×9       |
|    144 | 12×12     |

Each character represents a cell:

- `0` or `.`: an empty cell.
- `1`-`9` and `a`-`g` / `A`-`G`: a given value (the letters cover 10–16 for the
  larger boards).

Whitespace and linebreak are ignored, so a puzzle may be supplied as one line or
pasted/piped as a grid across several lines.

## Output

Output is split across two streams so the result is easy to capture or pipe:

- **stdout** — the solved puzzle in the same compact, one-line form as the input
  (machine-readable).
- **stderr** — the input and the solution pretty-printed (human-readable).

So `./sudoku/sudoku <puzzle> > solution.txt` writes only the compact solution to
the file, while the pretty-printed grids appear on the terminal.

## Exit codes

| Code | Meaning                                                                                 |
| ---: | --------------------------------------------------------------------------------------- |
|    0 | success                                                                                 |
|    1 | error: the puzzle could not be read, parsed, or solved (a message is printed to stderr) |
|    2 | usage error: wrong number of arguments, or a bad flag                                   |

## Example

```sh
$ ./build/sudoku '.2343.1243.1214.'
```

stdout (the compact solution):

```
1234341243212143
```

stderr (the input and the solution, rendered):

```
Sudoku
+-----+-----+
| . 2 | 3 4 |
| 3 . | 1 2 |
+-----+-----+
| 4 3 | . 1 |
| 2 1 | 4 . |
+-----+-----+
Solution
+-----+-----+
| 1 2 | 3 4 |
| 3 4 | 1 2 |
+-----+-----+
| 4 3 | 2 1 |
| 2 1 | 4 3 |
+-----+-----+
```
