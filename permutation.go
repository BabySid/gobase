package gobase

func Combinations[T any](arr []T, m int) [][]T {
	n := len(arr)
	combinations := [][]T{}

	for i := 0; i < (1 << n); i++ {
		if countBits(i) == m {
			// 生成当前组合
			comb := []T{}
			for j := 0; j < n; j++ {
				if i&(1<<j) != 0 {
					comb = append(comb, arr[j])
				}
			}
			combinations = append(combinations, comb)
		}
	}

	return combinations
}

// countBits 计算二进制表示中1的个数
func countBits(x int) int {
	count := 0
	for x > 0 {
		count += x & 1
		x >>= 1
	}
	return count
}
