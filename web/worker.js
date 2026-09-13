// worker.js owns the WebAssembly module. It loads the module, then it answers
// each solve request from the page.
//
// This worker runs on its own thread, so a solve blocks no part of the user
// interface. Go blocks inside solve, so this worker reads no message while a
// solve runs. To stop a solve, the page terminates the worker.
//
// This file is a classic worker script, and a module worker rejects
// importScripts.

importScripts("./wasm_exec.js");

// main() calls __sudokuReady when the global object sudoku appears. main() is
// still on the stack at that moment, and a call back into Go is unsafe there.
// This handler therefore resolves a promise and does nothing more. The first
// call into Go then arrives from onmessage, which is a separate task.
const ready = new Promise((resolve) => {
	self.__sudokuReady = resolve;
});

const go = new Go();

// instantiate compiles the module. It prefers the streaming form, which starts
// the compile before the download completes. A server that sends the wrong
// content type breaks that form, so the fallback reads the whole body first.
async function instantiate() {
	const url = "./sudoku.wasm";

	try {
		const streamed = await WebAssembly.instantiateStreaming(fetch(url), go.importObject);

		return streamed.instance;
	} catch {
		const bytes = await (await fetch(url)).arrayBuffer();
		const source = await WebAssembly.instantiate(bytes, go.importObject);

		return source.instance;
	}
}

async function start() {
	const instance = await instantiate();

	// go.run never settles, because main() ends in select {}. Do not await it.
	go.run(instance);

	await ready;
	self.postMessage({ type: "ready" });
}

self.onmessage = (event) => {
	const request = event.data;

	if (request.type !== "solve") {
		return;
	}

	if (!self.sudoku) {
		self.postMessage({ type: "error", reason: "the module is not ready" });

		return;
	}

	const began = performance.now();
	const result = self.sudoku.solve(request.puzzle);
	const elapsed = performance.now() - began;

	self.postMessage({
		type: "result",
		id: request.id,
		ok: result.ok,
		solution: result.solution,
		kind: result.kind,
		ms: elapsed,
	});
};

start().catch((error) => {
	self.postMessage({ type: "error", reason: String(error) });
});
