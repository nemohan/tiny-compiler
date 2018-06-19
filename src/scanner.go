package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

const (
	KEYWORD = iota
	ID
	NUMBER
	ADD
	MINUS
	MULTIPLY
	DIV
	EQUAL
	LESS
	LEFTPAREN
	RIGHTPAREN
	SEMICOLON
	ASSIGN
)

const (
	stateStart = iota
	stateComment
	stateNumber
	stateId
	stateAssign
	stateOperator
	stateOther
	stateDone
	stateErr
)

type tokenSymbol struct {
	file    string
	line    int
	comment string
}

var fileBuf []byte
var currentState = stateStart
var tokenBegin int
var tokenEnd int
var next int
var symbolTable map[int]string
var line = 1
var rdPos int
var lexeme string
var invalidToken = -1
var operatorTable = map[byte]int{'+': ADD, '-': MINUS, '*': MULTIPLY,
	'/': DIV, '=': EQUAL, '<': LESS}

var lineBuf = bytes.NewBuffer([]byte(""))
var tokenTable = map[int]string{
	KEYWORD:    "keyword",
	ID:         "id",
	NUMBER:     "number",
	ADD:        "add",
	MINUS:      "minus op",
	MULTIPLY:   "mutiply",
	DIV:        "div",
	EQUAL:      "equal",
	LESS:       "less",
	LEFTPAREN:  "leftparen(",
	RIGHTPAREN: "rightparen)",
	SEMICOLON:  "semicolon",
	ASSIGN:     "assign",
}
var token int

func init() {
	buf, err := ioutil.ReadFile("test.ty")
	if err != nil {
		panic(err)
	}
	fileBuf = buf
}

func GetToken() int {
	size := len(fileBuf)
	for i := rdPos; i < size; {
		rdPos = i
		c := fileBuf[i]
		switch currentState {
		case stateStart:
			currentState = handleStart(c)
		case stateComment:
			currentState = handleComment(c)
		case stateNumber:
			currentState = handleNumber(c)
		case stateId:
			currentState = handleId(c)
		case stateAssign:
			currentState = handleAssign(c)
		case stateOperator:
			currentState = handleOperator(c)
		case stateOther:
			currentState = handleOther(c)
		case stateDone:
			currentState = stateStart
			return token
		default:
		}
		i = rdPos
		i++
	}
	return invalidToken
}

func handleOperator(c byte) int {
	tokenBegin = rdPos
	lexeme = getLexeme(rdPos + 1)
	opToken, ok := operatorTable[c]
	if !ok {
		return stateErr
	}
	token = opToken
	return stateDone
}

func isOperator(c byte) bool {
	if _, ok := operatorTable[c]; ok {
		return true
	}
	return false
}

func handleOther(c byte) int {
	switch c {
	case '(':
	case ')':
	case ';':
	default:
	}
	tokenBegin = rdPos
	lexeme = getLexeme(rdPos + 1)
	return stateDone
}

//return the next state
func handleStart(c byte) int {
	switch c {
	case ':':
		putBack()
		return stateAssign
	case '{':
		return stateComment
	case '(':
		fallthrough
	case ')':
		fallthrough
	case ';':
		putBack()
		return stateOther
	case ' ':
		return stateStart
	case '\t':
		return stateStart
	case '\n':
		line++
		return stateStart
	default:
		if isDigit(c) {
			tokenBegin = rdPos
			return stateNumber
		}
		if isCharacter(c) {
			tokenBegin = rdPos
			return stateId
		}
		if isOperator(c) {
			putBack()
			return stateOperator
		}
	}
	return stateDone
}

func handleComment(c byte) int {
	switch c {
	case '}':
		return stateStart
	case '\n':
		line++
	default:
	}
	return stateComment
}

func handleNumber(c byte) int {
	switch c {
	case ' ':
		fallthrough
	case ';':
		fallthrough
	case '\t':
		fallthrough
	case '\n':
		fallthrough
	case ')':
		lexeme = getLexeme(rdPos)
		putBack()
		return stateDone
	default:
		if isOperator(c) {
			lexeme = getLexeme(rdPos)
			putBack()
			return stateDone
		}
		if !isDigit(c) {
			return stateErr
		}
	}
	return stateDone
}

func handleId(c byte) int {
	switch c {
	case '(':
		fallthrough
	case ')':
		fallthrough
	case ':':
		lexeme = string(fileBuf[tokenBegin:rdPos])
		token = ID
		rdPos--
		return stateDone
	case '\n':
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		fallthrough
	case ';':
		lexeme = getLexeme(rdPos)
		token = ID
		putBack()
		return stateDone
	default:
		if ok := isOperator(c); ok {
			lexeme = string(fileBuf[tokenBegin:rdPos])
			rdPos--
			token = ID
			return stateDone
		}
		if !isCharacter(c) {
			return stateErr
		}
	}
	return stateId
}

func getLexeme(end int) string {
	return string(fileBuf[tokenBegin:end])
}

func handleAssign(c byte) int {
	switch c {
	case ':':
		tokenBegin = rdPos
		return stateAssign
	case '=':
		lexeme = getLexeme(rdPos + 1)
		token = ASSIGN
		return stateDone
	default:
	}
	return stateErr
}

func isDigit(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func putBack() {
	rdPos--
}

func isCharacter(c byte) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	if c >= 'A' && c <= 'Z' {
		return true
	}
	return false
}

func main() {
	//lastLine := line
	for tokenType := GetToken(); tokenType != invalidToken; {
		fmt.Printf("line:%-4d token:[%-10s] \tlexeme:%-8v\n", line, tokenTable[tokenType], lexeme)
		tokenType = GetToken()
	}
}
