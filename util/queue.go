package util

import "errors"

// Queue is a generic FIFO queue.
type Queue[T any] struct {
	Items []T
}

// New returns a new, empty queue.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{Items: []T{}}
}

// Enqueue adds an item to the end of the queue.
func (q *Queue[T]) Enqueue(item T) {
	q.Items = append(q.Items, item)
}

// EnqueueFront adds an item to the front of the queue.
func (q *Queue[T]) EnqueueFront(item T) {
	q.Items = append([]T{item}, q.Items...)
}

// Dequeue removes and returns the item at the front of the queue.
// Returns an error if the queue is empty.
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	if len(q.Items) == 0 {
		return zero, errors.New("queue is empty")
	}
	item := q.Items[0]
	q.Items = q.Items[1:]
	return item, nil
}

// Peek returns the item at the front of the queue without removing it.
// Returns an error if the queue is empty.
func (q *Queue[T]) Peek() (T, error) {
	var zero T
	if len(q.Items) == 0 {
		return zero, errors.New("queue is empty")
	}
	return q.Items[0], nil
}

// Len returns the number of Items in the queue.
func (q *Queue[T]) Len() int {
	return len(q.Items)
}

// IsEmpty returns true if the queue has no Items.
func (q *Queue[T]) IsEmpty() bool {
	return len(q.Items) == 0
}

func (q *Queue[T]) Clear() {
	q.Items = q.Items[:0]
}
