package main

import (
	"fmt"
	"os"
	"strings"
)

var enableTrace = false
var loadFile = ""
var inputCh = make(chan string, 1)
var vmInput = make(chan string, 1)
var vmOutput = make(chan string, 1)

type tmInput struct {
	cmd   string
	param string
}

type tmCmd struct {
	name      string
	usage     string
	handler   func() bool
	param     string
	needInput bool
}

var cmdTable = map[string]tmCmd{
	"next":     tmCmd{name: "next", usage: "next step", handler: nextStep},
	"continue": tmCmd{name: "continue", usage: "leave single step mode", handler: leaveSingle},
	"exit":     tmCmd{name: "exit", usage: "exit tiny vm", handler: exit},
	"load":     tmCmd{name: "load", usage: "load tiny source file", handler: load},
	"single":   tmCmd{name: "single", usage: " enter single step mode", handler: enterSingle},
	"debug":    tmCmd{name: "debug", usage: "trace the whole process", handler: debug},
	"run":      tmCmd{name: "run", usage: "run loaded program", handler: run},
	"help":     tmCmd{name: "help", usage: "tiny vm help information", handler: nil},
}

func exit() bool {
	return true
}
func enterSingle() bool {
	tvm.enableSingleStep()
	return false
}

func leaveSingle() bool {
	tvm.disableSingleStep()
	return false
}

func nextStep() bool {
	tvm.nextStep()
	return false
}

func debug() bool {
	enableTrace = true
	initTrace()
	return false
}

func run() bool {
	go tvm.processor()
	return false
}
func load() bool {
	if err := readFile(loadFile); err != nil {
		fmt.Printf("%v\n", err)
		return false
	}
	ast := Parse()
	Analysis(ast)
	DumpSymbolTable()
	GenCode(ast)
	tvm.dumpInstructions()
	return false
}

func help() bool {
	for name, cmd := range cmdTable {
		fmt.Printf("\t%s\t%s\n", name, cmd.usage)
	}
	return true
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

func main() {
	fmt.Printf("tiny machine for tiny language 0.1.0\n")
	loop()
}

func loop() {
	inParam := ""
	for {
		fmt.Printf("tm-> ")
		fmt.Fscanln(os.Stdin, &inParam)
		in := parseInput(inParam)
		if tvm.Input(in.cmd) {
			continue
		}
		cmdName := in.cmd
		cmd, ok := cmdTable[cmdName]
		if cmd.name == "load" {
			loadFile = in.param
		}

		if !ok {
			fmt.Printf("tm-> unkown command %s\n", cmdName)
			continue
		}
		cmd.param = in.param
		if cmd.handler != nil && cmd.handler() {
			return
		}
		if cmd.handler == nil && cmd.name == "help" {
			help()
		}
	}
}
