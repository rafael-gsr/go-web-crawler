// Package datastructure
package datastructure

import (
	"maps"
	"slices"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Set[T]) Remove(item T) {
	delete(s, item)
}

func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) GetAll() Set[T] {
	return s
}

func (s Set[T]) GetList() []T {
	keys := maps.Keys(s)
	list := []T{}
	keys(func(value T) bool {
		list = append(list, value)
		return true
	})

	return list
}

func (s Set[T]) GetFirst() T {
	keys := maps.Keys(s)
	keySlice := slices.Collect(keys)

	return keySlice[0]
}
