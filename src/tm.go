package main

import (
	"fmt"
	"os"
)

var enableTrace = false

func load(file string) {
	readFile(file)
	ast := Parse()
	Analysis(ast)
	DumpSymbolTable()
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

func parseInput(in string) {

}

func loop() {
	cmd := ""
	for {
		fmt.Printf("tm->: ")
		fmt.Fscanln(os.Stdin, &cmd)
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

		}
		if cmd == "debug" {
			initTrace()
		}
		fmt.Printf("tm->:%s\n", cmd)
	}
}
