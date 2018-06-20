package main

import (
	"fmt"
)

type SyntaxTree struct {
	child    *SyntaxTree
	slibling *SyntaxTree
	parent   *SyntaxTree
	nodeKind int
	expKind  int
	stmtKind int
	expType  int
	token    *tokenSymbol
	height   int
}

func NewSyntaxTree(token *tokenSymbol, nodeKind int) *SyntaxTree {
	return &SyntaxTree{
		nodeKind: nodeKind,
		token:    token,
	}
}

func (st *SyntaxTree) AddSlibling(node *SyntaxTree) {
	if st == nil {
		return
	}
	next := st
	for ; next.slibling != nil; next = next.slibling {
		//empty
	}
	next.slibling = node

}
func (st *SyntaxTree) AddChild(node *SyntaxTree) {
	if st.child == nil {
		st.AddLeftChild(node)
		return
	}
	st.AddRightChild(node)
}
func (st *SyntaxTree) AddLeftChild(node *SyntaxTree) {
	st.child = node
}

func (st *SyntaxTree) AddRightChild(node *SyntaxTree) {
	if st.child == nil {
		msg := fmt.Sprintf("node:%v child is empty\n", st)
		panic(msg)
	}
	next := st.child
	for ; next.slibling != nil; next = next.slibling {
		//empty
	}
	next.slibling = node
}

func (st *SyntaxTree) DFSTraverse() {
	if st == nil {
		return
	}
	//stack := []*SyntaxTree{st.root}
	dfsTraverse(st)
}

func tabNum(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s = s + "\t"
	}
	return s
}
func dfsTraverse(node *SyntaxTree) {
	if node == nil {
		return
	}
	fmt.Printf("node:%d\n", node.nodeKind)
	dfsTraverse(node.child)
	dfsTraverse(node.slibling)
}

func (st *SyntaxTree) Traverse() {
	queue := []*SyntaxTree{st}
	traverse(queue)
}
func traverse(queue []*SyntaxTree) {
	if len(queue) == 0 {
		return
	}
	node := queue[0]
	queue = append(queue[:0], queue[1:]...)
	fmt.Printf("kind:%d     token:%v height:%d\n", node.nodeKind, node.token, node.height)
	next := node.child
	for ; next != nil; next = next.slibling {
		next.height = node.height + 1
		queue = append(queue, next)
	}
	traverse(queue)
}
