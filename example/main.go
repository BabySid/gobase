package main

import (
	"fmt"
	"github.com/BabySid/gobase"
)

func main() {
	arr := []interface{}{"1", "2", "3", "4", "5"}

	rs := gobase.RemoveAnyFromSlice(arr, "3")
	fmt.Println(rs)
}
