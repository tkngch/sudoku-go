module github.com/tkngch/sudoku-go

go 1.26

// npm installs the lint tooling for the browser scripts into node_modules. One
// package there ships a Go file, and the pattern ./... finds it. This
// repository owns no file in that directory.
ignore ./node_modules

require github.com/stretchr/testify v1.11.1

require (
	github.com/aclements/go-moremath v0.0.0-20210112150236-f10218a38794 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/perf v0.0.0-20260615155930-9e4b9ddef5b6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

tool golang.org/x/perf/cmd/benchstat
