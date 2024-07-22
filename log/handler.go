package log

import (
	"context"
	"io"
	"log/slog"
	"runtime"
)

var _ slog.Handler = (*logHandler)(nil)

type logHandler struct {
	slog.Handler
	out  io.Writer
	opt  *slog.HandlerOptions
	json bool
	skip int
}

func newLogHandler(json bool, skip int, out io.Writer, opt *slog.HandlerOptions) *logHandler {
	var handler slog.Handler
	if json {
		handler = slog.NewJSONHandler(out, opt)
	} else {
		handler = slog.NewTextHandler(out, opt)
	}

	return &logHandler{
		Handler: handler,
		out:     out,
		opt:     opt,
		json:    json,
		skip:    skip,
	}
}

func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.PC != 0 {
		// fs := runtime.CallersFrames([]uintptr{r.PC})
		// f, _ := fs.Next()

		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller, slog.LogAttrs function, slog.LogAttrs's caller]
		runtime.Callers(5+h.skip, pcs[:])
		r.PC = pcs[0]
	}

	return h.Handler.Handle(ctx, r)
}

func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handle := h.Handler.WithAttrs(attrs)
	return &logHandler{
		Handler: handle,
		out:     h.out,
		opt:     h.opt,
		json:    h.json,
	}
}
