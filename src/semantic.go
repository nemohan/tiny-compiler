package main

import (
	"errors"
	"fmt"
)

const (
	typeVoid = iota //this does't exist
	typeInt
	typeBool
)

type typeErr struct {
	err          error
	expectedType int
	realType     int
	token        *tokenSymbol
}

var typeErrs = make([]*typeErr, 0)

func addErr(token *tokenSymbol, err error) {
	e := &typeErr{
		token: token,
		err:   err,
	}
	typeErrs = append(typeErrs, e)
}

/*
func checkType(node *SyntaxTree) {
	if node.token == nil {
		return
	}
	switch node.nodeKind {
	case stmtK:
		checkStmt(node)
	case expK:
	}
}
*/
func checkStmt(node *SyntaxTree) {
	//the statement's type does not decide yet
	if node.childs[0].expType == typeVoid {
		return
	}
	t := node.token
	childType := node.childs[0].expType
	switch t.tokenType {
	case tokenIf:
		if node.childs[0].expType != typeBool {
			err := errors.New("can't not use integer type in if statement, expect bool type")
			addErr(t, err)
		}

	case tokenRepeat:
		if node.childs[0].expType != typeBool {
			err := errors.New("can't not use integer type in repeat statement, expect bool type")
			addErr(t, err)
		}
	case tokenWrite:
		if childType != typeInt {
			err := errors.New("can't not use bool type in write statement, expect int type")
			addErr(t, err)
		}
	case tokenAssign:
		if childType != node.childs[1].expType {
			err := errors.New("can't not use bool type in assign statement, expect int type")
			addErr(t, err)
		}
	}
}

func checkWriteExp(node *SyntaxTree) {

}

func checkAssignment(node *SyntaxTree) {

}

func checkIfStmt(node *SyntaxTree) {
}

func checkOperand(node *SyntaxTree) {

}

func isCmpOperator(opType int) bool {
	operators := []int{
		tokenEqual,
		tokenLess,
	}
	for _, op := range operators {
		if opType == op {
			return true
		}
	}
	return false
}

func isAlgoOperator(opType int) bool {
	operators := []int{
		tokenAdd,
		tokenMinus,
		tokenMultiply,
		tokenDiv,
	}
	for _, op := range operators {
		if opType == op {
			return true
		}
	}
	return false
}

func checkType(node *SyntaxTree) {
	t := node.token
	if t == nil {
		return
	}
	//fmt.Printf("chedk stmt:%s kind:%d==============\n", t.lexeme, node.nodeKind)
	if node.nodeKind == stmtK {
		checkStmt(node)
		return
	}

	if t.tokenType == NUMBER || t.tokenType == ID {
		node.expType = typeInt
	} else if isAlgoOperator(t.tokenType) {
		node.expType = typeInt
	} else if isCmpOperator(t.tokenType) {
		node.expType = typeBool
	}
}

func isAnyTypeErr() bool {
	return len(typeErrs) != 0
}

func dumpTypeErr() {
	for _, t := range typeErrs {
		Logf("file:%s line:%d %v\n", t.token.file, t.token.line, t.err)
	}
}
