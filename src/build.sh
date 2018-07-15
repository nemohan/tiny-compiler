#!/bin/sh
go build tinycc.go scanner.go parser.go tree.go symtab.go \
    semantic.go

