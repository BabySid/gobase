package main

import (
	"log/slog"

	"github.com/BabySid/gobase/log"
)

var l *log.SLogger

func main() {
	l = log.NewSLogger(log.WithOutFile(log.StdOut), log.WithSkipCaller(1))

	// info("this is a info log", slog.Int("number", 100), slog.Bool("flag", false))
	nl := l.WithOut(slog.Int("withNum", 200), slog.String("withStr", "hello")).WithErr(slog.Int("withNum", 300), slog.String("withStr", "helloworld"))
	nl.Info("this is a new info log with attrs")
	nl.Warn("this is a new warn log with attrs")

	// l.Warn("this is a origin log")
	run()

	info("this is a msg in nested func")
}

func run() {
	// l.Info("this is in run()")
	info("this is in run")
}

func info(msg string) {
	l.Info(msg)
}
