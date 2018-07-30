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
var dMem = make([]int, dMemSize, dMemSize)
var currentDMemPos = 0

type regManager struct {
	top      int
	freeRegs map[int]bool
	usedRegs []int
}

var regM = regManager{
	top: -1,
	freeRegs: map[int]bool{
		r0: true,
		r1: true,
	},
}

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
		genIf(node)
	case repeatK:
		genRepeat(node)
	case assignK:
		genAssign(node)
	case readK:
		dstReg := allocReg()
		offset := findSym(node.child.token.lexeme)
		emitROCode(opIn, dstReg, regNone, regNone)
		emitRMCode(opSt, dstReg, r5, offset)
		freeReg(dstReg)
	case writeK:
		Logf("gen code for write\n")
		genExp(node.child)
		r := popUsedReg()
		emitROCode(opOut, r, regNone, regNone)
		freeReg(r)
	}
}

func genIf(node *SyntaxTree) {
	genExp(node.child)
}

func genRepeat(node *SyntaxTree) {

}

func genAssign(node *SyntaxTree) {
	offset := findSym(node.child.token.lexeme)
	for next := node.child.slibling; next != nil; next = next.slibling {
		genExp(next)
	}
	Logf("gen assign\n")
	lastUsed := popUsedReg()
	emitRMCode(opSt, lastUsed, r5, offset)
	freeReg(lastUsed)
}

func genExpForBinOp(node *SyntaxTree) {
	child := node.child
	genExp(node.child)
	for next := child.slibling; next != nil; next = next.slibling {
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
		srcReg := popUsedReg()
		dstReg := popUsedReg()
		emitROCode(opAdd, dstReg, srcReg, dstReg)
		freeReg(srcReg)
	case tokenMinus:
		child := node.child
		genExp(child)
		for next := child.slibling; next != nil; next = next.slibling {
			genExp(next)
		}
		srcReg := popUsedReg()
		dstReg := popUsedReg()
		emitROCode(opSub, dstReg, srcReg, dstReg)
	case tokenMultiply:
	case tokenDiv:
	case tokenLess:
		genExpForBinOp(node)
		emitROCode(opSub, r0, r1, r0)
	case tokenEqual:
		child := node.child
		genExp(child)
		for next := child.slibling; next != nil; next = next.slibling {
			genExp(next)
		}
		emitROCode(opSub, r0, r1, r0)

	case tokenId:
		dstReg := allocReg()
		offset := findSym(node.token.lexeme)
		emitRMCode(opLd, dstReg, r5, offset)
	case tokenNumber:
		dstReg := allocReg()
		emitRMCode(opLdc, dstReg, regNone, strToInt(node.token.lexeme))
	default:
		panic("oops\n")
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
	op.regs = append(op.regs, dstReg)
	op.regs = append(op.regs, srcReg)
	op.regs = append(op.regs, offset)

	iMem = append(iMem, op)
	Logf("emit rmcode:%5s %-s, %-d(%s)\n", opTable[opcode], regTable[dstReg],
		offset, regTable[srcReg])
}

func popUsedReg() int {
	top := regM.top
	if top < 0 {
		panic("no more used regs")
	}
	Logf("pop at top:%d used:%v\n", top, regM.usedRegs)
	r := regM.usedRegs[top]
	regM.top--
	return r
}

func freeReg(reg int) {
	free := regM.freeRegs[reg]
	if free {
		Logf("double free register:%s\n", regTable[reg])
		panic("double free register")
	}
	Logf("free register:%s top:%d\n", regTable[reg], regM.top)
	regM.freeRegs[reg] = true
	if len(regM.usedRegs) == 1 {
		regM.top = -1
		regM.usedRegs = append(regM.usedRegs[:0], regM.usedRegs[1:]...)
		return
	}
	for i, r := range regM.usedRegs {
		if r != reg {
			continue
		}
		regM.usedRegs = append(regM.usedRegs[i-1:i], regM.usedRegs[i+1:]...)
		regM.top = len(regM.usedRegs) - 1
		break
	}
	Logf("after free top:%d\n", regM.top)
}

func allocReg() int {
	minReg := regNone
	for r, free := range regM.freeRegs {
		if !free {
			continue
		}
		if r < minReg {
			minReg = r
		}
	}
	if minReg == regNone {
		panic("no free register\n")
	}
	Logf("alloc register:%s top:%d\n", regTable[minReg], regM.top)
	regM.freeRegs[minReg] = false
	regM.usedRegs = append(regM.usedRegs, minReg)
	regM.top++
	return minReg
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
	Logf("exit==================\n")
}

func dumpDataMem() {

}
