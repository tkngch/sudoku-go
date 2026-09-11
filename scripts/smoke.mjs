// smoke.mjs loads the WebAssembly module and calls it. It proves that the
// module starts, that it installs the global object `sudoku`, and that `solve`
// answers. This script reads build/web/.

import { readFileSync } from "node:fs";

const webDir = new URL("../build/web/", import.meta.url);

const puzzle = ".2343.1243.1214.";
const solution = "1234341243212143";

let failures = 0;

function check(label, actual, expected) {
	if (actual !== expected) {
		console.error(`FAIL ${label}: want ${JSON.stringify(expected)}, got ${JSON.stringify(actual)}`);
		failures += 1;
	}
}

function runChecks() {
	const ok = globalThis.sudoku.solve(puzzle);
	check("solve.ok", ok.ok, true);
	check("solve.solution", ok.solution, solution);
	check("solve.kind", ok.kind, "");

	const bad = globalThis.sudoku.solve("z234" + "3.12" + "43.1" + "214.");
	check("character.ok", bad.ok, false);
	check("character.solution", bad.solution, "");
	check("character.kind", bad.kind, "character");

	const short = globalThis.sudoku.solve("123");
	check("size.kind", short.kind, "size");

	const noArg = globalThis.sudoku.solve();
	check("noArg.ok", noArg.ok, false);
	check("noArg.kind", noArg.kind, "unknown");

	if (failures > 0) {
		console.error(`smoke failed: ${failures} check(s)`);
		process.exit(1);
	}

	console.log("smoke ok");
	process.exit(0);
}

// Stop the job when the module does not report its readiness.
setTimeout(() => {
	console.error("smoke failed: __sudokuReady did not run within 30s");
	process.exit(1);
}, 30_000);

// Run the checks from a timer.
//
// Go calls __sudokuReady from main(), and main() is still on the stack.
// setTimeout puts runChecks on the task queue instead. __sudokuReady returns
// when main() reaches select {}, and Go becomes idle. The timer then fires, and
// runChecks calls solve as a new, top-level entry into Go.
globalThis.__sudokuReady = () => setTimeout(runChecks, 0);

// wasm_exec.js sets globalThis.Go.
await import(new URL("wasm_exec.js", webDir).href);

const go = new Go();
const { instance } = await WebAssembly.instantiate(
	readFileSync(new URL("sudoku.wasm", webDir)),
	go.importObject,
);

go.run(instance);
