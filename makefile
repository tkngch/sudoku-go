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
	go test -timeout 1m -run=^$$ -bench=. -benchmem ./...

BIN := build/sudoku
.PHONY: build
build: $(BIN)

$(BIN): $(wildcard cmd/sudoku/*) \
	$(wildcard internal/puzzle/*) \
	$(wildcard internal/solver/*) \
	$(wildcard internal/sudoku/*)
	go build -o $@ ./cmd/sudoku
