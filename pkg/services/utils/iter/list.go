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
