SHELL := bash
SOURCES := $(shell find . -name '*.go')

default: format vet lint test

.PHONY: format
format:
	golangci-lint fmt ./...
	golangci-lint run --fix ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	gofmt -l .
	golangci-lint run ./...

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
