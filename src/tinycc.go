package main

import (
	"flag"
	"fmt"
)

var sourceFile = ""

func main() {
	flag.StringVar(&sourceFile, "c", "", "tiny source file")
	flag.Parse()
	if sourceFile == "" {
		fmt.Printf("no source file\n usage: tinycc -c source\n")
		return
	}
	//dumpWithoutLine()
	//dumpWithLine()
	readFile(sourceFile)
	Parse()
}
