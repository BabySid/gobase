package gobase

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
