package main

import (
	"fmt"
	"github.com/BabySid/gobase"
)

func main() {
	file, fun, line := gobase.GetFileInfoOfCaller(gobase.CurCaller)
	fmt.Println(file, gobase.GetShortFuncName(fun), line)
}
func run() {
	file, fun, line := gobase.GetFileInfoOfCaller(gobase.CurCaller)
	fmt.Println(file, gobase.GetShortFuncName(fun), line)
}
