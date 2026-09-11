SHELL := bash
SOURCES := $(shell find . -name '*.go')

default: format vet lint test

.PHONY: format
format:
	golangci-lint fmt ./...
	golangci-lint run --fix ./...
	GOOS=js GOARCH=wasm golangci-lint fmt ./...
	GOOS=js GOARCH=wasm golangci-lint run --fix ./...

.PHONY: vet
vet:
	go vet ./...
	GOOS=js GOARCH=wasm go vet ./...

.PHONY: lint
lint:
	gofmt -l .
	golangci-lint run ./...
	GOOS=js GOARCH=wasm golangci-lint run ./...

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

.PHONY: web
web: $(WEB_DIR)/sudoku.wasm $(WEB_DIR)/wasm_exec.js $(WEB_STATIC)

$(WEB_DIR)/sudoku.wasm: $(SOURCES)
	@mkdir -p $(WEB_DIR)
	GOOS=js GOARCH=wasm go build -ldflags='-s -w' -o $@ ./cmd/sudoku-wasm

$(WEB_DIR)/wasm_exec.js: $(WASM_EXEC)
	@mkdir -p $(WEB_DIR)
	cp $< $@

$(WEB_DIR)/%: web/%
	@mkdir -p $(WEB_DIR)
	cp $< $@

.PHONY: smoke
smoke: web
	node scripts/smoke.mjs

.PHONY: clean
clean:
	rm -rf build

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
