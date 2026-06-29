default: format vet lint test

.PHONY: format
format:
	golangci-lint fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	go fmt ./...
	golangci-lint run ./...

.PHONY: test
test:
	go test -cover -race -timeout 30s ./...
