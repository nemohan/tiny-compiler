package main

import (
	"errors"
	"fmt"
)

var opcodeHandlerTable = map[int]func(*Instruction) (bool, error){
	opNone: emptyHandler,
	opHalt: haltHandler,
	opIn:   inHandler,
	opOut:  outHandler,
	opAdd:  addHandler,
	opSub:  emptyHandler,
	opMul:  emptyHandler,
	opDiv:  divHandler,
	opLd:   ldHandler,
	opLda:  ldaHandler,
	opLdc:  ldcHandler,
	opSt:   stHandler,
	opJlt:  emptyHandler,
	opJle:  emptyHandler,
	opJge:  emptyHandler,
	opJgt:  emptyHandler,
	opJEQ:  emptyHandler,
	opJNE:  emptyHandler,
}

//initEngine, init all registers, instruction memory and data memory
func initEngine() {
	for i, _ := range registers {
		registers[i] = 0
	}
	registers[r6] = dMemSize //mp
}

func processer() {
	for {
		pc := registers[regPC]
		op := iMem[pc]
		stop, err := execCode(op)
		if stop {
			break
		}
		if err != nil {
			Logf("processor error:%v\n", err)
			break
		}
		registers[regPC]++
	}
}

func execCode(op *Instruction) (bool, error) {
	code := op.opcode
	handler, ok := opcodeHandlerTable[code]
	if !ok {
		return true, fmt.Errorf("not find handler for code:%d %s", code, opTable[code])
	}
	return handler(op)
}

func haltHandler(op *Instruction) (bool, error) {

	return true, nil
}

func inHandler(op *Instruction) (bool, error) {
	return false, nil
}

func outHandler(op *Instruction) (bool, error) {
	return false, nil
}

func addHandler(op *Instruction) (bool, error) {
	return false, nil
}

func ldHandler(op *Instruction) (bool, error) {
	return false, nil
}

func ldaHandler(op *Instruction) (bool, error) {
	return false, nil
}

func ldcHandler(op *Instruction) (bool, error) {
	return false, nil
}

func stHandler(op *Instruction) (bool, error) {
	return false, nil
}

func divHandler(op *Instruction) (bool, error) {
	return false, nil
}

func emptyHandler(op *Instruction) (bool, error) {
	return true, errors.New("not support")
}
