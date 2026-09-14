// app.js paints the grid, reads the keyboard, and talks to the worker.
//
// The page supports the 9x9 grid only. The solver supports every size, and the
// command line tool exposes them all.
//
// This file is a classic script. It uses no import and no export, so
// `node --check` reads it as CommonJS.

"use strict";

(function () {
	const GRID_SIZE = 9;
	const BLOCK_ROWS = 3;
	const BLOCK_COLS = 3;
	const CELL_COUNT = GRID_SIZE * GRID_SIZE;

	// TIMEOUT_MS is the limit for one solve.
	const TIMEOUT_MS = 10000;

	const EMPTY = ".";
	const HASH_PREFIX = "#p=";
	const EMPTY_PUZZLE = EMPTY.repeat(CELL_COUNT);

	// SKIP_KEYS empty a cell and move to the next one. A user who copies a
	// puzzle then types it from left to right, with no pause at an empty cell.
	//
	// The dot and the zero both mark an empty cell in the input format, which
	// the README documents. The space bar is the third natural way to skip.
	const SKIP_KEYS = [EMPTY, "0", " "];

	// MESSAGES holds one sentence per error kind of internal/web.ErrorKind.
	const MESSAGES = {
		size: `The puzzle must hold ${CELL_COUNT} cells.`,
		character: "The puzzle holds a character that is not 1-9 or a dot.",
		value: `The puzzle holds a value that is too large for a ${GRID_SIZE}x${GRID_SIZE} grid.`,
		unsolvable: "This puzzle has no solution.",
		unknown: "The solver reported an unexpected fault.",
	};

	const TIMEOUT_MESSAGE = `The solver ran for more than ${TIMEOUT_MS / 1000} seconds, so the page stopped it.`;
	const LOAD_ERROR_MESSAGE = "The solver failed to load. Reload the page to try again.";
	const EXAMPLES_ERROR_MESSAGE = "The example puzzles failed to load.";
	const HASH_ERROR_MESSAGE = "The link holds a puzzle that the page cannot read.";

	const gridElement = document.getElementById("grid");
	const solveButton = document.getElementById("solve");
	const clearButton = document.getElementById("clear");
	const exampleSelect = document.getElementById("example");
	const statusElement = document.getElementById("status");
	const alertElement = document.getElementById("alert");

	let cells = [];
	let worker = null;

	// requestId counts the solve requests. pendingId names the request that the
	// page still waits for, or 0 when no request runs.
	let requestId = 0;
	let pendingId = 0;
	let watchdog = 0;

	// pendingPuzzle holds the grid that the open request asked about.
	let pendingPuzzle = "";

	// solutionShown reports whether the grid holds values from the solver.
	let solutionShown = false;

	function setStatus(text) {
		statusElement.textContent = text;
	}

	function showError(text) {
		alertElement.textContent = text;
	}

	function clearError() {
		alertElement.textContent = "";
	}

	// buildGrid creates the cells and returns them in row-major order.
	function buildGrid() {
		const fragment = document.createDocumentFragment();
		const inputs = [];

		gridElement.style.setProperty("--grid-size", String(GRID_SIZE));

		for (let index = 0; index < CELL_COUNT; index += 1) {
			const row = Math.floor(index / GRID_SIZE);
			const col = index % GRID_SIZE;
			const input = document.createElement("input");

			input.type = "text";
			input.className = "cell";
			input.inputMode = "numeric";
			input.maxLength = 1;
			input.autocomplete = "off";
			input.spellcheck = false;
			input.setAttribute("autocapitalize", "off");
			input.setAttribute("autocorrect", "off");
			input.setAttribute("aria-label", `row ${row + 1} column ${col + 1}`);
			input.dataset.index = String(index);

			if (row % BLOCK_ROWS === 0) {
				input.classList.add("block-top");
			}

			if (col % BLOCK_COLS === 0) {
				input.classList.add("block-left");
			}

			inputs.push(input);
			fragment.append(input);
		}

		gridElement.append(fragment);

		return inputs;
	}

	// readGrid returns the grid as one line, with a dot for an empty cell.
	function readGrid() {
		return cells.map((cell) => (cell.value === "" ? EMPTY : cell.value)).join("");
	}

	// checkPuzzle mirrors puzzle.Parse and returns the error kind that the
	// solver would report. It returns an empty string for a good puzzle.
	//
	// The page calls this before it fills the grid from the URL hash. The grid
	// itself always holds a good puzzle, so it needs no check.
	function checkPuzzle(puzzle) {
		if (puzzle.length !== CELL_COUNT) {
			return "size";
		}

		if (/[^0-9a-gA-G.]/.test(puzzle)) {
			return "character";
		}

		// A letter is a valid character, but its value is 10 or more, which a
		// 9x9 grid rejects.
		if (/[a-gA-G]/.test(puzzle)) {
			return "value";
		}

		return "";
	}

	// clearSolution removes every value that the solver added, so an edit always
	// starts from the puzzle that the user typed. keepIndex names a cell that
	// keeps its value, because the user just typed into it.
	function clearSolution(keepIndex) {
		if (!solutionShown) {
			return;
		}

		solutionShown = false;

		for (const cell of cells) {
			if (!cell.classList.contains("solved")) {
				continue;
			}

			cell.classList.remove("solved");

			if (Number(cell.dataset.index) !== keepIndex) {
				cell.value = "";
			}
		}
	}

	function showSolution(solution) {
		for (let index = 0; index < CELL_COUNT; index += 1) {
			if (cells[index].value === "") {
				cells[index].value = solution[index];
				cells[index].classList.add("solved");
			}
		}

		solutionShown = true;
	}

	function fillGrid(puzzle) {
		clearSolution();

		for (let index = 0; index < CELL_COUNT; index += 1) {
			const char = puzzle[index];

			cells[index].classList.remove("solved");
			cells[index].value = char === EMPTY || char === "0" ? "" : char;
		}

		updateHash();
	}

	// updateHash records the puzzle in the URL, so a link carries it. It uses
	// replaceState, so the back button leaves the page instead of the last edit.
	function updateHash() {
		const puzzle = readGrid();
		const url = puzzle === EMPTY_PUZZLE
			? window.location.pathname + window.location.search
			: HASH_PREFIX + puzzle;

		window.history.replaceState(null, "", url);
	}

	// restoreFromHash fills the grid from the URL.
	function restoreFromHash(replaced) {
		const hash = window.location.hash;

		if (!hash.startsWith(HASH_PREFIX)) {
			return;
		}

		try {
			const puzzle = decodeURIComponent(hash.slice(HASH_PREFIX.length));
		} catch {
  		showError(HASH_ERROR_MESSAGE);
  		return;
		}

		const kind = checkPuzzle(puzzle);

		if (kind !== "") {
			showError(MESSAGES[kind]);

			return;
		}

		fillGrid(puzzle);
		clearError();

		if (replaced) {
			setStatus("");
		}
	}

	function focusCell(index, rowDelta, colDelta) {
		const row = Math.floor(index / GRID_SIZE) + rowDelta;
		const col = (index % GRID_SIZE) + colDelta;

		if (row < 0 || row >= GRID_SIZE || col < 0 || col >= GRID_SIZE) {
			return;
		}

		cells[row * GRID_SIZE + col].focus();
	}

	// focusNext moves to the next cell in row-major order. It stops at the last
	// cell, so the focus stays inside the grid.
	function focusNext(index) {
		if (index + 1 < CELL_COUNT) {
			cells[index + 1].focus();
		}
	}

	// focusPrevious moves to the cell before this one, in row-major order. It
	// stops at the first cell, so the focus stays inside the grid.
	function focusPrevious(index) {
		if (index > 0) {
			cells[index - 1].focus();
		}
	}

	function setCell(index, value) {
		clearSolution(index);
		cells[index].value = value;
		clearError();
		updateHash();
	}

	function onKeyDown(event) {
		const input = event.target;

		if (!input.classList.contains("cell")) {
			return;
		}

		const index = Number(input.dataset.index);
		const key = event.key;

		if (key === "ArrowUp" || key === "ArrowDown" || key === "ArrowLeft" || key === "ArrowRight") {
			event.preventDefault();
			focusCell(index, arrowRow(key), arrowCol(key));

			return;
		}

		// Backspace on an empty cell moves back, and it removes no value. One
		// more press then empties the cell that it moved to. A user therefore
		// walks backward through a row, and each press does one thing.
		if (key === "Backspace" && input.value === "") {
			event.preventDefault();
			focusPrevious(index);

			return;
		}

		// Backspace and Delete correct a cell, so the focus stays on it.
		if (key === "Backspace" || key === "Delete") {
			event.preventDefault();
			setCell(index, "");

			return;
		}

		if (SKIP_KEYS.includes(key)) {
			event.preventDefault();
			setCell(index, "");
			focusNext(index);

			return;
		}

		if (key >= "1" && key <= "9") {
			event.preventDefault();
			setCell(index, key);
			focusNext(index);

			return;
		}

		// Keep the keys that navigate or that hold a modifier. Reject every
		// other character, so no cell holds a value that the solver refuses.
		if (key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey) {
			event.preventDefault();
		}
	}

	function arrowRow(key) {
		if (key === "ArrowUp") {
			return -1;
		}

		return key === "ArrowDown" ? 1 : 0;
	}

	function arrowCol(key) {
		if (key === "ArrowLeft") {
			return -1;
		}

		return key === "ArrowRight" ? 1 : 0;
	}

	// onInput catches the paths that keydown misses, such as a mobile keyboard
	// and a paste into one cell. It keeps the last digit and drops the rest.
	function onInput(event) {
		const input = event.target;

		if (!input.classList.contains("cell")) {
			return;
		}

		const index = Number(input.dataset.index);
		const digits = input.value.replace(/[^1-9]/g, "");

		setCell(index, digits.slice(-1));

		// Move to the next cell after an insertion, so a mobile keyboard
		// behaves like a physical one. A deletion keeps the focus, which agrees
		// with Backspace and Delete above.
		if (String(event.inputType).startsWith("insert")) {
			focusNext(index);
		}
	}

	function onSolve() {
		clearSolution();
		clearError();

		// Read the grid after clearSolution, so the request holds the puzzle
		// that the user typed and no value from an earlier solve.
		const puzzle = readGrid();

		requestId += 1;
		pendingId = requestId;
		pendingPuzzle = puzzle;

		solveButton.disabled = true;
		setStatus("The solver runs.");

		watchdog = window.setTimeout(() => restart(TIMEOUT_MESSAGE), TIMEOUT_MS);
		worker.postMessage({ type: "solve", id: pendingId, puzzle });
	}

	function onClear() {
		fillGrid(EMPTY_PUZZLE);
		clearError();
		setStatus("");
		cells[0].focus();
	}

	function onExample() {
		const puzzle = exampleSelect.value;

		// Reset the control, so the user may load the same example twice.
		exampleSelect.selectedIndex = 0;

		if (puzzle === "") {
			return;
		}

		fillGrid(puzzle);
		clearError();
		setStatus("");
	}

	function endRequest() {
		window.clearTimeout(watchdog);
		watchdog = 0;
		pendingId = 0;
		pendingPuzzle = "";
	}

	// restart stops the worker and starts a new one. The worker reads no message
	// while Go runs inside solve, so termination is the only way to stop it.
	function restart(message) {
		endRequest();

		worker.terminate();
		worker = newWorker();

		solveButton.disabled = true;
		showError(message);
		setStatus("The solver restarts. Wait a moment.");
	}

	function onWorkerMessage(event) {
		const reply = event.data;

		if (reply.type === "ready") {
			solveButton.disabled = false;
			setStatus("The solver is ready.");

			return;
		}

		if (reply.type === "error") {
			endRequest();
			solveButton.disabled = true;
			showError(LOAD_ERROR_MESSAGE);
			setStatus("");

			return;
		}

		// Drop a reply from a request that the page abandoned.
		if (reply.id !== pendingId) {
			return;
		}

		// Drop a reply when the user changed the grid after the request.
		const stale = pendingPuzzle !== readGrid();

		endRequest();
		solveButton.disabled = false;

		if (stale) {
			setStatus("");

			return;
		}

		if (reply.ok) {
			showSolution(reply.solution);
			setStatus(`The solver finished in ${reply.ms.toFixed(1)} ms.`);

			return;
		}

		showError(MESSAGES[reply.kind] || MESSAGES.unknown);
		setStatus("");
	}

	function newWorker() {
		const created = new Worker("./worker.js");

		created.onmessage = onWorkerMessage;

		// End the open request here, as the error branch of onWorkerMessage
		// does. Without this call the watchdog stays active. It then replaces
		// this message with the timeout message and starts a second worker.
		created.onerror = () => {
			endRequest();
			solveButton.disabled = true;
			showError(LOAD_ERROR_MESSAGE);
			setStatus("");
		};

		return created;
	}

	async function loadExamples() {
		const response = await fetch("./examples.json");

		if (!response.ok) {
			throw new Error(`examples.json: ${response.status}`);
		}

		for (const entry of await response.json()) {
			const option = document.createElement("option");

			option.value = entry.puzzle;
			option.textContent = entry.name;
			exampleSelect.append(option);
		}

		exampleSelect.disabled = false;
	}

	function init() {
		cells = buildGrid();

		gridElement.addEventListener("keydown", onKeyDown);
		gridElement.addEventListener("input", onInput);
		solveButton.addEventListener("click", onSolve);
		clearButton.addEventListener("click", onClear);
		exampleSelect.addEventListener("change", onExample);
		window.addEventListener("hashchange", () => restoreFromHash(true));

		// Start the worker before the page reads the hash. The grid is already
		// on the screen, so it paints long before the module arrives. This
		// order also keeps the solver alive after a fault in the hash.
		worker = newWorker();

		restoreFromHash(false);

		loadExamples().catch(() => showError(EXAMPLES_ERROR_MESSAGE));
	}

	init();
})();
