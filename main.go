package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Hash []byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func hash(data []byte) Hash {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

type Hashable interface {
	hash() Hash
}

type Block string

func (b Block) hash() Hash {
	return hash([]byte(b[:]))
}

type EmptyBlock struct {
}

func (_ EmptyBlock) hash() Hash {
	return []byte{}
}

type Node struct {
	left  Hashable
	right Hashable
}

func (n Node) hash() Hash {
	var l, r []byte
	l = n.left.hash()
	r = n.right.hash()
	return hash(append(l[:], r[:]...))
}

func buildTree(parts []Hashable) []Hashable {
	var nodes []Hashable

	var i int

	for i = 0; i < len(parts); i += 2 {
		if i+1 < len(parts) {
			nodes = append(nodes, Node{left: parts[i], right: parts[i+1]})
		} else {
			nodes = append(nodes, Node{left: parts[i], right: EmptyBlock{}})
		}
	}

	if len(nodes) == 1 {
		return nodes
	} else {
		return buildTree(nodes)
	}
}

func printTree(node Node) {
	printNode(node, 0)
}

func printNode(node Node, level int) {
	fmt.Printf("(%d) %s %s\n", level, strings.Repeat(" ", level), node.hash())
	if l, ok := node.left.(Node); ok {
		printNode(l, level+1)
	} else if l, ok := node.left.(Block); ok {
		fmt.Printf("(%d) %s %s (data: %s)\n", level+1, strings.Repeat(" ", level+1), l.hash(), l)
	}
	if r, ok := node.right.(Node); ok {
		printNode(r, level+1)
	} else if r, ok := node.right.(Block); ok {
		fmt.Printf("(%d) %s %s (data: %s)\n", level+1, strings.Repeat(" ", level+1), r.hash(), r)
	}
}

func main() {
	printTree(buildTree([]Hashable{Block("a"), Block("b"), Block("c"), Block("d")})[0].(Node))
}
