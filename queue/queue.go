package queue

import (
	"sync"
)

type Queue[T any] struct {
	items []T
	mu    sync.Mutex
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
	}
}

func (q *Queue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		var zero T
		return zero, false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue[T]) Peek() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		var zero T
		return zero, false
	}

	return q.items[0], true
}

func (q *Queue[T]) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

func (q *Queue[T]) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items) == 0
}

func (q *Queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = make([]T, 0)
}
