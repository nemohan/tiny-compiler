package main

import (
	"errors"
	//"fmt"
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
	nodeNone = iota
	stmtK
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
	voidK = iota
	opK
	constK
	idK
)

//TODO: use a struct may be better
//var currentLexeme int
var currentToken *tokenSymbol
var parseErr error

func Parse() *SyntaxTree {
	currentToken = GetToken()
	/*
		astRoot := NewSyntaxTree(nil, fileK, voidK)
		parseStmtSequence(astRoot)
	*/
	astRoot := NewSyntaxTree(nil, fileK, voidK)
	child := parseStmtSeq()
	astRoot.AddLeftChild(child)
	astRoot.Traverse()
	astRoot.DFSTraverse()
	return astRoot
}

func parseStmtSeq() *SyntaxTree {
	//for parseErr == nil && currentToken.tokenType != tokenEof {
	if parseErr != nil {
		return nil
	}
	switch currentToken.tokenType {
	case tokenId:
		left := NewSyntaxTree(currentToken, expK, idK)
		match(tokenId)
		lastToken := currentToken
		if !match(tokenAssign) {
			parseErr = errors.New(" tokenAssign")
			//return tree
			return nil
		}
		assignTree := NewSyntaxTree(lastToken, stmtK, assignK)
		assignTree.AddLeftChild(left)
		right := parseAssignStmt()
		assignTree.AddRightChild(right)
		//root.AddChild(assignTree)
		//assignTree.Traverse()
		//don't care match or not
		match(tokenSemi)
		assignTree.AddSlibling(parseStmtSeq())
		return assignTree
	case tokenIf:
		lastToken := currentToken
		//TODO: error
		match(tokenIf)
		ifTree := NewSyntaxTree(lastToken, stmtK, ifK)
		ifTree.AddLeftChild(parseIfStmt())
		//root.AddChild(ifTree)
		//TODO: semicolone should not follow if-stmt
		match(tokenSemi)
		ifTree.AddSlibling(parseStmtSeq())
		return ifTree
	case tokenRead:
		lastToken := currentToken
		match(tokenRead)
		readTree := NewSyntaxTree(lastToken, stmtK, readK)
		readTree.AddLeftChild(parseReadStmt())
		//root.AddChild(readTree)
		match(tokenSemi)
		readTree.AddSlibling(parseStmtSeq())
		return readTree
	case tokenRepeat:
		lastToken := currentToken
		match(tokenRepeat)
		repeatTree := NewSyntaxTree(lastToken, stmtK, repeatK)
		//root.AddChild(repeatTree)
		parseRepeatStmt(repeatTree)
		repeatTree.AddSlibling(parseStmtSeq())
		return repeatTree
	case tokenWrite:
		lastToken := currentToken
		match(tokenWrite)
		writeTree := NewSyntaxTree(lastToken, stmtK, writeK)
		//root.AddChild(writeTree)
		writeTree.AddChild(parseWriteStmt())
		match(tokenSemi)
		writeTree.AddSlibling(parseStmtSeq())
		return writeTree
	case tokenEOF:
		return nil
	default:
		//panic("default parseStmtSeq")
		return nil
	}
	//}
	return nil
}

func parseStmtSequence(root *SyntaxTree) *SyntaxTree {
	var tree *SyntaxTree
	for parseErr == nil && currentToken.tokenType != tokenEOF {
		switch currentToken.tokenType {
		case tokenId:
			left := NewSyntaxTree(currentToken, expK, idK)
			match(tokenId)
			lastToken := currentToken
			if !match(tokenAssign) {
				parseErr = errors.New(" tokenAssign")
				return tree
			}
			assignTree := NewSyntaxTree(lastToken, stmtK, assignK)
			assignTree.AddLeftChild(left)
			right := parseAssignStmt()
			assignTree.AddRightChild(right)
			root.AddChild(assignTree)
			assignTree.Traverse()
			//don't care match or not
			match(tokenSemi)
		default:
			lastToken := currentToken
			//TODO: error
			if match(tokenIf) {
				ifTree := NewSyntaxTree(lastToken, stmtK, ifK)
				ifTree.AddLeftChild(parseIfStmt())
				root.AddChild(ifTree)
				match(tokenSemi)
				continue
			}

			if match(tokenRead) {
				readTree := NewSyntaxTree(lastToken, stmtK, readK)
				readTree.AddLeftChild(parseReadStmt())
				root.AddChild(readTree)
				match(tokenSemi)
				continue
			}
			if match(tokenRepeat) {
				repeatTree := NewSyntaxTree(lastToken, stmtK, repeatK)
				root.AddChild(repeatTree)
				parseRepeatStmt(repeatTree)
				continue
			}
			if match(tokenWrite) {
				writeTree := NewSyntaxTree(lastToken, stmtK, writeK)
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
	//lastToken := currentToken
	if !match(tokenThen) {
		//TODO:
	}
	//TODO: don't need then any more, need the body
	/*
		slibling := NewSyntaxTree(lastToken, expK, 0)
		tree.AddSlibling(slibling)
		parseStmtSequence(slibling)
	*/
	thenBody := parseStmtSeq()
	tree.AddSlibling(thenBody)
	//optional
	if match(tokenElse) {
		parseStmtSequence(nil)
	}

	if !match(tokenEndBlock) {

	}
	return tree
}

func parseRepeatStmt(root *SyntaxTree) {
	parseStmtSequence(root)
	if !match(tokenUntil) {
		//TODO:
	}
	root.AddChild(parseExp())
	match(tokenSemi)
}

func parseAssignStmt() *SyntaxTree {
	tree := parseExp()
	return tree
}

func parseReadStmt() *SyntaxTree {
	switch currentToken.tokenType {
	case tokenId:
		tree := NewSyntaxTree(currentToken, expK, idK)
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
		node := NewSyntaxTree(currentToken, expK, opK)
		node.AddLeftChild(expTree)
		match(tokenLess)
		expTree = handleExp(node)
		expTree.Traverse()

	case tokenEqual:
		node := NewSyntaxTree(currentToken, expK, opK)
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
	for tokenType == tokenAdd || tokenType == tokenMinus {
		switch tokenType {
		case tokenAdd:
			tree := NewSyntaxTree(currentToken, expK, opK)
			tree.AddLeftChild(leftTree)
			match(tokenAdd)
			tree.AddRightChild(parseTerm())
			leftTree = tree

		case tokenMinus:
			tree := NewSyntaxTree(currentToken, expK, opK)
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
			opTree := NewSyntaxTree(currentToken, expK, opK)
			match(tokenMultiply)
			opTree.AddLeftChild(leftTree)
			opTree.AddRightChild(parseFactor())
			leftTree = opTree

		case tokenDiv:
			opTree := NewSyntaxTree(currentToken, expK, opK)
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
		tree := NewSyntaxTree(currentToken, expK, constK)
		match(tokenNumber)
		return tree
	case tokenId:
		tree := NewSyntaxTree(currentToken, expK, idK)
		match(tokenId)
		return tree
	default:
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
