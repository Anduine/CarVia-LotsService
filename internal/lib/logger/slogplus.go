package logger

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"

	"github.com/fatih/color"
)

type PlusHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PlusHandler struct {
	slog.Handler
	log   *stdLog.Logger
	attrs []slog.Attr
}

func InitGlobalLogger(writer io.Writer, level slog.Level) {
	opts := PlusHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}

	handler := opts.NewPlusHandler(writer)
	logger := slog.New(handler)

	slog.SetDefault(logger)
}

func (opts PlusHandlerOptions) NewPlusHandler(out io.Writer) *PlusHandler {
	h := &PlusHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		log:     stdLog.New(out, "", 0),
	}

	return h
}

func (h *PlusHandler) Handle(_ context.Context, rec slog.Record) error {
	level := rec.Level.String() + ":"

	switch rec.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]any, rec.NumAttrs())

	rec.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := rec.Time.Format("15:04:05.000")
	msg := color.CyanString(rec.Message)

	h.log.Println(
		timeStr,
		">",
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (h *PlusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PlusHandler{
		Handler: h.Handler,
		log:     h.log,
		attrs:   attrs,
	}
}

func (h *PlusHandler) WithGroup(name string) slog.Handler {
	return &PlusHandler{
		Handler: h.Handler.WithGroup(name),
		log:     h.log,
	}
}
