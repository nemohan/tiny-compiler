package main

import (
	"fmt"
	"os"
	"strings"
)

var enableTrace = false

type tmInput struct {
	cmd   string
	param string
}

func load(file string) {
	readFile(file)
	ast := Parse()
	Analysis(ast)
	DumpSymbolTable()
	GenCode(ast)
	dumpInstructions()
}

func main() {
	fmt.Printf("tiny machine for tiny language 0.1.0\n")
	loop()
}

func usage() {
	fmt.Printf("\tload(file) \tload a tiny source file\n")
	fmt.Printf("\texit  \texit tiny machine\n")
	fmt.Printf("\trun  \trun loaded file\n")
	fmt.Printf("\tdebug \t enable generate compile trace file\n")
}

func parseInput(src string) *tmInput {
	in := &tmInput{}
	paramBegin := strings.Index(src, "(")
	if paramBegin == -1 {
		in.cmd = src
	} else {
		in.cmd = string(src[:paramBegin])
		in.param = string(src[paramBegin+1 : len(src)-1])
	}
	return in
}

func loop() {
	cmd := ""
	for {
		fmt.Printf("tm-> ")
		fmt.Fscanln(os.Stdin, &cmd)
		in := parseInput(cmd)
		cmd = in.cmd
		if cmd == "exit" {
			return
		}
		if cmd == "help" {
			usage()
			continue
		}
		if cmd == "run" {

		}
		if cmd == "load" {
			load(in.param)
		}
		if cmd == "debug" {
			enableTrace = true
			initTrace()
		}
		fmt.Printf("tm->:%s\n", cmd)
	}
}
