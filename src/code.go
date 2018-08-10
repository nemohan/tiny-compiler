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
	opJeq
	opJne
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
	opJeq:  "jeq",
	opJne:  "jne",
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

var iMemPatches = make([]int, 0)

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

func isRMCode(code int) bool {
	rmCodes := map[int]bool{
		opLd:  true,
		opLda: true,
		opLdc: true,
		opSt:  true,
		opJlt: true,
		opJle: true,
		opJge: true,
		opJgt: true,
		opJeq: true,
		opJne: true,
	}
	if _, ok := rmCodes[code]; ok {
		return true
	}
	return false
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
		//node = root.child
		node = root.childs[0]
	}
	if node.nodeKind == stmtK {
		genStmt(node)
	}
	if node.nodeKind == expK {
		Logf("gen exp\n")
		genExp(node)
	}
	for next := node.sibling; next != nil; next = next.sibling {
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
		offset := findSym(node.Left().token.lexeme)
		emitROCode(opIn, dstReg, regNone, regNone)
		emitRMCode(opSt, dstReg, r5, offset)
		freeReg(dstReg)
	case writeK:
		Logf("gen code for write\n")
		genExp(node.Left())
		r := popUsedReg()
		emitROCode(opOut, r, regNone, regNone)
		freeReg(r)
	}
}

func genIf(node *SyntaxTree) {
	Logf("gen if\n")
	patchCode()
	left := node.Left()
	genExp(left)
	if left.token.tokenType == tokenLess {
		destReg := popUsedReg()
		pos := emitRMCode(opJlt, destReg, regNone, regNone)
		//note: we should free destReg now
		freeReg(destReg)
		iMemPatches = append(iMemPatches, pos)
		Logf("need patch at imem location:%d\n", pos)
	}

	//TODO: if else  and if else if
	//TODO: use child and slibling can't distinguish the body part of if and the else part
	/*
		for next := node.child.slibling; next != nil; next = next.slibling {
			genStmt(next)
			if next.token.tokenType == tokenElse {
				patchCode()
			}
		}
	*/
	//note: just for then else
	//genStmt(node.Right())
	genCode(node.Right())
	patchCode()
}

func genRepeat(node *SyntaxTree) {

}

func genAssign(node *SyntaxTree) {
	//offset := findSym(node.child.token.lexeme)
	left := node.Left()
	offset := findSym(left.token.lexeme)
	/*
		for next := node.child.slibling; next != nil; next = next.slibling {
			genExp(next)
		}
	*/
	right := node.Right()
	genExp(right)
	Logf("gen assign\n")
	lastUsed := popUsedReg()
	emitRMCode(opSt, lastUsed, r5, offset)
	freeReg(lastUsed)
}

func genExpForBinOp(node *SyntaxTree) {
	/*
		child := node.child
		genExp(node.child)
		for next := child.slibling; next != nil; next = next.slibling {
			genExp(next)
		}
	*/
	genExp(node.Left())
	genExp(node.Right())
}

func genExp(node *SyntaxTree) {
	if node == nil {
		return
	}
	switch node.token.tokenType {
	case tokenAdd:
		//child := node.childs[0]
		genExp(node.Left())
		/*
			for next := child.slibling; next != nil; next = next.slibling {
				genExp(next)
			}
		*/
		genExp(node.Right())
		srcReg := popUsedReg()
		dstReg := popUsedReg()
		emitROCode(opAdd, dstReg, srcReg, dstReg)
		freeReg(srcReg)
	case tokenMinus:
		/*
			child := node.child
			genExp(child)
			for next := child.slibling; next != nil; next = next.slibling {
				genExp(next)
			}
		*/
		genExp(node.Left())
		genExp(node.Right())
		srcReg := popUsedReg()
		dstReg := popUsedReg()
		emitROCode(opSub, dstReg, srcReg, dstReg)
	case tokenMultiply:
	case tokenDiv:
	case tokenLess:
		genExpForBinOp(node)
		srcReg := popUsedReg() //left
		dstReg := popUsedReg() //right
		// dstReg = srcReg -dstReg
		emitROCode(opSub, dstReg, srcReg, dstReg)
		freeReg(srcReg)
	case tokenEqual:
		/*
			child := node.child
			genExp(child)
			for next := child.slibling; next != nil; next = next.slibling {
				genExp(next)
			}
		*/
		genExp(node.Left())
		genExp(node.Right())
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

/*
func nextSibling(node *SyntaxTree) func() *SyntaxTree {
	return func() *SyntaxTree {
		next := node.slibling
		if next != nil {
			node = next.slibling
		}
		return next
	}
}
*/

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

func getIMemLocation() int {
	return len(iMem) - 1
}

func getIMemNextLoc() int {
	return len(iMem)
}

//func patchCode(srcPos, dstPos int) {
func patchCode() {
	Logf("want patch some code\n")
	pos := len(iMemPatches)
	if pos == 0 {
		return
	}
	srcPos := iMemPatches[pos-1]
	dstPos := getIMemNextLoc()
	op := iMem[srcPos]
	op.regs[2] = dstPos
	Logf("add patch src pos:%d dst pos:%d\n", srcPos, dstPos)
	iMemPatches = append(iMemPatches[:pos-1], iMemPatches[pos:]...)
}

func emitROCode(opcode int, dstReg, srcReg, srcReg2 int) int {
	op := &Instruction{
		opcode: opcode,
	}
	op.regs = append(op.regs, dstReg)
	op.regs = append(op.regs, srcReg)
	op.regs = append(op.regs, srcReg2)
	iMem = append(iMem, op)
	Logf("\t-------addr:%d emit rocode:%5s %-s, %-s, %-s\n", len(iMem)-1, opTable[opcode],
		regTable[dstReg], regTable[srcReg], regTable[srcReg2])
	return len(iMem) - 1
}

func emitRMCode(opcode int, dstReg, srcReg int, offset int) int {
	op := &Instruction{
		opcode: opcode,
	}
	op.regs = append(op.regs, dstReg)
	op.regs = append(op.regs, srcReg)
	op.regs = append(op.regs, offset)

	iMem = append(iMem, op)
	Logf("\t------addr:%d emit rmcode:%5s %-s, %-d(%s)\n", len(iMem)-1, opTable[opcode],
		regTable[dstReg], offset, regTable[srcReg])
	return len(iMem) - 1
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
	for i, v := range iMem {
		if v == nil || v.opcode == opHalt {
			break
		}

		if isRMCode(v.opcode) {
			Logf("%04d: %-6s %s, %d(%s)\n", i, opTable[v.opcode], regTable[v.regs[0]],
				v.regs[2], regTable[v.regs[1]])
		} else {
			Logf("%04d: %-6s %s, %s, %s\n", i, opTable[v.opcode], regTable[v.regs[0]],
				regTable[v.regs[1]], regTable[v.regs[2]])
		}

	}
	Logf("exit==================\n")
}

func dumpDataMem() {

}
