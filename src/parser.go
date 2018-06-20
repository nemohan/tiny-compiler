package main

import (
	"errors"
	"fmt"
)

const (
	ifStmt = iota
	assignStmt
	repeatStmt
	readStmt
	writeStmt
	algoExp
	relationExp
)

//NodeKind
const (
	stmtK = iota
	expK
	fileK
)

//statement kind
const (
	ifK = iota
	repeatK
	assignK
	readK
	writeK
)

// exp kind
const (
	opK = iota
	constK
	idK
)

//TODO: use a struct may be better
//var currentLexeme int
var currentToken *tokenSymbol
var parseErr error

func Parse() {
	currentToken = GetToken()
	fmt.Printf("current token:%v\n", currentToken)
	astRoot := NewSyntaxTree(nil, fileK)
	parseStmtSequence(astRoot)
	astRoot.Traverse()
}

func parseStmtSequence(root *SyntaxTree) *SyntaxTree {
	var tree *SyntaxTree
	for parseErr == nil && currentToken.tokenType != tokenEOF {
		switch currentToken.tokenType {
		case tokenId:
			left := NewSyntaxTree(currentToken, expK)
			match(tokenId)
			if !match(tokenAssign) {
				parseErr = errors.New(" tokenAssign")
				return tree
			}
			assignTree := NewSyntaxTree(currentToken, expK)
			assignTree.AddLeftChild(left)
			right := parseAssignStmt()
			assignTree.AddRightChild(right)
			root.AddChild(assignTree)
			//don't care match or not
			matchStr(";")
		default:
			lastToken := currentToken
			//TODO: error
			if matchStr("if") {
				ifTree := NewSyntaxTree(lastToken, stmtK)
				fmt.Printf("match if\n")
				ifTree.AddLeftChild(parseIfStmt())
				root.AddChild(ifTree)
				matchStr(";")
				continue
			}

			if matchStr("read") {
				readTree := NewSyntaxTree(lastToken, stmtK)
				readTree.AddLeftChild(parseReadStmt())
				root.AddChild(readTree)
				matchStr(";")
				fmt.Printf("parse read done\n")
				continue
			}
			if matchStr("repeat") {
				parseRepeatStmt()
				continue

			}
			if matchStr("write") {
				parseWriteStmt()
				continue
			}
			return tree
		}
	}
	return tree
}

func parseStmt() {

}

func parseIfStmt() *SyntaxTree {
	tree := parseExp()
	lastToken := currentToken
	if !matchStr("then") {
		//TODO:
	}
	slibling := NewSyntaxTree(lastToken, expK)
	tree.AddSlibling(slibling)
	fmt.Printf("match then\n")
	parseStmtSequence(slibling)
	//slibling.AddLeftChild(child)

	//optional
	if matchStr("else") {
		parseStmtSequence(nil)
	}

	if !matchStr("end") {

	}
	fmt.Printf("match end\n")
	return tree
}

func parseRepeatStmt() {
	parseStmtSequence(nil)
	matchStr("until")
	parseExp()
}

func parseAssignStmt() *SyntaxTree {
	fmt.Printf("parse assignStmt\n")
	tree := parseExp()
	fmt.Printf("parse assignStmt end:%v\n", currentToken)
	return tree
}

func parseReadStmt() *SyntaxTree {
	switch currentToken.tokenType {
	case tokenId:
		tree := NewSyntaxTree(currentToken, expK)
		fmt.Printf("parseReadStmt:%v\n", currentToken)
		match(tokenId)
		return tree
	default:
		//TODO: error
	}
	return nil
}

func parseWriteStmt() {
	parseExp()
}

func handleExp(tree *SyntaxTree) *SyntaxTree {
	right := parseSimpleExp()
	//tree.AddRightChild(right)
	tree.AddChild(right)
	return tree
}

func parseExp() *SyntaxTree {
	expTree := parseSimpleExp()
	switch currentToken.tokenType {
	case tokenLess:
		fmt.Printf("match token less\n")
		node := NewSyntaxTree(currentToken, expK)
		node.AddLeftChild(expTree)
		match(tokenLess)
		expTree = handleExp(node)
		fmt.Printf("==================\n")
		expTree.Traverse()

	case tokenEqual:
		match(tokenEqual)
		expTree = handleExp(expTree)
	default:
		//panic("parseExp")
	}
	return expTree
}

func parseSimpleExp() *SyntaxTree {
	leftTree := parseTerm()
	tokenType := currentToken.tokenType
	for tokenType == tokenAdd || tokenType == tokenMinus {
		switch tokenAdd {
		case tokenAdd:
			tree := NewSyntaxTree(currentToken, expK)
			tree.AddLeftChild(leftTree)
			match(tokenAdd)
			tree.AddRightChild(parseSimpleExp())
			leftTree = tree
		case tokenMinus:
			tree := NewSyntaxTree(currentToken, expK)
			match(tokenMinus)
			tree.AddRightChild(parseSimpleExp())
			leftTree = tree
		default:
			panic("parseSimpleExp")
		}
	}
	return leftTree
}

func parseTerm() *SyntaxTree {
	leftTree := parseFactor()
	tokenType := currentToken.tokenType
	for tokenType == tokenMultiply || tokenType == tokenDiv {
		switch tokenType {
		case tokenMultiply:
			opTree := NewSyntaxTree(currentToken, expK)
			match(tokenMultiply)
			opTree.AddLeftChild(leftTree)
			opTree.AddRightChild(parseFactor())
			leftTree = opTree

		case tokenDiv:
			opTree := NewSyntaxTree(currentToken, expK)
			match(tokenDiv)
			opTree.AddLeftChild(leftTree)
			opTree.AddRightChild(parseFactor())
			leftTree = opTree
		}
		tokenType = currentToken.tokenType
	}
	return leftTree
}

func parseFactor() *SyntaxTree {
	switch currentToken.tokenType {
	case tokenLeftParen:
		match(tokenLeftParen)
		tree := parseExp()
		if !match(tokenRightParen) {

		}
		return tree
	case tokenNumber:
		fmt.Printf("match number:%s\n", currentToken.lexeme)
		tree := NewSyntaxTree(currentToken, expK)
		match(tokenNumber)
		return tree
	case tokenId:
		fmt.Printf("match id:%s\n", currentToken.lexeme)
		tree := NewSyntaxTree(currentToken, expK)
		match(tokenId)
		return tree
	default:
		fmt.Printf("parse factor default\n")
	}
	return nil
}

func matchStr(lexeme string) bool {
	if lexeme == currentToken.lexeme {
		currentToken = GetToken()
		return true
	}
	return false
}
func match(token int) bool {
	if token == currentToken.tokenType {
		currentToken = GetToken()
		return true
	}
	return false
}
