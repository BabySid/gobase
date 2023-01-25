package main

import (
	"fmt"
	"github.com/BabySid/gobase"
)

func main() {
	arr := []interface{}{1, 2, 3, 4, 5}

	rs := gobase.RemoveItemFromSlice(arr, 1)
	fmt.Println(rs)
}
