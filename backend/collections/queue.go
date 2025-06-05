package collections

import "errors"

type Queue[T any] struct {
	values[]T
}

func (queue *Queue[T]) Push(value T) {
	if queue.values == nil {
		queue.values = []T{value}
		return
	}

	queue.values = append(queue.values, value)
}

func (queue *Queue[T]) Pop() (T, error) {
	if queue.IsEmpty() {
		var result T
		return result, errors.New("Cannot invoke Pop() on empty Queue.")
	}

	value := queue.values[0]
	queue.values = queue.values[1:]
	return value, nil
}

func (queue *Queue[T]) Peek() (T, error) {
	if queue.IsEmpty() {
		var result T
		return result, errors.New("Cannot invoke Peek() on empty Queue.")
	}
	return queue.values[0], nil
}

func (queue *Queue[T]) Length() int {
	if queue.values == nil {
		return 0
	}

	return len(queue.values)
}

func (queue *Queue[T]) IsEmpty() bool {
	return queue.values == nil || len(queue.values) == 0
}
