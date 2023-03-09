package gobase

import (
	"golang.org/x/exp/constraints"
	"sort"
)

type Map[T constraints.Ordered, V any] struct {
	Key   T
	Value V
}
type MapSlice[T constraints.Ordered, V any] []Map[T, V]

func SortMap[T constraints.Ordered, V any](m map[T]V) MapSlice[T, V] {
	keys := make(orderedSlice[T], 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}

	sort.Sort(keys)

	ms := make(MapSlice[T, V], len(keys))
	for i, k := range keys {
		ms[i] = Map[T, V]{Key: k, Value: m[k]}
	}

	return ms
}

type orderedSlice[T constraints.Ordered] []T

func (cs orderedSlice[T]) Less(i, j int) bool {
	return cs[i] < cs[j]
}

func (cs orderedSlice[T]) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

func (cs orderedSlice[T]) Len() int {
	return len(cs)
}
