#!/bin/sh
echo "build tinycc"
go build tinycc.go scanner.go parser.go tree.go symtab.go \
    semantic.go trace.go code.go

echo "build tiny machine"
go build tm.go scanner.go parser.go tree.go symtab.go \
    semantic.go code.go trace.go

