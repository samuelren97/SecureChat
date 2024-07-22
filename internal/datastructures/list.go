package datastructures

import (
	"errors"
	"fmt"
)

type List[T any] struct {
	Count  int
	values []T
}

func NewList[T any]() *List[T] {
	return &List[T]{values: make([]T, 0), Count: 0}
}

// Public
func (l *List[T]) Add(t T) {
	l.values = append(l.values, t)
	l.Count++
}

func (l *List[T]) Remove(index int) {
	l.validateIndex(&index)
	l.values = append(l.values[:index], l.values[index+1:]...)
	l.Count--
}

func (l *List[T]) Get(index int) T {
	l.validateIndex(&index)
	return l.values[index]
}

func (l *List[T]) FindIndex(fn func(T) bool) int {
	for i := 0; i < len(l.values); i++ {
		if fn(l.values[i]) {
			return i
		}
	}

	return -1
}

func (l *List[T]) ForEach(f func(T)) {
	for _, t := range l.values {
		f(t)
	}
}

func (l *List[T]) BubbleSort(f func(a, b T) bool) {
	for i := 0; i < l.Count-1; i++ {
		for k := 0; k < l.Count-i-1; k++ {
			if f(l.values[k], l.values[k+1]) {
				tmp := l.values[k]
				l.values[k] = l.values[k+1]
				l.values[k+1] = tmp
			}
		}
	}
}

// Private
func (l *List[T]) validateIndex(index *int) {
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
