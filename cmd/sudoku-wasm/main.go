//go:build js && wasm

// Command sudoku-wasm exposes the solver to the browser. It installs the global
// object `sudoku`, which holds one method, `solve`.
package main

import (
	"syscall/js"

	"github.com/tkngch/sudoku-go/internal/web"
)

func main() {
	js.Global().Set("sudoku", js.ValueOf(map[string]any{
		"solve": js.FuncOf(solve),
	}))

	// Report that the solver is ready.
	//
	// The browser loads this module asynchronously, so a caller cannot know
	// when the global object sudoku appears. A caller therefore can install the
	// function __sudokuReady before it loads the module.
	if ready := js.Global().Get("__sudokuReady"); ready.Type() == js.TypeFunction {
		ready.Invoke()
	}

	// Keep the module alive, so the browser may call solve later. This blocks
	// no event: the js/wasm scheduler starts an event handler whenever it
	// idles.
	select {}
}

// solve reads the puzzle from the first argument and returns an object with
// three fields:
//
//   - ok: true after a success, and false after a failure.
//   - solution: the solved puzzle, or an empty string after a failure.
//   - kind: an empty string after a success, or the label of the fault. The
//     labels are the values of web.ErrorKind.
func solve(_ js.Value, args []js.Value) any {
	const wantArgCount = 1

	if len(args) != wantArgCount || args[0].Type() != js.TypeString {
		return result("", web.ErrorKindUnknown)
	}

	return result(web.Solve(args[0].String()))
}

// result builds the object that solve returns. It converts the kind to a plain
// string, because js.ValueOf accepts no named string type.
func result(solution string, errKind web.ErrorKind) map[string]any {
	return map[string]any{
		"ok":       errKind == "",
		"solution": solution,
		"kind":     string(errKind),
	}
}
