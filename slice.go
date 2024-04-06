package gobase

import "reflect"

func RemoveItemFromSlice[T comparable](src []T, elem T) []T {
	i := 0
	for _, v := range src {
		if v != elem {
			src[i] = v
			i++
		}
	}
	return src[:i]
}

func RemoveAnyFromSlice(src []interface{}, elem interface{}) []interface{} {
	i := 0
	for _, v := range src {
		if v != elem {
			src[i] = v
			i++
		}
	}
	return src[:i]
}

func ConvertSlice[T any, R any](s []T, f func(T) R) []R {
	result := make([]R, len(s), len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

func GetNotNil[T comparable](s ...T) T {
	for _, v := range s {
		if !isNil(v) {
			return v
		}
	}
	var zero T
	return zero
}

func isNil[T any](t T) bool {
	v := reflect.ValueOf(t)
	kind := v.Kind()
	return (kind == reflect.Ptr || kind == reflect.Interface || kind == reflect.Slice || kind == reflect.Map || kind == reflect.Chan) && v.IsNil()
}
