package main

import (
	"errors"
	"fmt"
	"os"
)

//var advancePC = true

var opcodeHandlerTable = map[int]func(*Instruction) (bool, error){
	opNone: emptyHandler,
	opHalt: haltHandler,
	opIn:   inHandler,
	opOut:  outHandler,
	opAdd:  addHandler,
	opSub:  subHandler,
	opMul:  mulHandler,
	opDiv:  divHandler,
	opLd:   ldHandler,
	opLda:  ldaHandler,
	opLdc:  ldcHandler,
	opSt:   stHandler,
	opJlt:  jltHandler,
	opJle:  emptyHandler,
	opJge:  emptyHandler,
	opJgt:  emptyHandler,
	opJeq:  jeqHandler,
	opJne:  jneHandler,
	opUjp:  ujpHandler,
}

type TinyVM struct {
	imem           []*Instruction
	dmem           []int
	regs           []int
	regM           *regManager
	advancePC      bool
	singleStep     bool
	signalCh       chan int
	currentDMemPos int
}

var tvm = NewTinyVM()

func NewTinyVM() *TinyVM {
	return &TinyVM{
		imem:      make([]*Instruction, 0, iMemSize),
		dmem:      make([]int, dMemSize, dMemSize),
		regs:      make([]int, regNum),
		advancePC: true,
		regM: &regManager{
			top: -1,
			freeRegs: map[int]bool{
				r0: true,
				r1: true,
			},
		},
	}
}

func (tm *TinyVM) allocReg() int {
	regM := tm.regM
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

func (tm *TinyVM) freeReg(reg int) {
	regM := tm.regM
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
		Logf("free reg at:%d num:%d\n", i, len(regM.usedRegs))
		regM.usedRegs = append(regM.usedRegs[i-1:i], regM.usedRegs[i+1:]...)
		regM.top = len(regM.usedRegs) - 1
		break
	}
	Logf("after free top:%d\n", regM.top)
}

func (tm *TinyVM) popUsedReg() int {
	regM := tm.regM
	top := regM.top
	if top < 0 {
		panic("no more used regs")
	}
	Logf("pop at top:%d used:%v\n", top, regM.usedRegs)
	r := regM.usedRegs[top]
	regM.top--
	return r
}
func (tm *TinyVM) emitROCode(opcode int, dstReg, srcReg, srcReg2 int) int {
	imem := tm.imem
	op := &Instruction{
		opcode: opcode,
	}
	op.regs = append(op.regs, dstReg)
	op.regs = append(op.regs, srcReg)
	op.regs = append(op.regs, srcReg2)
	imem = append(imem, op)
	tm.imem = imem
	Logf("\t-------addr:%d emit rocode:%5s %-s, %-s, %-s\n", len(imem)-1, opTable[opcode],
		regTable[dstReg], regTable[srcReg], regTable[srcReg2])
	return len(imem) - 1
}

func (tm *TinyVM) emitRMCode(opcode int, dstReg, srcReg int, offset int) int {
	imem := tm.imem
	op := &Instruction{
		opcode: opcode,
	}
	op.regs = append(op.regs, dstReg)
	op.regs = append(op.regs, srcReg)
	op.regs = append(op.regs, offset)

	imem = append(imem, op)
	tm.imem = imem
	Logf("\t------addr:%d emit rmcode:%5s %-s, %-d(%s)\n", len(imem)-1, opTable[opcode],
		regTable[dstReg], offset, regTable[srcReg])
	return len(imem) - 1
}

func (tm *TinyVM) enterDMem(node *SyntaxTree) int {
	location := findSym(node.token.lexeme)
	tm.dmem[tm.currentDMemPos] = location
	oldPos := tm.currentDMemPos
	tm.currentDMemPos++
	return oldPos
}

func (tm *TinyVM) getIMemLocation() int {
	return len(tm.imem) - 1
}

func (tm *TinyVM) getIMemNextLoc() int {
	return len(tm.imem)
}

//func patchCode(srcPos, dstPos int) {
func (tm *TinyVM) patchCode(pos int) {
	Logf("want patch some code at pos:%d\n", pos)
	if pos == invalidPos {
		pos = len(iMemPatches)
		if pos == 0 {
			return
		}
		pos -= 1
	}
	srcPos := iMemPatches[pos]
	dstPos := getIMemNextLoc()
	op := tm.imem[srcPos]
	op.regs[2] = dstPos
	Logf("add patch src pos:%d dst pos:%d\n", srcPos, dstPos)
	iMemPatches = append(iMemPatches[:pos], iMemPatches[pos+1:]...)
}

func (tm *TinyVM) DumpRegister() {
	for i, v := range tm.regs {
		Logf("reg:%s %d\n", regTable[i], v)
	}
}

func (tm *TinyVM) dumpInstructions() {
	for i, v := range tm.imem {
		if v == nil {
			break
		}
		if v.opcode == opHalt {
			Logf("%04d: %-6s %s, %d(%s)\n", i, opTable[v.opcode], regTable[v.regs[0]],
				v.regs[2], regTable[v.regs[1]])
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

//initEngine, init all registers, instruction memory and data memory
func initEngine() {
	/*
		for i, _ := range registers {
			registers[i] = 0
		}
		registers[r6] = dMemSize //mp
	*/
}

func (tm *TinyVM) nextStep() {
	tm.signalCh <- 1
}

func (tm *TinyVM) disableSingleStep() {
	tm.singleStep = false
	tm.nextStep()
}

func (tm *TinyVM) enableSingleStep() {
	tm.singleStep = true
}

func (tm *TinyVM) processor() {
	for {
		pc := tm.regs[regPC]
		Logf("processor, pc:%d begin\n", pc)
		op := tm.imem[pc]
		op.vm = tm
		stop, err := execCode(op)
		if stop {
			Logf("processor stop. %v\n", err)
			break
		}
		if tm.singleStep {
			<-tm.signalCh
		}
		if err != nil {
			Logf("processor error:%v\n", err)
			break
		}
		if tm.advancePC {
			tm.regs[regPC]++
		}
		tm.advancePC = true
		Logf("processor, pc:%d\n", tm.regs[regPC])
	}
}

func isValidMemAddr() bool {

	return true
}

func checkRegister(op *Instruction) bool {
	for i, r := range op.regs {
		if !isValidRegister(r) {
			Logf("invalid register at:%d op:%s regs:%v\n", i, opTable[op.opcode], op.regs)
			return false
		}
	}
	return true
}
func isValidRegister(reg int) bool {
	if reg < r0 || reg >= regNone {
		return false
	}
	return true
}

func execCode(op *Instruction) (bool, error) {
	code := op.opcode
	handler, ok := opcodeHandlerTable[code]
	if !ok {
		return true, fmt.Errorf("not find handler for code:%d %s", code, opTable[code])
	}
	Logf("prepare code:%s\n", opTable[code])
	return handler(op)
}

func haltHandler(op *Instruction) (bool, error) {
	return true, nil
}

func jltHandler(op *Instruction) (bool, error) {
	//NOTE: absolute address
	vm := op.vm
	r := op.regs[0]
	if !isValidRegister(r) {
		err := fmt.Errorf("invalid reg:%d in opcode <in> regs:%v", r, op.regs)
		return false, err
	}
	res := vm.regs[r]
	if res < 0 {
		Logf("jlt res:%d jump to addr:%d\n", res, op.regs[2])
		vm.regs[regPC] = op.regs[2]
		vm.advancePC = false
	}
	return false, nil
}

func ujpHandler(op *Instruction) (bool, error) {
	vm := op.vm
	currentPos := vm.regs[regPC]
	vm.regs[regPC] = op.regs[2]
	Logf("ujp from location:%d to:%d\n", currentPos, op.regs[2])
	vm.advancePC = false
	return false, nil
}

func jeqHandler(op *Instruction) (bool, error) {
	vm := op.vm
	if vm.regs[op.regs[0]] != 0 {
		return false, nil
	}
	vm.regs[regPC] = op.regs[2]
	oldPos := vm.regs[regPC]
	vm.advancePC = false
	Logf("jeq from locaton:%d to:%d\n", oldPos, op.regs[2])
	return false, nil
}

func jneHandler(op *Instruction) (bool, error) {
	vm := op.vm
	if vm.regs[op.regs[0]] == 0 {
		return false, nil
	}
	oldPos := vm.regs[regPC]
	vm.regs[regPC] = op.regs[2]
	vm.advancePC = false
	Logf("jne from locaton:%d to:%d\n", oldPos, op.regs[2])
	return false, nil
}

func inHandler(op *Instruction) (bool, error) {
	reg := op.regs[0]
	if !isValidRegister(reg) {
		err := fmt.Errorf("invalid reg:%d in opcode <in> regs:%v", reg, op.regs)
		return false, err
	}
	in := 0
	fmt.Printf("input integer:\n")
	n, err := fmt.Fscanf(os.Stdin, "%d", &in)
	if err != nil {
		return false, fmt.Errorf("invalid input in opcode <in>. num:%d err:%v", n, err)
	}

	Logf("input :%d\n", in)
	op.vm.regs[reg] = in
	return false, nil
}

func outHandler(op *Instruction) (bool, error) {
	vm := op.vm
	r := op.regs[0]
	if !isValidRegister(r) {
		return false, fmt.Errorf("invalid register in <out> regs:%v", op.regs)
	}
	fmt.Printf("out:%d\n", vm.regs[r])
	Logf("exec <out> (%d,%d)\n", r, vm.regs[r])
	return false, nil
}

func ldHandler(op *Instruction) (bool, error) {
	if !isValidRegister(op.regs[0]) || !isValidRegister(op.regs[1]) {
		return false, fmt.Errorf("invalid register <ld> regs:%v", op.regs)
	}
	vm := op.vm
	reg := op.regs[0]
	srcReg := op.regs[1]
	offset := op.regs[2]
	addr := vm.regs[srcReg] + offset
	vm.regs[reg] = vm.dmem[addr]
	Logf(" exec <ld> load:%d at dmem:%d to reg:%d \n", vm.dmem[addr], addr, reg)
	return false, nil
}

func ldaHandler(op *Instruction) (bool, error) {
	return false, nil
}

func ldcHandler(op *Instruction) (bool, error) {
	if !isValidRegister(op.regs[0]) {
		return false, fmt.Errorf("invalid register in opcode:%s regs:%v", opTable[op.opcode],
			op.regs)
	}
	vm := op.vm
	reg := op.regs[0]
	num := op.regs[2]
	vm.regs[reg] = num
	Logf(" exec <ldc> load num:%d  reg:%d \n", num, reg)
	return false, nil
}

func stHandler(op *Instruction) (bool, error) {
	if !isValidRegister(op.regs[0]) || !isValidRegister(op.regs[1]) {
		return false, fmt.Errorf("invalid register in opcode:%s regs:%v", opTable[op.opcode],
			op.regs)
	}
	vm := op.vm
	dstReg := op.regs[1]
	srcReg := op.regs[0]
	offset := op.regs[2]
	addr := vm.regs[dstReg] + offset

	Logf("exec <st>  store:%d to:%d in dmem\n", vm.regs[srcReg], addr)
	vm.dmem[addr] = vm.regs[srcReg]
	return false, nil
}

//===================================arithmetic operation
func mulHandler(op *Instruction) (bool, error) {
	if !checkRegister(op) {
		return true, nil
	}
	vm := op.vm
	dstReg := op.regs[0]
	srcReg := op.regs[1]
	vm.regs[dstReg] = vm.regs[dstReg] * vm.regs[srcReg]
	return false, nil

}

func divHandler(op *Instruction) (bool, error) {
	//return false, nil
	return emptyHandler(op)
}

func subHandler(op *Instruction) (bool, error) {
	registers := op.vm.regs
	dstReg := op.regs[0]
	srcReg := op.regs[1]
	old := registers[dstReg]
	registers[dstReg] = registers[srcReg] - registers[dstReg]
	Logf("sub %s - %s res:%d = %d - %d\n", regTable[dstReg], regTable[srcReg],
		registers[dstReg], registers[srcReg], old)
	return false, nil
}

func addHandler(op *Instruction) (bool, error) {
	for _, r := range op.regs {
		if !isValidRegister(r) {
			return false, fmt.Errorf("invalid register in <add> regs:%v", op.regs)
		}
	}
	registers := op.vm.regs
	dstReg := op.regs[0]
	old := registers[dstReg]
	srcReg := op.regs[1]
	registers[dstReg] = old + registers[srcReg]
	Logf(" exec <add> v1:%d v2:%d to reg:%d \n", old, registers[srcReg], dstReg)
	return false, nil
}
func emptyHandler(op *Instruction) (bool, error) {
	return true, errors.New("not support")
}
