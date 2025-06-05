package collections

import "errors"

type Stack[T any] struct {
	values[]T
}

func (stack *Stack[T]) Push(value T) {
	if stack.values == nil {
		stack.values = []T{value}
		return
	}

	stack.values = append(stack.values, value)
}

func (stack *Stack[T]) Pop() (T, error) {
	if stack.IsEmpty() {
		var result T
		return result, errors.New("Cannot invoke Pop() on empty Stack.")
	}

	value := stack.values[len(stack.values) - 1]
	stack.values = stack.values[:len(stack.values) - 1]
	return value, nil
}

func (stack *Stack[T]) Peek() (T, error) {
	if stack.IsEmpty() {
		var result T
		return result, errors.New("Cannot invoke Peek() on empty Stack.")
	}
	return stack.values[len(stack.values) - 1], nil
}

func (stack *Stack[T]) Length() int {
	if stack.values == nil {
		return 0
	}

	return len(stack.values)
}

func (stack *Stack[T]) IsEmpty() bool {
	return stack.values == nil || len(stack.values) == 0
}
