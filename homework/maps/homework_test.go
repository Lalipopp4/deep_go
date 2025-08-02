package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node[K cmp.Ordered, V any] struct {
	Key         K
	Value       V
	Left, Right *Node[K, V]
}

type OrderedMap[K cmp.Ordered, V any] struct {
	size int
	root *Node[K, V]
}

func NewOrderedMap[K cmp.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func (m *OrderedMap[K, V]) insertPrev(prev, node *Node[K, V], isLeft bool) {
	switch {
	case prev == nil:
		m.root = node
	case isLeft:
		prev.Left = node
	default:
		prev.Right = node
	}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	var (
		isLeft  bool
		curNode = m.root
		prev    = curNode
	)

	for curNode != nil {
		switch {
		case curNode.Key > key:
			prev = curNode
			curNode = curNode.Left
			isLeft = true
		case curNode.Key < key:
			prev = curNode
			curNode = curNode.Right
			isLeft = false
		}
	}

	m.size++
	m.insertPrev(prev, &Node[K, V]{Key: key, Value: value}, isLeft)
}

func (m *OrderedMap[K, V]) deleteNode(prev, node *Node[K, V], isLeft bool) {
	switch {
	case node.Right == nil && node.Left == nil:
		m.insertPrev(prev, nil, isLeft)
	case node.Right == nil:
		m.insertPrev(prev, node.Left, isLeft)
	case node.Left == nil:
		m.insertPrev(prev, node.Right, isLeft)
	default:
		curNode := node.Right
		prevToMove := curNode
		for curNode.Left != nil {
			prevToMove = curNode
			curNode = curNode.Left
		}

		prevToMove.Left = curNode.Right
		m.insertPrev(prev, curNode, isLeft)
	}
}

func (m *OrderedMap[K, V]) Erase(key K) {
	var (
		curNode = m.root
		prev    = curNode
		isLeft  bool
	)

	for curNode != nil {
		switch {
		case curNode.Key > key:
			prev = curNode
			curNode = curNode.Left
			isLeft = true
		case curNode.Key < key:
			prev = curNode
			curNode = curNode.Right
			isLeft = false
		default:
			m.deleteNode(prev, curNode, isLeft)
			m.size--
			return
		}
	}
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	curNode := m.root
	for curNode != nil {
		switch {
		case curNode.Key > key:
			curNode = curNode.Left
		case curNode.Key < key:
			curNode = curNode.Right
		default:
			return true
		}
	}

	return false
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func forEach[K cmp.Ordered, V any](node *Node[K, V], action func(K, V)) {
	if node == nil {
		return
	}

	forEach(node.Left, action)
	action(node.Key, node.Value)
	forEach(node.Right, action)
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	forEach(m.root, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
