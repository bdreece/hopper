package iter

type Iterator[T any] interface {
	Next() (*T, error)
}

type MapFunc[T, U any] func(*T) U
type Map[T, U any] struct {
	iter Iterator[T]
	fn   MapFunc[T, U]
}

func Collect[T any](iter Iterator[T]) []T {
	var (
		err  error
		vals []T = make([]T, 1)
	)

	for err == nil {
		next, err := iter.Next()
		if err == nil {
			vals = append(vals, *next)
		}
	}

	return vals
}

func NewMap[T, U any](iter Iterator[T], fn MapFunc[T, U]) Iterator[U] {
	return &Map[T, U]{
		iter,
		fn,
	}
}

func (m *Map[T, U]) Next() (*U, error) {
	next, err := m.iter.Next()
	if err != nil {
		return nil, err
	}
	val := m.fn(next)
	return &val, nil
}
