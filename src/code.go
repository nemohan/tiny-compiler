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
	regPc
	regNone
)

var regTable = map[int]string{
	r0:    "r0",
	r1:    "r1",
	r2:    "r2",
	r3:    "r3",
	r4:    "r4",
	r5:    "r5",
	r6:    "r6",
	regPC: "pc",
}

type Instruction struct {
	opcode int
	regs   []int
}

const regNum = 8
const iMemSize = 1024
const dMemSize = 1024

var registers = make([]int, regNum)
var iMem = make([]*Instruction, iMemSize)
var dMem = make([]int, dMemSize)
var currentDMemPos = 0

func genCode(root *SyntaxTree) {
	node := root
	if root.nodeKind == fileK {
		node = root.child
	}
	if node.nodeKind == stmtK {
		genStmt(node)
	}
	if node.nodeKind == expK {
		genExp(k)
	}
}

func genStmt(node *SyntaxTree) {
	switch node.stmtKind {
	case ifK:
	case repeatK:
	case assignK:
		genAssign(node)
	case readK:
		// read a;
		// lda r0, offset(r5)
		// in r1
		// st r1, offset(r5)
		offset := enterDMem(node.child)
		emitROCode(opIn, r1, regNone, regNone)
		emitRMCode(opSt, r1, r5, offset)
	case writeK:
	}
}

func genAssign(node *SyntaxTree) {
	offset := enterDMem(node.child)
	// lda r0, offfset(r5)
	//emitRMCode(opLda, r0, r5, offset)
	for next := node.child.slibling; next != nil; next = next.slibling {
		genExp(next)
	}
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
		emitROCode(op, opAdd, r0, r1)
		/*
			child := node.child
			offset := findSym(child.token.lexeme)
			emitRMCode(opLd, r0, r5, offset)
			genExp(child.slibling)
			emitROCode(opAdd, r0, r0, r1)
		*/
	case tokenMinus:
	case tokenMutiply:
	case tokenDiv:
	case tokenId:
		offset := findSym(child.token.lexeme)
		emitRMCode(opLd, r0, r5, offset)
	case tokenNumber:
		emitRMCode(opLdc, r1, strToInt(node.token.lexeme), regNone)
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
}

func emitRMCode(opcode int, dstReg, srcReg int, offset int) {
	op := &Instruction{
		opcode: opcode,
	}
	iMem = append(iMem, op)
}

func dumpRegister() {
	for i, v := range registers {
		fmt.Printf("reg:%s %d\n", regTable[i], v)
	}
}

func dumpInstrucions() {
	for i, v := range iMem {
		fmt.Printf("%s\n", opTable[v.opcode])
	}
}

func dumpDataMem() {

}
