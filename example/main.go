package main

import (
	"fmt"
	"github.com/BabySid/gobase"
)

func main() {
	arr := []string{"1", "2", "3", "4", "5"}

	rs := gobase.RemoveItemFromSlice(arr, "3")
	fmt.Println(rs)
}
