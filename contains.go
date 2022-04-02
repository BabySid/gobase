package gobase

func ContainsString(array []string, val string) int {
	for i, item := range array {
		if item == val {
			return i
		}
	}

	return -1
}

func ContainsInterface(array []interface{}, val interface{}) int {
	for i, item := range array {
		if item == val {
			return i
		}
	}

	return -1
}

func ContainsInt64(array []int64, val int64) int {
	for i, item := range array {
		if item == val {
			return i
		}
	}

	return -1
}
