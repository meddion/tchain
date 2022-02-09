package database

import (
	"bytes"
	"math/rand"
)

type nodeLists struct {
	key, value []byte
	isHead     bool
	levels     []*nodeLists
}

func newNodeLists(key, value []byte, height int, isHead bool) *nodeLists {
	return &nodeLists{
		key:    key,
		value:  value,
		isHead: isHead,
		levels: make([]*nodeLists, height),
	}
}

type SkipLists struct {
	height int
	head   *nodeLists
}

func NewSkipLists() *SkipLists {
	return &SkipLists{
		height: 1,
		head:   newNodeLists(nil, nil, 1, true),
	}
}

func randomHeight(p float64) int {
	height := 1
	for rand.Float64() < p {
		height++
	}

	return height
}

// Search returns a *node holding the key,
// otherwise a pointer to witch append the key.
func (s *SkipLists) search(key []byte) ([]*nodeLists, *nodeLists) {
	prevNodes := make([]*nodeLists, s.height)

	curNode := s.head
	level := s.height - 1

	for level > 0 {
		nextNode := curNode.levels[level]

		if nextNode != nil && bytes.Equal(key, nextNode.key) {
			curNode = nextNode
			continue
		}

		prevNodes[level] = curNode
		level--
	}

	return prevNodes, curNode
}

func (s *SkipLists) Insert(key, value []byte) *nodeLists {
	prevLists, p := s.search(key)
	if bytes.Equal(p.key, key) {
		return nil
	}

	newHeight := randomHeight(0.5)
	target := newNodeLists(key, value, newHeight, false)

	// append those lists
	if newHeight > s.height {
		for i := newHeight - 1; i >= 0; i-- {
			// s.head = append(s.head.levels, prevLists...)
		}

		s.height = newHeight
	}

	for _, n := range prevLists {

	}

	return target
}

func insertAfter(cur, target *nodeLists) {
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
