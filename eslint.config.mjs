// eslint.config.mjs checks the browser scripts and the smoke test.
//
// Each file below runs in a different environment, so each one gets its own
// block.

import js from "@eslint/js";
import globals from "globals";

// wasmExecGlobals holds the names that wasm_exec.js adds. That file comes from
// the Go distribution, and the makefile copies it into the build directory. No
// file in this repository declares these names.
const wasmExecGlobals = {
	Go: "readonly",
};

export default [
	{
		// The build directory holds a copy of the scripts and the Go runtime.
		ignores: ["build/"],
	},
	{
		// web/app.js runs on the main thread. It is a classic script, so it
		// uses no import and no export.
		files: ["web/app.js"],
		...js.configs.recommended,
		languageOptions: {
			ecmaVersion: "latest",
			sourceType: "script",
			globals: globals.browser,
		},
	},
	{
		// web/worker.js runs in a classic worker, and importScripts loads
		// wasm_exec.js into it.
		files: ["web/worker.js"],
		...js.configs.recommended,
		languageOptions: {
			ecmaVersion: "latest",
			sourceType: "script",
			globals: {
				...globals.worker,
				...wasmExecGlobals,
			},
		},
	},
	{
		// scripts/smoke.mjs runs on node, and it is an ES module. It imports
		// wasm_exec.js for its side effect.
		files: ["scripts/*.mjs"],
		...js.configs.recommended,
		languageOptions: {
			ecmaVersion: "latest",
			sourceType: "module",
			globals: {
				...globals.nodeBuiltin,
				...wasmExecGlobals,
			},
		},
	},
];
