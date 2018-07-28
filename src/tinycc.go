package main

import (
	"flag"
	"fmt"
)

var sourceFile = ""
var enableTrace = false

func main() {
	flag.StringVar(&sourceFile, "c", "", "tiny source file")
	flag.BoolVar(&enableTrace, "t", false, "trace compile process")
	//flag.StringVar(&logFile, "l", "", "trace output file")
	flag.Parse()
	if sourceFile == "" {
		fmt.Printf("no source file\n usage: tinycc -c source\n")
		return
	}
	if enableTrace {
		initTrace()
	}
	//dumpWithoutLine()
	//dumpWithLine()
	readFile(sourceFile)
	ast := Parse()
	Analysis(ast)
	DumpSymbolTable()

	Logf("\ngenerated code ==================\n")
	genCode(ast)
	dumpInstructions()

}
