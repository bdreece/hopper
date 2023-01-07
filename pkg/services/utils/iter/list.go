package iter

import "errors"

var (
	ErrIndexOutOfBounds = errors.New("Index out of bounds")
)

type ListIterator[T any] struct {
	list  *[]T
	index int
}

func NewListIterator[T any](size, capacity uint) Iterator[T] {
	list := make([]T, size, capacity)
	return &ListIterator[T]{
		list:  &list,
		index: 0,
	}
}

func FromSlice[T any](slice *[]T) Iterator[T] {
	return &ListIterator[T]{
		list:  slice,
		index: 0,
	}
}

func (l *ListIterator[T]) Next() (*T, error) {
	if l.index >= len(*l.list) {
		return nil, ErrIndexOutOfBounds
	}

	val := (*l.list)[l.index]
	l.index += 1
	return &val, nil
}
