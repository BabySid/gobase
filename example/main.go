package main

import (
	"log/slog"

	"github.com/BabySid/gobase"
	"github.com/BabySid/gobase/log"
	"github.com/BabySid/gobase/log_sub"
)

func main() {
	l := log.NewSLogger(log.WithOutFile(log.StdOut), log.WithErrFile("./err.log"), log.WithJsonFormat())
	con, err := log_sub.NewConsumer(log_sub.Config{
		Location: &log_sub.SeekInfo{
			FileName: "20240314.log",
			Offset:   0,
			Whence:   0,
		},
		DateTimeLogLayout: &log_sub.DateTimeLayout{Layout: "20060102.log"},
	})

	gobase.True(err == nil)

	run(l)

	i := 0
	for {
		line := <-con.Lines
		if line.Err != nil {
			l.Info("got an error from consumeLog", slog.Any("err", err))
			continue
		}
		v := i % 2
		l.Info("got a line", slog.String("line", line.Text), slog.Int("counter", v))
		if v == 0 {
			l.Trace("there is a trace here, but donot print", slog.Int("counter", i))
			l.SetLevel(log.LevelTrace)
		}
		if v == 1 {
			l.Warn("this is a warn msg")
		}
		l.Warn("this is a warn msg in new line", slog.Int("counterInWarn", i))
		l.Info("this is a info msg in new line", slog.Int("counterInInfo", i))
		i += 1
	}
}

func run(l *log.SLogger) {
	l.Info("this is in run()")
	l.Info("this is in run", slog.String("key", "value"))
}
