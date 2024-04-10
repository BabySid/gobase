package gobase

import (
	"path/filepath"
	"runtime"
	"strings"
)

var (
	CurCaller         = 1
	DefaultSkipCaller = 2
	DefaultMaxCaller  = 15
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

type CallFrame struct {
	Line     int
	Function string
	File     string
}

func GetCallerFrames(max int, skip int, fullPath bool) []CallFrame {
	var frames []CallFrame
	pcs := make([]uintptr, max)
	depth := runtime.Callers(skip, pcs)
	fs := runtime.CallersFrames(pcs[:depth])

	for f, again := fs.Next(); again; f, again = fs.Next() {
		cf := CallFrame{
			Line:     f.Line,
			Function: f.Function,
			File:     f.File,
		}
		if !fullPath {
			cf.File = filepath.Base(cf.File)
		}
		frames = append(frames, cf)
	}
	return frames
}
