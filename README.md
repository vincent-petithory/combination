# combination [![Build Status](https://travis-ci.org/vincent-petithory/combination.svg?branch=master)](https://travis-ci.org/vincent-petithory/combination)

Combination is a small program to help generating test tables data.

It is primarily intended for Go as it's common to use test-tables for tests.
Hence the output is valid Go syntax ready to be pasted in your buffer.

To install, first [install Go](http://golang.org/doc/install) then run:

    go get github.com/vincent-petithory/combination

Example:

    cat << EOF | combination
    card: ["\"Heart\"", "\"Tile\"", "\"Clover\"", "\"Pike\""]
    figure: ["\"Jack\"", "\"Queen\"", "\"King\""]
    EOF

For its usage, see [![GoDoc](https://godoc.org/github.com/vincent-petithory/combination?status.svg)](https://godoc.org/github.com/vincent-petithory/combination)
