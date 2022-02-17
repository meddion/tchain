package database

import (
	"bytes"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type node struct {
	key, value []byte
	isHead     bool
	levels     []*node
}

func newNode(key, value []byte, isHead bool) *node {
	return &node{
		key:    key,
		value:  value,
		isHead: isHead,
		levels: make([]*node, 0),
	}
}

func (n *node) Get(level int) *node {
	if level >= len(n.levels) {
		return nil
	}

	return n.levels[level]
}

func (n *node) Set(level int, target *node) {
	if level >= len(n.levels) {
		panic("out of range indexing in node.Set")
	}

	n.levels[level] = target
}

func (n *node) Append(nodes ...*node) {
	n.levels = append(n.levels, nodes...)
}

type SkipLists struct {
	*node
}

func NewSkipLists() *SkipLists {
	return &SkipLists{
		node: newNode(nil, nil, true),
	}
}

func randomHeight(p float64) int {

	height := 1
	for rand.Float64() < p {
		height++
	}

	return height
}

// Search returns a slice of nodes which lead to the target node,
// as the second argument the matching node is returned (key == n.key)
// or if a node doesn't exist -- return a node after which to append
func (s *SkipLists) search(key []byte) ([]*node, *node) {
	pathNodes := make([]*node, s.Height())

	curNode := s.node
	level := s.Height() - 1

	for level > 0 {
		nextNode := curNode.Get(level)

		// key > nextNode.key
		if nextNode != nil && bytes.Compare(key, nextNode.key) == 1 {
			curNode = nextNode
			continue
		}

		pathNodes[level] = curNode
		level--
	}

	return pathNodes, curNode
}

func (s *SkipLists) Insert(key, value []byte) *node {
	prevLists, p := s.search(key)
	if bytes.Equal(p.key, key) {
		return nil
	}

	newHeight := randomHeight(0.5)
	target := newNode(key, value, false)

	for i := s.Height(); i < newHeight; i++ {
		s.Append(target)
	}

	for level := 0; level < newHeight; level++ {
		prevLists[level].Set(level, target)
	}

	return target
}

func (s *SkipLists) Height() int {
	return len(s.levels)
}

func insertAfter(cur, target *node) {
	// save pointers to each of the following nodes in currentNode
	// before inserting a newNode
	// allocate enough space in levels to hold "the following nodes"

	// if newHeight > curHeight then
	// 1) set new  height
	// 2) set the prev pointers to the new heights

	// tmp := cur.next
	// cur.next = target
	// target.prev = tmp
}

//
// func (s *SkipLists) Delete(key []byte) bool {
//     return true
// }
