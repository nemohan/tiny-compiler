package main

import (
	"bytes"
	"flag"
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
	tokenEOF
)

const (
	tokenKeyWord    = KEYWORD
	tokenId         = ID
	tokenNumber     = NUMBER
	tokenAdd        = ADD
	tokenMinus      = MINUS
	tokenMultiply   = MULTIPLY   // *
	tokenDiv        = DIV        // \
	tokenEqual      = EQUAL      // =
	tokenLess       = LESS       // <
	tokenLeftParen  = LEFTPAREN  // (
	tokenRightParen = RIGHTPAREN // )
	tokenSemicolon  = SEMICOLON  // ;
	tokenAssign     = ASSIGN     // :=
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
	file      string
	line      int
	comment   string
	tokenType int
	lexeme    string
}

var fileBuf []byte
var currentState = stateStart
var tokenBegin int
var tokenEnd int
var next int

//var symbolTable map[int]string
var line = 1
var rdPos int
var lexeme string
var invalidToken = -1
var lines = make([]uint32, 0)
var operatorTable = map[byte]int{'+': ADD, '-': MINUS, '*': MULTIPLY,
	'/': DIV, '=': EQUAL, '<': LESS}

var lineBuf = bytes.NewBuffer([]byte(""))
var lineBegin = 0
var lineEnd = 0
var symbolTable = make(map[int][]*tokenSymbol)
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

var keywordTable = []string{
	"if",
	"read",
	"write",
	"until",
	"then",
	"else",
	"repeat",
	"end",
}

var testFile1 = "if.ty"
var testFile2 = "test.ty"

func init() {
	/*
		buf, err := ioutil.ReadFile(testFile1)
		if err != nil {
			panic(err)
		}
		fileBuf = buf
		fmt.Printf("%c\n", fileBuf[len(buf)-1])
	*/
}

func newToken() *tokenSymbol {
	t := &tokenSymbol{
		file:      "test.ty",
		line:      line,
		comment:   "",
		tokenType: token,
		lexeme:    lexeme,
	}
	for _, k := range keywordTable {
		if lexeme == k {
			t.tokenType = KEYWORD
			break
		}
	}
	return t
}

func GetToken() *tokenSymbol {
	lexeme = ""
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
			return newToken()
		default:
		}
		i = rdPos
		i++
	}
	token = tokenEOF
	return newToken()
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
		token = LEFTPAREN
	case ')':
		token = RIGHTPAREN
	case ';':
		token = SEMICOLON
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
		begin := uint32(lineBegin) & 0xffff
		end := (uint32(rdPos)) & 0xffff
		lines = append(lines, begin|(end<<16))
		lineBegin = rdPos + 1
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
		begin := uint32(lineBegin) & 0xffff
		end := (uint32(rdPos) + 1) & 0xffff
		lines = append(lines, begin|(end<<16))
		lineBegin = rdPos + 1
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
		token = NUMBER
		return stateDone
	default:
		if isOperator(c) {
			lexeme = getLexeme(rdPos)
			putBack()
			token = NUMBER
			return stateDone
		}
		//TODO: should continue
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

func dumpWithLine() {
	for t := GetToken(); t.tokenType != tokenEOF; t = GetToken() {
		tokenLine, ok := symbolTable[line]
		if !ok {
			tokenLine = make([]*tokenSymbol, 0)
			symbolTable[line] = tokenLine
		}
		/*
			tokenLine = append(tokenLine, tokenSymbol{lexeme: lexeme, tokenType: tokenType,
				line: line})
		*/
		tokenLine = append(tokenLine, t)
		symbolTable[line] = tokenLine
	}
	fmt.Printf("line num:%d sym:%d\n", len(lines), len(symbolTable))
	for lineNum, pos := range lines {
		lexemes, ok := symbolTable[lineNum+1]
		if !ok {
			continue
		}
		begin := pos & 0xffff
		end := (pos >> 16) & 0xffff
		fmt.Printf("\n[line:%d %s]\n", lineNum+1, string(fileBuf[begin:end]))
		for _, l := range lexemes {
			fmt.Printf("line:%-4d token:[%-10s] \tlexeme:%-8v\n",
				l.line, tokenTable[l.tokenType], l.lexeme)
		}
	}
}
func dumpWithoutLine() {
	for t := GetToken(); t.tokenType != tokenEOF; t = GetToken() {
		fmt.Printf("line:%-4d token:[%-10s] \tlexeme:%-8v\n",
			t.line, tokenTable[t.tokenType], t.lexeme)
	}

}

var sourceFile = ""

func readFile() {
	fmt.Printf("parse source file:%s\n", sourceFile)
	buf, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		panic(err)
	}
	fileBuf = buf
}
func main() {
	flag.StringVar(&sourceFile, "f", "", "tiny source file")
	flag.Parse()
	if sourceFile == "" {
		fmt.Printf("no source file\n usage: scanner -f source\n")
		return
	}
	//dumpWithoutLine()
	//dumpWithLine()
	readFile()
	Parse()
}
