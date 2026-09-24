SHELL := bash

# find skips node_modules, because npm installs the lint tooling there and one
# package ships a Go file.
SOURCES := $(shell find . -path ./node_modules -prune -o -name '*.go' -print)

default: format vet lint test smoke

.PHONY: format
format:
	golangci-lint fmt ./...
	golangci-lint run --fix ./...
	GOOS=js GOARCH=wasm golangci-lint fmt ./...
	GOOS=js GOARCH=wasm golangci-lint run --fix ./...
	prettier -w web/

.PHONY: vet
vet:
	go vet ./...
	GOOS=js GOARCH=wasm go vet ./...

.PHONY: lint
lint: node_modules
	@# gofmt -l prints the file names, but it exits with code 0. Test the
	@# output, so a bad format fails this target.
	@files=$$(gofmt -l $(SOURCES)); \
	if [ -n "$$files" ]; then \
		echo "These files are not gofmt-formatted:"; \
		echo "$$files"; \
		exit 1; \
	fi
	golangci-lint run ./...
	GOOS=js GOARCH=wasm golangci-lint run ./...
	npm run lint

.PHONY: test
test:
	go test -cover -race -timeout 30s ./...

.PHONY: bench
bench:
	@set -o pipefail; go test -run='^$$' -bench='^BenchmarkSolve$$' -benchmem -count=10 -timeout 20m ./internal/solver/ | go tool benchstat -

BIN := build/sudoku
.PHONY: build
build: $(BIN)

$(BIN): $(SOURCES)
	go build -o $@ ./cmd/sudoku

WEB_DIR    := build/web
WASM_EXEC  := $(shell go env GOROOT)/lib/wasm/wasm_exec.js
WEB_STATIC := $(patsubst web/%,$(WEB_DIR)/%,$(wildcard web/*))
PORT ?= 8080

.PHONY: web
web: $(WEB_DIR)/sudoku.wasm $(WEB_DIR)/wasm_exec.js $(WEB_STATIC)

# The page needs an HTTP server. A file:// URL blocks fetch, the worker, and
# the WebAssembly instantiation. Python 3.10 and later map .wasm to
# application/wasm, so instantiateStreaming works.
.PHONY: serve
serve: web
	python3 -m http.server --directory $(WEB_DIR) $(PORT)

.PHONY: smoke
smoke: web
	node scripts/smoke.mjs

.PHONY: clean
clean:
	rm -rf build

$(WEB_DIR)/sudoku.wasm: $(SOURCES)
	@mkdir -p $(WEB_DIR)
	GOOS=js GOARCH=wasm go build -ldflags='-s -w' -o $@ ./cmd/sudoku-wasm

$(WEB_DIR)/wasm_exec.js: $(WASM_EXEC)
	@mkdir -p $(WEB_DIR)
	cp $< $@

$(WEB_DIR)/%: web/%
	@mkdir -p $(WEB_DIR)
	cp $< $@

node_modules: package.json package-lock.json
	npm ci
	@touch $@

PROF_DIR  := build/prof
PROF_CPU := $(PROF_DIR)/cpu.prof
PROF_MEM := $(PROF_DIR)/mem.prof
PROF_OUT := $(PROF_DIR)/solver.test
PROF_SUM := $(PROF_DIR)/summary.txt

.PHONY: profile
profile: $(PROF_SUM)

$(PROF_CPU) $(PROF_MEM) $(PROF_OUT): $(SOURCES)
	@mkdir -p $(PROF_DIR)
	go test -run='^$$' -bench='^BenchmarkSolve$$' -benchmem -count=1 -benchtime=2s \
                -cpuprofile=$(PROF_CPU) -memprofile=$(PROF_MEM) \
                -o $(PROF_OUT) ./internal/solver/

$(PROF_SUM): $(PROF_CPU) $(PROF_MEM)
	@echo '== CPU: top by cum (expensive subtrees) ==' > $@
	@go tool pprof -top -cum -nodecount=10 $(PROF_CPU) >> $@
	@echo '== CPU: who allocates on-CPU ==' >> $@
	@go tool pprof -peek='mallocgc$$' $(PROF_CPU) >> $@
	@echo '== MEM: allocation COUNT by function =='  >> $@
	@go tool pprof -top -nodecount=10 -sample_index=alloc_objects $(PROF_MEM) >> $@
	@cat $@
	@echo ''
	@echo 'To profile a function line by line, call `go tool pprof -list=<function name> $(PROF_MEM)`'
