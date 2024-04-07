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

	level slog.LevelVar

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

	out *slog.Logger
	err *slog.Logger
}

func NewSLogger(opts ...Option) *SLogger {
	log := SLogger{}

	for _, opt := range opts {
		opt(&log.opt)
	}

	log.outWriter = log.getWriter(log.opt.outFile)
	log.errWriter = log.getWriter(log.opt.errFile)

	gobase.TrueF(log.outWriter != nil || log.errWriter != nil, "outFile or errFile must be set at least one")

	slogOpt := slog.HandlerOptions{
		AddSource: true,
		Level:     &log.opt.level,
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

	log.out = log.getSlogger(log.outWriter, &slogOpt)
	log.err = log.getSlogger(log.errWriter, &slogOpt)

	log.out = gobase.GetNotNil(log.out, log.err)
	log.err = gobase.GetNotNil(log.err, log.out)
	gobase.TrueF(log.out != nil && log.err != nil, "log.out=%v log.err=%v", log.out, log.err)

	return &log
}

func (d *SLogger) getSlogger(out io.Writer, opt *slog.HandlerOptions) *slog.Logger {
	if out != nil {
		var handler slog.Handler = &textHandler{
			slog.NewTextHandler(out, opt),
		}
		if d.opt.json {
			handler = &jsonHandler{
				slog.NewJSONHandler(out, opt),
			}
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

func (d *SLogger) OutLogger() *slog.Logger {
	return d.out
}

func (d *SLogger) ErrLogger() *slog.Logger {
	return d.err
}
