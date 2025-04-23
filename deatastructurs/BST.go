package main

import (
	"fmt"
)

type Node struct {
	value int
	left  *Node
	right *Node
}

func (n *Node) find(v int) *Node {

	if n == nil {
		return nil
	} else if n.value == v {
		return n
	}

	if v < n.value {
		return n.left.find(v)
	} else {
		return n.right.find(v)
	}
}

func (n *Node) InOrderTraversal() {
	if n == nil {
		return
	}

	n.left.InOrderTraversal()
	fmt.Println(n.value)
	n.right.InOrderTraversal()
}

func (n *Node) add(v int) {

	if n == nil {
		return
	}

	if v < n.value {
		if n.left == nil {
			n.left = &Node{v, nil, nil}
		} else {
			n.left.add(v)
		}
	} else {
		if n.right == nil {
			n.right = &Node{v, nil, nil}
		} else {
			n.right.add(v)
		}
	}
}

func (n *Node) delete(v int) *Node {
	var predessesor *Node
	if n == nil {
		return nil
	} else if n.value == v {

		if n.left == nil && n.right == nil {
			return nil
		} else if n.left == nil && n.right != nil {
			n.value = n.right.value
			predessesor = n.right
			n.right = predessesor.right
			predessesor = nil
		} else if n.right == nil && n.left != nil {
			n.value = n.left.value
			predessesor = n.left
			n.left = predessesor.left
			predessesor = nil
		} else {
			n.value = n.right.value
			predessesor = n.right
			n.right = predessesor.right
			predessesor = nil
		}
	} else if v < n.value {
		n.left = n.left.delete(v)
	} else {
		n.right = n.right.delete(v)
	}

	return n
}

type Tree struct {
	root   *Node
	length int
}

func (t Tree) add(v int) {
	if t.root == nil {
		return
	} else {
		t.root.add(v)
	}
}

func (t *Tree) InOrderTraversal() {
	t.root.InOrderTraversal()
}

func (t *Tree) delete(v int) {
	t.root = t.root.delete(v)
}
