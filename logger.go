package gobase

import (
	"fmt"
	"os"
	"strings"
)

type Logger interface {
	Tracef(format string, args ...interface{})
	Trace(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
}

var _ Logger = (*StdErrLogger)(nil)

type StdErrLogger struct {
}

func (s *StdErrLogger) Tracef(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
}

func (s *StdErrLogger) Trace(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func (s *StdErrLogger) Debugf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
}

func (s *StdErrLogger) Debug(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func (s *StdErrLogger) Infof(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
}

func (s *StdErrLogger) Info(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func (s *StdErrLogger) Warnf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
}

func (s *StdErrLogger) Warn(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func (s *StdErrLogger) Fatalf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func (s *StdErrLogger) Fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}
