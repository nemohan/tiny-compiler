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
	astRoot := NewSyntaxTree(nil, fileK, voidK)
	child := parseStmtSeq()
	astRoot.AddChild(child)
	astRoot.Traverse()
	astRoot.DFSTraverse()
	return astRoot
}

//TODO: handle syntax error ---------------
func parseStmtSeq() *SyntaxTree {
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
		assignTree.AddChild(left)
		right := parseAssignStmt()
		assignTree.AddChild(right)
		//root.AddChild(assignTree)
		//don't care match or not
		match(tokenSemi)
		assignTree.AddSibling(parseStmtSeq())
		return assignTree
	case tokenIf:
		ifTree := parseIfStmt()
		match(tokenSemi)
		ifTree.AddSibling(parseStmtSeq())
		return ifTree

	case tokenRead:
		readTree := NewSyntaxTree(currentToken, stmtK, readK)
		match(tokenRead)
		readTree.AddChild(parseReadStmt())
		//note: semicolon as seperator
		//TODO: this will assume a semicolon must follow read
		match(tokenSemi)
		readTree.AddSibling(parseStmtSeq())
		return readTree

	case tokenRepeat:
		repTree := parseRepeatStmt()
		sibling := parseStmtSeq()
		repTree.AddSibling(sibling)
		return repTree
	case tokenWrite:
		writeTree := NewSyntaxTree(currentToken, stmtK, writeK)
		match(tokenWrite)
		writeTree.AddChild(parseWriteStmt())
		if currentToken.tokenType == tokenSemi {
			match(tokenSemi)
		}
		writeTree.AddSibling(parseStmtSeq())
		return writeTree
	case tokenEOF:
		return nil
	default:
		//panic("default parseStmtSeq")
		Logf("default: %s\n", currentToken.lexeme)
		return nil
	}
	//}
	return nil
}

func parseIfStmt() *SyntaxTree {
	ifTree := NewSyntaxTree(currentToken, stmtK, ifK)
	match(tokenIf)

	tree := parseExp()
	ifTree.AddChild(tree)
	if !match(tokenThen) {
		//TODO:
	}

	//TODO: don't need then any more, need the body
	thenBody := parseStmtSeq()
	ifTree.AddChild(thenBody)
	if currentToken.tokenType == tokenElse {
		match(tokenElse)
		elseBody := parseStmtSeq()
		ifTree.AddChild(elseBody)
	}
	if !match(tokenEndBlock) {

	}
	return ifTree
}

func parseRepeatStmt() *SyntaxTree {
	repeatTree := NewSyntaxTree(currentToken, stmtK, repeatK)
	match(tokenRepeat)
	repBody := parseStmtSeq()
	repeatTree.AddChild(repBody)
	if !match(tokenUntil) {
		//TODO:
	}
	repeatTree.AddChild(parseExp())
	match(tokenSemi)
	return repeatTree
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
	tree.AddChild(right)
	return tree
}

func parseExp() *SyntaxTree {
	expTree := parseSimpleExp()
	switch currentToken.tokenType {
	case tokenLess:
		node := NewSyntaxTree(currentToken, expK, opK)
		node.AddChild(expTree)
		match(tokenLess)
		expTree = handleExp(node)
		return node
	case tokenEqual:
		node := NewSyntaxTree(currentToken, expK, opK)
		node.AddChild(expTree)
		match(tokenEqual)
		expTree = handleExp(expTree)
		return node
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
			tree.AddChild(leftTree)
			match(tokenAdd)
			tree.AddChild(parseTerm())
			leftTree = tree

		case tokenMinus:
			tree := NewSyntaxTree(currentToken, expK, opK)
			tree.AddChild(leftTree)
			match(tokenMinus)
			tree.AddChild(parseTerm())
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
			opTree.AddChild(leftTree)
			opTree.AddChild(parseFactor())
			leftTree = opTree

		case tokenDiv:
			opTree := NewSyntaxTree(currentToken, expK, opK)
			match(tokenDiv)
			opTree.AddChild(leftTree)
			opTree.AddChild(parseFactor())
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
