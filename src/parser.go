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
		//TODO: how about move all this to function parseIfStmt
		ifTree := NewSyntaxTree(currentToken, stmtK, ifK)
		match(tokenIf)
		parseIfStmt(ifTree)
		//note: semicolone should not follow if-stmt
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

/*
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
*/
func parseStmt() {

}

func parseIfStmt(parent *SyntaxTree) *SyntaxTree {
	tree := parseExp()
	parent.AddChild(tree)
	if !match(tokenThen) {
		//TODO:
	}
	//TODO: don't need then any more, need the body
	/*
		slibling := NewSyntaxTree(lastToken, expK, 0)
		tree.AddSibling(slibling)
		parseStmtSequence(slibling)
	*/
	thenBody := parseStmtSeq()
	parent.AddChild(thenBody)
	//optional
	/*
		if match(tokenElse) {
			parseStmtSequence(nil)
		}
	*/
	//note: this is better than code "if match(tokenElse)"
	if currentToken.tokenType == tokenElse {
		match(tokenElse)
		elseBody := parseStmtSeq()
		parent.AddChild(elseBody)
	}
	if !match(tokenEndBlock) {

	}
	return tree
}

//New version
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

//old version
/*
func parseRepeatStmt(root *SyntaxTree) {
	//parseStmtSequence(root)
	repBody := parseStmtSeq()
	root.AddChild(repBody)
	if !match(tokenUntil) {
		//TODO:
	}
	root.AddChild(parseExp())
	match(tokenSemi)
}
*/
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
