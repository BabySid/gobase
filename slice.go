package gobase

func RemoveItemFromSlice(src []interface{}, elem interface{}) []interface{} {
	i := 0
	for _, v := range src {
		if v != elem {
			src[i] = v
			i++
		}
	}
	return src[:i]
}
