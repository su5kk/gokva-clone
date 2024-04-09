package container

import "sync"

type Queue[T any] struct {
	sync.Mutex

	items []T
}

func (q *Queue[T]) Enqueue(item T) {
	q.Lock()
	q.items = append(q.items, item)
	q.Unlock()
}

func (q *Queue[T]) Dequeue() T {
	q.Lock()
	val := q.items[0]
	q.items = q.items[1:]
	q.Unlock()

	return val
}
