package log

import (
	"context"
	"log/slog"
	"runtime"
)

var _ slog.Handler = (*textHandler)(nil)
var _ slog.Handler = (*jsonHandler)(nil)

type textHandler struct {
	*slog.TextHandler
}

func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.PC != 0 {
		// fs := runtime.CallersFrames([]uintptr{r.PC})
		// f, _ := fs.Next()

		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller, slog.LogAttrs function, slog.LogAttrs's caller]
		runtime.Callers(5, pcs[:])
		r.PC = pcs[0]
	}

	return h.TextHandler.Handle(ctx, r)
}

type jsonHandler struct {
	*slog.JSONHandler
}

func (h *jsonHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.PC != 0 {
		// fs := runtime.CallersFrames([]uintptr{r.PC})
		// f, _ := fs.Next()

		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller, slog.LogAttrs function, slog.LogAttrs's caller]
		runtime.Callers(5, pcs[:])
		r.PC = pcs[0]
	}
	return h.JSONHandler.Handle(ctx, r)
}
