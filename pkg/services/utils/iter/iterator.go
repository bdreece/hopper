/*
 * hopper - A gRPC API for collecting IoT device event messages
 * Copyright (C) 2022 Brian Reece

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
		vals []T = make([]T, 0, 1)
	)

	for err == nil {
		next, err := iter.Next()
		if err != nil {
			break
		}

		vals = append(vals, *next)
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
