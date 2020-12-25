package main

import "fmt"

type Person struct {
	FirstName, MiddleName, LastName string
}

func (p *Person) Names() []string {
	return []string{p.FirstName, p.MiddleName, p.LastName}
}

type Node struct {
	Value               int
	left, right, parent *Node
}

func NewNode(value int, left *Node, right *Node) *Node {
	n := &Node{
		Value: value,
		left:  left,
		right: right,
	}
	left.parent = n
	right.parent = n
	return n
}

func NewTerminalNode(value int) *Node {
	return &Node{Value: value}
}

type InOrderIterator struct {
	Current     *Node
	root        *Node
	returnStart bool
}

func NewInOrderIterator(root *Node) *InOrderIterator {
	i := &InOrderIterator{
		Current:     root,
		root:        root,
		returnStart: false,
	}
	for i.Current.left != nil {
		i.Current = i.Current.left
	}
	return i
}

func (i *InOrderIterator) Reset() {
	i.Current = i.root
	i.returnStart = false
}

func (i *InOrderIterator) MoveNext() bool {
	if i.Current == nil {
		return false
	}
	if !i.returnStart {
		i.returnStart = true
		return true
	}

	// todo
}

func main() {
	p := Person{
		"A", "B", "C",
	}
	for _, n := range p.Names() {
		fmt.Println(n)
	}
}
