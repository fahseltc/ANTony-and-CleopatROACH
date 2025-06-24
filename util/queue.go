package util

import "errors"

// Queue is a generic FIFO queue.
type Queue[T any] struct {
	items []T
}

// New returns a new, empty queue.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{items: []T{}}
}

// Enqueue adds an item to the end of the queue.
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the item at the front of the queue.
// Returns an error if the queue is empty.
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	if len(q.items) == 0 {
		return zero, errors.New("queue is empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// Peek returns the item at the front of the queue without removing it.
// Returns an error if the queue is empty.
func (q *Queue[T]) Peek() (T, error) {
	var zero T
	if len(q.items) == 0 {
		return zero, errors.New("queue is empty")
	}
	return q.items[0], nil
}

// Len returns the number of items in the queue.
func (q *Queue[T]) Len() int {
	return len(q.items)
}

// IsEmpty returns true if the queue has no items.
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}
