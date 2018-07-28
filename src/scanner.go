package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

const (
	KEYWORD = iota
	IF
	ELSE
	READ
	THEN
	REPEAT
	UNTIL
	WRITE
	END
	ID
	NUMBER
	ADD
	MINUS
	MULTIPLY
	DIV
	EQUAL
	LESS
	LPAREN
	RPAREN
	SEMI
	ASSIGN
	tokenEOF
)

const (
	tokenKeyWord  = KEYWORD
	tokenIf       = IF
	tokenElse     = ELSE
	tokenThen     = THEN
	tokenRead     = READ
	tokenWrite    = WRITE
	tokenRepeat   = REPEAT
	tokenUntil    = UNTIL
	tokenEndBlock = END
	tokenId       = ID
	tokenNumber   = NUMBER
	tokenAdd      = ADD
	tokenMinus    = MINUS
	tokenMultiply = MULTIPLY // *
	tokenDiv      = DIV      // \
	tokenEqual    = EQUAL    // =
	tokenLess     = LESS     // <
	tokenLParen   = LPAREN   // (
	tokenRParen   = RPAREN   // )
	tokenSemi     = SEMI     // ;
	tokenAssign   = ASSIGN   // :=
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
var internalTokenTable = make(map[int][]*tokenSymbol)
var tokenTable = map[int]string{
	KEYWORD:  "keyword",
	ID:       "id",
	NUMBER:   "number",
	ADD:      "add",
	MINUS:    "minus op",
	MULTIPLY: "mutiply",
	DIV:      "div",
	EQUAL:    "equal",
	LESS:     "less",
	LPAREN:   "leftparen(",
	RPAREN:   "rightparen)",
	SEMI:     "semicolon",
	ASSIGN:   "assign",
}

var token int

var keywordTable = map[string]int{
	"if":     IF,
	"read":   READ,
	"write":  WRITE,
	"until":  UNTIL,
	"then":   THEN,
	"else":   ELSE,
	"repeat": REPEAT,
	"end":    END,
}

var currentSrcFile = ""

func newToken() *tokenSymbol {
	t := &tokenSymbol{
		file:      currentSrcFile,
		line:      line,
		comment:   "",
		tokenType: token,
		lexeme:    lexeme,
	}
	tokenType, ok := keywordTable[lexeme]
	if ok {
		t.tokenType = tokenType
	}
	return t
}

func (t *tokenSymbol) SimpleStr() string {
	return fmt.Sprintf("{lexeme:%s}", t.lexeme)
}
func (t *tokenSymbol) String() string {
	return fmt.Sprintf("{f:%s l:%d type:%d lexeme:%s}",
		t.file, t.line, t.tokenType, t.lexeme)
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
		token = LPAREN
	case ')':
		token = RPAREN
	case ';':
		token = SEMI
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
		tokenLine, ok := internalTokenTable[line]
		if !ok {
			tokenLine = make([]*tokenSymbol, 0)
			internalTokenTable[line] = tokenLine
		}
		/*
			tokenLine = append(tokenLine, tokenSymbol{lexeme: lexeme, tokenType: tokenType,
				line: line})
		*/
		tokenLine = append(tokenLine, t)
		internalTokenTable[line] = tokenLine
	}
	Logf("line num:%d sym:%d\n", len(lines), len(internalTokenTable))
	for lineNum, pos := range lines {
		lexemes, ok := internalTokenTable[lineNum+1]
		if !ok {
			continue
		}
		begin := pos & 0xffff
		end := (pos >> 16) & 0xffff
		Logf("\n[line:%d %s]\n", lineNum+1, string(fileBuf[begin:end]))
		for _, l := range lexemes {
			Logf("line:%-4d token:[%-10s] \tlexeme:%-8v\n",
				l.line, tokenTable[l.tokenType], l.lexeme)
		}
	}
}
func dumpWithoutLine() {
	for t := GetToken(); t.tokenType != tokenEOF; t = GetToken() {
		Logf("line:%-4d token:[%-10s] \tlexeme:%-8v\n",
			t.line, tokenTable[t.tokenType], t.lexeme)
	}

}

func readFile(srcFile string) error {
	Logf("parse source file:%s\n", srcFile)
	buf, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}
	fileBuf = buf
	currentSrcFile = srcFile
	return nil
}
