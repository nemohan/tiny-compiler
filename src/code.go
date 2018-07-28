package main

import (
	"strconv"
)

const (
	opNone = iota
	opHalt
	opIn
	opOut
	opAdd
	opSub
	opMul
	opDiv
	opLd
	opLda
	opLdc
	opSt
	opJlt
	opJle
	opJge
	opJgt
	opJEQ
	opJNE
)

var opTable = map[int]string{
	opNone: "none",
	opHalt: "halt",
	opIn:   "in",
	opOut:  "out",
	opAdd:  "add",
	opSub:  "sub",
	opMul:  "mul",
	opDiv:  "div",
	opLd:   "ld",
	opLda:  "lda",
	opLdc:  "ldc",
	opSt:   "st",
	opJlt:  "jlt",
	opJle:  "jle",
	opJge:  "jge",
	opJgt:  "jgt",
	opJEQ:  "jeq",
	opJNE:  "jne",
}

//registers
const (
	r0 = iota
	r1
	r2
	r3
	r4
	r5
	r6
	regPC
	regNone
)

var regTable = map[int]string{
	r0:      "r0",
	r1:      "r1",
	r2:      "r2",
	r3:      "r3",
	r4:      "r4",
	r5:      "r5",
	r6:      "r6",
	regPC:   "pc",
	regNone: "regNone",
}

type Instruction struct {
	opcode  int
	regs    []int
	handler func(*Instruction) (bool, error)
}

const regNum = 8
const iMemSize = 1024
const dMemSize = 1024

var registers = make([]int, regNum)
var iMem = make([]*Instruction, 0, iMemSize)
var dMem = make([]int, 0, dMemSize)
var currentDMemPos = 0

func GenCode(root *SyntaxTree) {
	genCode(root)
	emitROCode(opHalt, regNone, regNone, regNone)
}

func genCode(root *SyntaxTree) {
	if root == nil {
		return
	}
	node := root
	if root.nodeKind == fileK {
		node = root.child
	}
	if node.nodeKind == stmtK {
		genStmt(node)
	}
	if node.nodeKind == expK {
		Logf("gen exp\n")
		genExp(node)
	}
	for next := node.slibling; next != nil; next = next.slibling {
		genCode(next)
	}
}

func genStmt(node *SyntaxTree) {
	switch node.stmtKind {
	case ifK:
	case repeatK:
	case assignK:
		genAssign(node)
	case readK:
		offset := findSym(node.child.token.lexeme)
		emitROCode(opIn, r1, regNone, regNone)
		emitRMCode(opSt, r1, r5, offset)
	case writeK:
		genExp(node)
		emitROCode(opOut, r0, regNone, regNone)
	}
}

func genAssign(node *SyntaxTree) {
	offset := findSym(node.child.token.lexeme)
	for next := node.child.slibling; next != nil; next = next.slibling {
		genExp(next)
	}
	emitRMCode(opSt, r0, r5, offset)
}

func genExp(node *SyntaxTree) {
	if node == nil {
		return
	}
	switch node.token.tokenType {
	case tokenAdd:
		child := node.child
		genExp(node.child)
		for next := child.slibling; next != nil; next = next.slibling {
			genExp(next)
		}
		emitROCode(opAdd, r0, r1, r0)
	case tokenMinus:
	case tokenMultiply:
	case tokenDiv:
	case tokenId:
		offset := findSym(node.token.lexeme)
		emitRMCode(opLd, r0, r5, offset)
	case tokenNumber:
		//emitRMCode(opLdc, r1, strToInt(node.token.lexeme), regNone)
		emitRMCode(opLdc, r1, regNone, strToInt(node.token.lexeme))
	}
}

func strToInt(n string) int {
	v, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(v)
}

func enterDMem(node *SyntaxTree) int {
	location := findSym(node.token.lexeme)
	//dMem = append(dMem, location)
	dMem[currentDMemPos] = location
	oldPos := currentDMemPos
	currentDMemPos++
	return oldPos
}

func emitROCode(opcode int, dstReg, srcReg, srcReg2 int) {
	op := &Instruction{
		opcode: opcode,
	}
	op.regs = append(op.regs, dstReg)
	op.regs = append(op.regs, srcReg)
	op.regs = append(op.regs, srcReg2)
	iMem = append(iMem, op)
	Logf("emit rocode:%5s %-s, %-s, %-s\n", opTable[opcode], regTable[dstReg],
		regTable[srcReg], regTable[srcReg2])
}

func emitRMCode(opcode int, dstReg, srcReg int, offset int) {
	op := &Instruction{
		opcode: opcode,
	}
	iMem = append(iMem, op)
	Logf("emit rmcode:%5s %-s, %-d(%s)\n", opTable[opcode], regTable[dstReg],
		offset, regTable[srcReg])
}

func dumpRegister() {
	for i, v := range registers {
		Logf("reg:%s %d\n", regTable[i], v)
	}
}

func dumpInstructions() {
	for _, v := range iMem {
		if v == nil || v.opcode == opHalt {
			break
		}
		Logf("%s\n", opTable[v.opcode])
	}
}

func dumpDataMem() {

}
