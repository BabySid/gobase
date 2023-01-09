package gobase

import (
	"runtime"
	"strings"
)

var (
	CurCaller = 1
)

func GetFileInfoOfCaller(skip int) (string, string, int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "???", "???", 0
	}

	return file, runtime.FuncForPC(pc).Name(), line
}

func GetShortFuncName(long string) string {
	arr := strings.Split(long, ".")
	return arr[len(arr)-1]
}
