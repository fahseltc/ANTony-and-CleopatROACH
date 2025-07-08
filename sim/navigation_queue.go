package sim

import (
	"fmt"
	"gamejam/vec2"
)

type NavQueue struct {
	Items []*vec2.T
}

func (q *NavQueue) Enqueue(data *vec2.T) {
	q.Items = append(q.Items, data)
}

func (q *NavQueue) Peek() (*vec2.T, error) {
	if q.IsEmpty() {
		return nil, fmt.Errorf("queue is empty")
	}
	return q.Items[0], nil
}

func (q *NavQueue) Dequeue() (*vec2.T, error) {
	if q.IsEmpty() {
		return nil, fmt.Errorf("queue is empty")
	}
	item := q.Items[0]
	q.Items = q.Items[1:]
	return item, nil
}

func (q *NavQueue) IsEmpty() bool {
	return len(q.Items) == 0
}

func (q *NavQueue) Clear() {
	q.Items = q.Items[:0]
}

func (q *NavQueue) Print() {
	for _, item := range q.Items {
		fmt.Print(item, " ")
	}
	fmt.Println()
}
