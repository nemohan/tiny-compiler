package main

import (
	"errors"
	"fmt"
	"os"
)

var advancePC = true
var singleStep = false
var signalCh = make(chan int, 1)
var opcodeHandlerTable = map[int]func(*Instruction) (bool, error){
	opNone: emptyHandler,
	opHalt: haltHandler,
	opIn:   inHandler,
	opOut:  outHandler,
	opAdd:  addHandler,
	opSub:  subHandler,
	opMul:  emptyHandler,
	opDiv:  divHandler,
	opLd:   ldHandler,
	opLda:  ldaHandler,
	opLdc:  ldcHandler,
	opSt:   stHandler,
	opJlt:  jltHandler,
	opJle:  emptyHandler,
	opJge:  emptyHandler,
	opJgt:  emptyHandler,
	opJeq:  emptyHandler,
	opJne:  emptyHandler,
}

//initEngine, init all registers, instruction memory and data memory
func initEngine() {
	for i, _ := range registers {
		registers[i] = 0
	}
	registers[r6] = dMemSize //mp
}

func nextStep() {
	signalCh <- 1
}

func disableSingleStep() {
	singleStep = false
	nextStep()
}

func enableSingleStep() {
	singleStep = true
}

func processor() {
	for {
		pc := registers[regPC]
		Logf("processor, pc:%d begin\n", pc)
		op := iMem[pc]
		stop, err := execCode(op)
		if stop {
			Logf("processor stop. %v\n", err)
			break
		}
		if singleStep {
			<-signalCh
		}
		if err != nil {
			Logf("processor error:%v\n", err)
			break
		}
		if advancePC {
			registers[regPC]++
		}
		advancePC = true
		Logf("processor, pc:%d\n", registers[regPC])
	}
}

func isValidMemAddr() bool {

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
	r := op.regs[0]
	if !isValidRegister(r) {
		err := fmt.Errorf("invalid reg:%d in opcode <in> regs:%v", r, op.regs)
		return false, err
	}
	res := registers[r]
	if res < 0 {
		Logf("jlt jump to addr:%d\n", op.regs[2])
		registers[regPC] = op.regs[2]
		advancePC = false
	}
	//advancePC = false
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
	registers[reg] = in
	return false, nil
}

func outHandler(op *Instruction) (bool, error) {
	r := op.regs[0]
	if !isValidRegister(r) {
		return false, fmt.Errorf("invalid register in <out> regs:%v", op.regs)
	}
	fmt.Printf("out:%d\n", registers[r])
	Logf("exec <out> (%d,%d)\n", r, registers[r])
	return false, nil
}

func addHandler(op *Instruction) (bool, error) {
	for _, r := range op.regs {
		if !isValidRegister(r) {
			return false, fmt.Errorf("invalid register in <add> regs:%v", op.regs)
		}
	}
	dstReg := op.regs[0]
	old := registers[dstReg]
	srcReg := op.regs[1]
	registers[dstReg] = old + registers[srcReg]
	Logf(" exec <add> v1:%d v2:%d to reg:%d \n", old, registers[srcReg], dstReg)
	return false, nil
}

func ldHandler(op *Instruction) (bool, error) {
	if !isValidRegister(op.regs[0]) || !isValidRegister(op.regs[1]) {
		return false, fmt.Errorf("invalid register <ld> regs:%v", op.regs)
	}
	reg := op.regs[0]
	srcReg := op.regs[1]
	offset := op.regs[2]
	addr := registers[srcReg] + offset
	registers[reg] = dMem[addr]
	Logf(" exec <ld> load:%d at dmem:%d to reg:%d \n", dMem[addr], addr, reg)
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
	reg := op.regs[0]
	num := op.regs[2]
	registers[reg] = num
	Logf(" exec <ldc> load num:%d  reg:%d \n", num, reg)
	return false, nil
}

func stHandler(op *Instruction) (bool, error) {
	if !isValidRegister(op.regs[0]) || !isValidRegister(op.regs[1]) {
		return false, fmt.Errorf("invalid register in opcode:%s regs:%v", opTable[op.opcode],
			op.regs)
	}
	dstReg := op.regs[1]
	srcReg := op.regs[0]
	offset := op.regs[2]
	addr := registers[dstReg] + offset

	Logf("exec <st>  store:%d to:%d in dmem\n", registers[srcReg], addr)
	dMem[addr] = registers[srcReg]
	return false, nil
}

func divHandler(op *Instruction) (bool, error) {
	return false, nil
}

func subHandler(op *Instruction) (bool, error) {
	dstReg := op.regs[0]
	srcReg := op.regs[1]

	registers[dstReg] = registers[dstReg] - registers[srcReg]
	return false, nil
}

func emptyHandler(op *Instruction) (bool, error) {
	return true, errors.New("not support")
}
