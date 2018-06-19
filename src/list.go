package main

type node struct {
	next **node
}

type list struct {
	head *node
	size int
}
