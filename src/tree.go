package main

import (
	"fmt"
)

const maxChildNum = 3

type SyntaxTree struct {
	childs   [maxChildNum]*SyntaxTree
	sibling  *SyntaxTree
	parent   *SyntaxTree
	nodeKind int
	expKind  int
	stmtKind int
	expType  int
	token    *tokenSymbol
	height   int
	childIdx int
}

var treeLine = 0

type traverseProc func(*SyntaxTree)

func printTraverseProc(node *SyntaxTree) {
	if node.token == nil {
		return
	}
	Logf("lexeme: %s\n", node.token.lexeme)
}

func emptyTraverseProc(node *SyntaxTree) {
}

func NewSyntaxTree(token *tokenSymbol, nodeKind int, kind int) *SyntaxTree {
	node := &SyntaxTree{
		nodeKind: nodeKind,
		token:    token,
	}
	if nodeKind == stmtK {
		node.stmtKind = kind
	}
	return node
}

func (st *SyntaxTree) AddSibling(node *SyntaxTree) {
	if st == nil {
		panic("add sibling to empty node")
	}
	next := st
	for ; next.sibling != nil; next = next.sibling {
		//empty
	}
	next.sibling = node
}

func (st *SyntaxTree) AddChild(node *SyntaxTree) {
	if st.childIdx >= maxChildNum {
		msg := fmt.Sprintf("child number:%d beyond max:%d\n", st.childIdx, maxChildNum)
		panic(msg)
	}
	st.childs[st.childIdx] = node
	st.childIdx++
}

func (st *SyntaxTree) DFSTraverse() {
	if st == nil {
		return
	}
	dfsTraverse(st)
	treeLine = 0
}

func tabNum(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s = s + "\t"
	}
	return s
}

//preOrder traverse
func dfsTraverse(node *SyntaxTree) {
	if node == nil {
		return
	}
	//TODO: the root node's token is nil
	tokenStr := ""
	if node.token != nil {
		//tokenStr = node.token.String()
		tokenStr = node.token.SimpleStr()
	}
	height := node.height
	Logf("%d %s node:%d height:%d %s\n", treeLine, tabNum(height),
		node.nodeKind, height, tokenStr)
	treeLine++
	for _, c := range node.childs {
		if c == nil {
			continue
		}
		dfsTraverse(c)
		for next := c.sibling; next != nil; next = next.sibling {
			dfsTraverse(next)
		}
	}
}

func GenTraverse(root *SyntaxTree, preProc, postProc traverseProc) {
	if root == nil {
		return
	}
	preProc(root)
	/*
		for next := root.child; next != nil; next = next.sibling {
			GenTraverse(next, preProc, postProc)
		}
	*/
	for _, c := range root.childs {
		if c == nil {
			continue
		}
		GenTraverse(c, preProc, postProc)
		for next := c.sibling; next != nil; next = next.sibling {
			GenTraverse(c, preProc, postProc)
		}
	}
	postProc(root)
}

func (st *SyntaxTree) Traverse() {
	queue := []*SyntaxTree{st}
	traverse(queue)
}

func newIter(node *SyntaxTree) func(**SyntaxTree) bool {
	first := node
	return func(next **SyntaxTree) bool {
		if first == nil || first.sibling == nil {
			return false
		}
		*next = first.sibling
		//buggy
		//first = (*next).sibling
		first = first.sibling
		return true
	}
}

func traverse(queue []*SyntaxTree) {
	if len(queue) == 0 {
		return
	}
	node := queue[0]
	queue = append(queue[:0], queue[1:]...)
	Logf("kind:%d     token:%v height:%d\n", node.nodeKind, node.token, node.height)
	for _, child := range node.childs {
		if child == nil {
			continue
		}

		moreSibling := newIter(child)
		child.height = node.height + 1
		queue = append(queue, child)
		var sibling *SyntaxTree
		for moreSibling(&sibling) {
			(sibling).height = node.height + 1
			queue = append(queue, sibling)
		}
	}
	traverse(queue)
}
