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
			lastToken := currentToken
			if !match(tokenAssign) {
				parseErr = errors.New(" tokenAssign")
				return tree
			}
			assignTree := NewSyntaxTree(lastToken, expK)
			assignTree.AddLeftChild(left)
			right := parseAssignStmt()
			assignTree.AddRightChild(right)
			root.AddChild(assignTree)
			assignTree.Traverse()
			fmt.Printf("before match:%v\n", currentToken)
			//don't care match or not
			match(tokenSemi)
			fmt.Printf("token:%v--------------\n", currentToken)
		default:
			lastToken := currentToken
			//TODO: error
			if match(tokenIf) {
				ifTree := NewSyntaxTree(lastToken, stmtK)
				fmt.Printf("match if\n")
				ifTree.AddLeftChild(parseIfStmt())
				root.AddChild(ifTree)
				match(tokenSemi)
				continue
			}

			if match(tokenRead) {
				readTree := NewSyntaxTree(lastToken, stmtK)
				readTree.AddLeftChild(parseReadStmt())
				root.AddChild(readTree)
				match(tokenSemi)
				fmt.Printf("parse read done\n")
				continue
			}
			if match(tokenRepeat) {
				repeatTree := NewSyntaxTree(lastToken, stmtK)
				root.AddChild(repeatTree)
				parseRepeatStmt(repeatTree)
				continue
			}
			if match(tokenWrite) {
				writeTree := NewSyntaxTree(lastToken, stmtK)
				root.AddChild(writeTree)
				writeTree.AddChild(parseWriteStmt())
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
	if !match(tokenThen) {
		//TODO:
	}
	slibling := NewSyntaxTree(lastToken, expK)
	tree.AddSlibling(slibling)
	fmt.Printf("match then\n")
	parseStmtSequence(slibling)
	//slibling.AddLeftChild(child)

	//optional
	if match(tokenElse) {
		parseStmtSequence(nil)
	}

	if !match(tokenEndBlock) {

	}
	fmt.Printf("match end\n")
	return tree
}

func parseRepeatStmt(root *SyntaxTree) {
	fmt.Printf("parseReat xxxxxxxxxxxxxxxxxxxx\n")
	parseStmtSequence(root)
	if !match(tokenUntil) {
		//TODO:
	}
	root.AddChild(parseExp())
	match(tokenSemi)
	fmt.Printf("parse Repeat========\n")
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

func parseWriteStmt() *SyntaxTree {
	return parseExp()
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
		fmt.Printf("match token equal\n")
		node := NewSyntaxTree(currentToken, expK)
		node.AddLeftChild(expTree)
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
	fmt.Printf("ParseSimpleExp:%v add:%d minus:%d\n", currentToken, tokenAdd, tokenMinus)
	for tokenType == tokenAdd || tokenType == tokenMinus {
		fmt.Printf("fuck:%v\n", currentToken)
		switch tokenType {
		case tokenAdd:
			tree := NewSyntaxTree(currentToken, expK)
			tree.AddLeftChild(leftTree)
			match(tokenAdd)
			tree.AddRightChild(parseTerm())
			leftTree = tree

		case tokenMinus:
			fmt.Printf("match token minus\n")
			tree := NewSyntaxTree(currentToken, expK)
			tree.AddLeftChild(leftTree)
			match(tokenMinus)
			tree.AddRightChild(parseTerm())
			leftTree = tree
		default:
			panic("parseSimpleExp")
		}
		tokenType = currentToken.tokenType
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
	case tokenLParen:
		match(tokenLParen)
		tree := parseExp()
		if !match(tokenRParen) {

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

func match(token int) bool {
	if token == currentToken.tokenType {
		currentToken = GetToken()
		return true
	}
	return false
}
