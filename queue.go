package minifac

import "golang.org/x/exp/slices"

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		values: []T{},
	}
}

type Queue[T any] struct {
	values []T
}

func (q *Queue[T]) Enqueue(t T) {
	q.values = append(q.values, t)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.values) == 0 {
		var t T
		return t, false
	}
	elt := q.values[0]
	q.values = slices.Delete(q.values, 0, 1)
	return elt, true
}

func (q *Queue[T]) Len() int {
	return len(q.values)
}

func (q *Queue[T]) Peek() (T, bool) {
	if len(q.values) == 0 {
		var t T
		return t, false
	}
	return q.values[0], true
}
