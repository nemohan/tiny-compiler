package main

import (
	"container/list"
	"fmt"
)

type symbol struct {
	name     string
	line     int
	location int
}

var symbolTable = make(map[string]*symbol)
var table = make(map[string]*list.List)
var symLocation = 0

func insertSym(name string, line int) {
	_, ok := symbolTable[name]
	if ok {
		l := table[name]
		l.PushBack(&symbol{name: name, line: line})
		return
	}
	sym := &symbol{
		name:     name,
		line:     line,
		location: symLocation,
	}
	symLocation++
	symbolTable[name] = sym
	l := list.New()
	l.PushFront(sym)
	table[name] = l
}

func findSym(name string) int {
	s, ok := symbolTable[name]
	if !ok {
		panic(name)
	}
	return s.location
}

func delSym() {

}

func DumpSymbolTable() {
	fmt.Printf("variable\tName\tLocation\tLine\n")
	for name, s := range symbolTable {
		fmt.Printf("%s\t%d ", name, s.location)
		l := table[name]
		for e := l.Front(); e != nil; e = e.Next() {
			sym := e.Value.(*symbol)
			fmt.Printf("\t%d ", sym.line)
		}
		fmt.Printf("\n")
	}
}

func HandleSymProc(node *SyntaxTree) {
	if node.token == nil || node.token.tokenType != ID {
		return
	}
	insertSym(node.token.lexeme, node.token.line)
}

func Analysis(node *SyntaxTree) {
	GenTraverse(node, HandleSymProc, emptyTraverseProc)
	GenTraverse(node, emptyTraverseProc, checkType)
	GenTraverse(node, emptyTraverseProc, printTraverseProc)
	if isAnyTypeErr() {
		dumpTypeErr()
	}
}
