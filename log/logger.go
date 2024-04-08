package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/BabySid/gobase"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type Logger interface {
	Trace(msg string, attrs ...slog.Attr)
	Debug(msg string, attrs ...slog.Attr)
	Info(msg string, attrs ...slog.Attr)
	Warn(msg string, attrs ...slog.Attr)
	Error(msg string, attrs ...slog.Attr)
	SetLevel(level slog.Level)

	WithOut(attrs ...slog.Attr) Logger
	WithErr(attrs ...slog.Attr) Logger
}

var _ Logger = (*SLogger)(nil)

const (
	LevelTrace = slog.Level(-8)
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

var levelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
}

type rotateByTime struct {
	pattern    string
	maxAge     time.Duration
	rotateTime time.Duration
}

type option struct {
	outFile string
	errFile string

	level *slog.LevelVar

	skipCaller int

	json         bool
	rotateByTime *rotateByTime
}

func defaultRotateByTime() *rotateByTime {
	return &rotateByTime{
		pattern:    HourPattern,
		maxAge:     time.Hour * 24 * 7,
		rotateTime: time.Hour,
	}
}

type Option func(*option)

const (
	StdOut = "stdout"
	StdErr = "stderr"

	HourPattern = ".%Y%m%d%H"
)

func WithOutFile(out string) Option {
	return func(opt *option) {
		opt.outFile = out
	}
}

func WithErrFile(out string) Option {
	return func(opt *option) {
		opt.errFile = out
	}
}

func WithJsonFormat() Option {
	return func(opt *option) {
		opt.json = true
	}
}

func WithSkipCaller(skip int) Option {
	return func(opt *option) {
		opt.skipCaller = skip
	}
}

func WithTimeRotate(pattern string, maxAge time.Duration, rotateTime time.Duration) Option {
	return func(opt *option) {
		opt.rotateByTime = &rotateByTime{}
		opt.rotateByTime.pattern = pattern
		opt.rotateByTime.maxAge = maxAge
		opt.rotateByTime.rotateTime = rotateTime
	}
}

func WithLevel(lvl string) Option {
	return func(opt *option) {
		level := strings.ToLower(lvl)
		switch level {
		case "trace":
			opt.level.Set(LevelTrace)
		case "debug":
			opt.level.Set(slog.LevelDebug)
		case "info":
			opt.level.Set(slog.LevelInfo)
		case "warn":
			opt.level.Set(slog.LevelWarn)
		case "error":
			opt.level.Set(slog.LevelError)
		default:
			panic(fmt.Errorf("invalid level:%s", lvl))
		}
	}
}

type SLogger struct {
	opt option

	outWriter io.Writer
	errWriter io.Writer

	slogOpt slog.HandlerOptions
	out     *slog.Logger
	err     *slog.Logger
}

func NewSLogger(opts ...Option) *SLogger {
	log := SLogger{
		opt: option{
			level: &slog.LevelVar{},
		},
	}

	for _, opt := range opts {
		opt(&log.opt)
	}

	log.outWriter = log.getWriter(log.opt.outFile)
	log.errWriter = log.getWriter(log.opt.errFile)
	gobase.TrueF(log.outWriter != nil || log.errWriter != nil, "outFile or errFile must be set at least one")

	log.outWriter = gobase.GetNotNil(log.outWriter, log.errWriter)
	log.errWriter = gobase.GetNotNil(log.errWriter, log.outWriter)
	gobase.TrueF(log.outWriter != nil && log.errWriter != nil, "log.outWriter=%v log.errWriter=%v", log.outWriter, log.errWriter)

	log.slogOpt = slog.HandlerOptions{
		AddSource: true,
		Level:     log.opt.level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := levelNames[level]
				if !exists {
					levelLabel = level.String()
				}
				a.Value = slog.StringValue(levelLabel)
			}
			return a
		},
	}

	log.out = log.getSlogger(log.outWriter)
	log.err = log.getSlogger(log.errWriter)

	gobase.TrueF(log.out != nil && log.err != nil, "log.out=%v log.err=%v", log.out, log.err)

	return &log
}

func (d *SLogger) getSlogger(out io.Writer, attrs ...slog.Attr) *slog.Logger {
	if out != nil {
		var handler slog.Handler = newLogHandler(d.opt.json, d.opt.skipCaller, out, &d.slogOpt)
		if len(attrs) > 0 {
			handler = handler.WithAttrs(attrs)
		}

		return slog.New(handler)
	}
	return nil
}

func (d *SLogger) getWriter(file string) io.Writer {
	switch file {
	case StdOut:
		return os.Stdout
	case StdErr:
		return os.Stderr
	default:
		if file != "" {
			if d.opt.rotateByTime == nil {
				d.opt.rotateByTime = defaultRotateByTime()
			}
			out, err := rotatelogs.New(
				file+d.opt.rotateByTime.pattern,
				rotatelogs.WithLinkName(file),
				rotatelogs.WithMaxAge(d.opt.rotateByTime.maxAge),
				rotatelogs.WithRotationTime(d.opt.rotateByTime.rotateTime),
			)
			gobase.TrueF(err == nil, "init slog failed. err=%v", err)
			return out
		}
	}
	return nil
}

// Trace implements Logger.
func (d *SLogger) Trace(msg string, attrs ...slog.Attr) {
	d.out.LogAttrs(context.Background(), LevelTrace, msg, attrs...)
}

// Debug implements Logger.
func (d *SLogger) Debug(msg string, attrs ...slog.Attr) {
	d.out.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
}

// Info implements Logger.
func (d *SLogger) Info(msg string, attrs ...slog.Attr) {
	d.out.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
}

// Warn implements Logger.
func (d *SLogger) Warn(msg string, attrs ...slog.Attr) {
	d.err.LogAttrs(context.Background(), slog.LevelWarn, msg, attrs...)
}

// Error implements Logger.
func (d *SLogger) Error(msg string, attrs ...slog.Attr) {
	d.err.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
}

// SetLevel implements Logger.
func (d *SLogger) SetLevel(level slog.Level) {
	d.opt.level.Set(level)
}

func (d *SLogger) WithOut(attrs ...slog.Attr) Logger {
	if len(attrs) == 0 {
		return d
	}

	n := d.clone()
	n.out = n.getSlogger(n.outWriter, attrs...)
	gobase.TrueF(n.out != nil && n.err != nil, "log.out=%v log.err=%v", n.out, n.err)
	return n
}

func (d *SLogger) WithErr(attrs ...slog.Attr) Logger {
	if len(attrs) == 0 {
		return d
	}

	n := d.clone()
	n.err = n.getSlogger(n.errWriter, attrs...)
	gobase.TrueF(n.out != nil && n.err != nil, "log.out=%v log.err=%v", n.out, n.err)
	return n
}

func (d *SLogger) clone() *SLogger {
	n := *d
	return &n
}
