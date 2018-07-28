package main

//initEngine, init all registers, instruction memory and data memory
func initEngine() {
	for i, _ := range registers {
		registers[i] = 0
	}
	r6 := dMemSize //mp
}

func processer() {
	for {
		pc := registers[regPc]
		op := iMem[pc]

		registers[regPc]++
	}
}
