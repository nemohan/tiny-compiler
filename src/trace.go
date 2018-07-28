package main

import (
	"fmt"
	"os"
)

var logFile = "log"
var logFileHandle *os.File

func initTrace() {
	file, err := os.Create(logFile)
	if err != nil {
		panic(err)
	}
	logFileHandle = file
}

func Logf(format string, args ...interface{}) {
	if !enableTrace {
		return
	}
	if logFileHandle == nil {
		fmt.Printf(format, args...)
		return
	}
	fmt.Fprintf(logFileHandle, format, args...)
}
