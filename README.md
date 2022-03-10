# l5x

A parser for the RSLogix5000 L5X format

## `l5x-to-go`

This parser was designed to generate go structures for L5X tags for use with [`go-plc`](https://github.com/stellentus/go-plc).

### Example Usage

*Generate a file `gen.go` in a package `types` which contains all of the types in your program file.*
`go run github.com/stellentus/l5x/cmd/l5x-to-go -package types -in PATH/TO/PROGRAM.L5X -out gen.go`

*To generate the same as above, but using a `go generate` command in your project, add the following comment at the beginning of a file in the `types` package.*
`//go:generate go run github.com/stellentus/l5x/cmd/l5x-to-go -package types -in PATH/TO/PROGRAM.L5X -out gen.go`

## `tag-descriptions`

Tag descriptions can be generated as a CSV. This does not print fully qualified tag names. It might be incomplete.

### Example Usage

*Generate a file `tag-descriptions.csv` which is a 2-column CSV of the names and descriptions of all tags, including program tags, in your program file.*
`go run cmd/comments-for-tags/main.go -in PATH/TO/PROGRAM.L5X -out tag-descriptions.csv`
