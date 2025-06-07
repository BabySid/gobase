package main

import (
	"fmt"
	"log/slog"

	"github.com/BabySid/gobase"
	"github.com/BabySid/gobase/log"
)

var l *log.SLogger

func main() {
	l = log.NewSLogger(log.WithOutFile(log.StdOut), log.WithSkipCaller(1), log.WithColorful())
	l.SetLevel(log.LevelTrace)

	// info("this is a info log", slog.Int("number", 100), slog.Bool("flag", false))
	nl := l.WithOut(slog.Int("withNum", 200), slog.String("withStr", "hello")).WithErr(slog.Int("withNum", 300), slog.String("withStr", "helloworld"))
	nl.Trace("this is a new info log with attrs")
	nl.Debug("this is a new info log with attrs")
	nl.Info("this is a new info log with attrs")
	nl.Warn("this is a new info log with attrs")
	nl.Error("this is a new warn log with attrs")

	// l.Warn("this is a origin log")
	run()

	info("this is a msg in nested func")
	out := gobase.Combinations([]int{1, 3, 5, 7}, 2)
	fmt.Println(out)
}

func run() {
	// l.Info("this is in run()")
	info("this is in run")
}

func info(msg string) {
	l.Info(msg)
}
