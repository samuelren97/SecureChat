package datastructures

import (
	"errors"
	"fmt"
)

// Node struct def

type Node[T any] struct {
	Value    T
	NextNode *Node[T]
}

// Node private constructor

func newNode[T any](value T) *Node[T] {
	return &Node[T]{Value: value, NextNode: nil}
}

// LinkList struct def

type LinkList[T any] struct {
	Count     int
	firstNode *Node[T]
	lastNode  *Node[T]
}

// LinkList Constructor

func NewLinkList[T any]() *LinkList[T] {
	return &LinkList[T]{Count: 0, firstNode: nil, lastNode: nil}
}

// Public

func (l *LinkList[T]) Add(value T) {
	node := newNode(value)
	if l.Count == 0 {
		l.firstNode = node
		l.lastNode = l.firstNode
		l.Count++
		return
	}

	l.lastNode.NextNode = node
	l.lastNode = node
	l.Count++
}

func (l *LinkList[T]) Get(index int) T {
	l.validateIndex(&index)

	if index == 0 {
		return l.firstNode.Value
	}

	if index == l.Count-1 {
		return l.lastNode.Value
	}

	currentNode := l.firstNode
	for i := 0; i < index; i++ {
		currentNode = currentNode.NextNode
	}
	return currentNode.Value
}

func (l *LinkList[T]) Remove(index int) {
	l.validateIndex(&index)

	if l.Count == 1 {
		l.firstNode = nil
		l.lastNode = nil
		l.Count--
		return
	}

	if l.Count == 2 && index == 0 {
		l.firstNode = l.lastNode
		l.Count--
		return
	}

	if l.Count == 2 && index == 1 {
		l.firstNode.NextNode = nil
		l.lastNode = l.firstNode
		l.Count--
		return
	}

	currentNode := l.firstNode
	for i := 0; i < index-1; i++ {
		currentNode = currentNode.NextNode
	}

	nextNode := currentNode.NextNode.NextNode
	currentNode.NextNode = nextNode
	l.Count--
}

func (l *LinkList[T]) ForEach(f func(T)) {
	if l.Count != 0 {
		if l.firstNode.NextNode == nil {
			f(l.firstNode.Value)
			return
		}

		currentNode := l.firstNode
		for currentNode != nil {
			f(currentNode.Value)
			currentNode = currentNode.NextNode
		}
	}
}

func (l *LinkList[T]) Enqueue(t T) {
	newNode := newNode(t)

	if l.Count == 0 {
		l.firstNode = newNode
		l.lastNode = newNode
		l.Count++
		return
	}

	l.lastNode.NextNode = newNode
	l.lastNode = newNode
	l.Count++
}

func (l *LinkList[T]) Dequeue() T {
	value := l.firstNode.Value
	l.firstNode = l.firstNode.NextNode
	l.Count--
	return value
}

func (l *LinkList[T]) Peek() T {
	return l.firstNode.Value
}

func (l *LinkList[T]) IsEmpty() bool {
	return l.Count == 0
}

func (l *LinkList[T]) Clear() {
	l.firstNode = nil
	l.lastNode = nil
}

// Private
func (l *LinkList[T]) validateIndex(index *int) {
	if l.Count == 0 {
		panic(errors.New("list is empty"))
	}

	if *index < 0 {
		panic(errors.New("index out of bounds, cannot get a negative value"))
	}

	if *index >= l.Count {
		err := errors.New("index out of bounds, cannot get " + fmt.Sprint(index) + ", nb of elements: " + fmt.Sprint(l.Count))
		panic(err)
	}
}

func (l *LinkList[T]) getNodeAt(index int) *Node[T] {
	l.validateIndex(&index)
	if index == 0 {
		return l.firstNode
	}

	if index == l.Count-1 {
		return l.lastNode
	}

	currentNode := l.firstNode
	for i := 0; i < index; i++ {
		currentNode = currentNode.NextNode
	}
	return currentNode
}
