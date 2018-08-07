package main

import (
	"fmt"
	"testing"
)

const stmtK = 0

func Logf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func TestTree(t *testing.T) {
	rootToken := newToken()
	rootToken.lexeme = "root"
	root := NewSyntaxTree(rootToken, 1, 2)

	funcToken := newToken()
	funcToken.lexeme = "func"
	funcChild := NewSyntaxTree(funcToken, 1, 2)
	root.AddChild(funcChild)

	param1Token := newToken()
	param1Token.lexeme = "int"
	paramChild := NewSyntaxTree(param1Token, 1, 2)
	funcChild.AddChild(paramChild)

	param2Token := newToken()
	param2Token.lexeme = "string"
	param2Child := NewSyntaxTree(param2Token, 1, 2)
	paramChild.AddSibling(param2Child)

	retToken := newToken()
	retToken.lexeme = "void"
	retChild := NewSyntaxTree(retToken, 1, 2)
	funcChild.AddChild(retChild)

	floatToken := newToken()
	floatToken.lexeme = "float"
	retChild1 := NewSyntaxTree(floatToken, 1, 2)
	retChild.AddSibling(retChild1)

	//where my char
	charToken := newToken()
	charToken.lexeme = "char"
	retChild2 := NewSyntaxTree(charToken, 1, 2)
	retChild1.AddSibling(retChild2)

	forToken := newToken()
	forToken.lexeme = "for"
	forChild := NewSyntaxTree(forToken, 1, 2)
	root.AddChild(forChild)

	expToken := newToken()
	expToken.lexeme = "exp"
	expChild := NewSyntaxTree(expToken, 1, 2)
	forChild.AddChild(expChild)
	root.Traverse()
	root.DFSTraverse()
}

func TestPreOrder(t *testing.T) {
	rootToken := newToken()
	rootToken.lexeme = "root"
	root := NewSyntaxTree(rootToken, 1, 2)

	funcToken := newToken()
	funcToken.lexeme = "func"
	funcChild := NewSyntaxTree(funcToken, 1, 2)
	root.AddChild(funcChild)

	param1Token := newToken()
	param1Token.lexeme = "int"
	paramChild := NewSyntaxTree(param1Token, 1, 2)
	funcChild.AddChild(paramChild)

	param2Token := newToken()
	param2Token.lexeme = "string"
	param2Child := NewSyntaxTree(param2Token, 1, 2)
	paramChild.AddSibling(param2Child)

	retToken := newToken()
	retToken.lexeme = "void"
	retChild := NewSyntaxTree(retToken, 1, 2)
	funcChild.AddChild(retChild)
	root.DFSTraverse()
}
