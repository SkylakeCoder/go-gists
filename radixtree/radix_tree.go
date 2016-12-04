package main

import "fmt"
import "strings"

type RadixTree struct {
	tree *node
}

type node struct {
	label    byte
	prefix   string
	parent   *node
	children []*node
}

func NewRadixTree() *RadixTree {
	return &RadixTree{
		tree: &node{},
	}
}

func (r *RadixTree) Dump() {
	children := r.tree.children
	if len(children) == 0 {
		return
	}
	leafNodes := []*node{}
	for len(children) > 0 {
		node := children[0]
		children = children[1:]
		if len(node.children) > 0 {
			for _, n := range node.children {
				children = append(children, n)
			}
		} else {
			leafNodes = append(leafNodes, node)
		}
	}
	for _, n := range leafNodes {
		dump := ""
		for n != nil {
			dump = n.prefix + " -> " + dump
			n = n.parent
		}
		dump = dump[:len(dump)-len(" -> ")]
		fmt.Println(dump)
	}
}

func (r *RadixTree) Find(key string) bool {
	root := r.tree
	if key[0] != root.label {
		return false
	}
	if !strings.HasPrefix(key, root.prefix) {
		return false
	}
	key = key[len(root.prefix):]
	node := root.findNext(key)
	for node != nil && !node.isLeaf() {
		key = key[len(node.prefix):]
		node = node.findNext(key)
	}
	if node != nil && node.isLeaf() {
		key = key[len(node.prefix):]
	}
	return node != nil && key == ""
}

// Insert : The following code was copied from echo/router.go/insert(...) method.
// But it have been simplied a lot here.
func (r *RadixTree) Insert(key string) {
	cn := r.tree
	search := key
	for {
		sl := len(search)
		pl := len(cn.prefix)
		l := 0

		// LCP
		max := pl
		if sl < max {
			max = sl
		}
		for ; l < max && search[l] == cn.prefix[l]; l++ {
		}
		if l == 0 {
			// At root node
			cn.label = search[0]
			cn.prefix = search
		} else if l < pl {
			// Split node
			n := newNode(cn.prefix[l:], cn, cn.children)

			// Reset parent node
			cn.label = cn.prefix[0]
			cn.prefix = cn.prefix[:l]
			cn.children = nil

			cn.addChild(n)

			if l == sl {
				// At parent node
			} else {
				// Create child node
				n = newNode(search[l:], cn, nil)
				cn.addChild(n)
			}
		} else if l < sl {
			search = search[l:]
			c := cn.findChildWithLabel(search[0])
			if c != nil {
				// Go deeper
				cn = c
				continue
			}
			// Create child node
			n := newNode(search, cn, nil)
			cn.addChild(n)
		}
		break
	}
}

func newNode(pre string, p *node, c []*node) *node {
	n := &node{
		label:    pre[0],
		prefix:   pre,
		parent:   p,
		children: c,
	}
	for _, cn := range c {
		cn.parent = n
	}
	return n
}

func (n *node) addChild(c *node) {
	n.children = append(n.children, c)
}

func (n *node) findChildWithLabel(l byte) *node {
	for _, c := range n.children {
		if c.label == l {
			return c
		}
	}
	return nil
}

func (n *node) findNext(key string) *node {
	for _, c := range n.children {
		if key[0] != c.label {
			continue
		}
		if !strings.HasPrefix(key, c.prefix) {
			continue
		}
		return c
	}
	return nil
}

func (n *node) isLeaf() bool {
	return len(n.children) == 0
}
