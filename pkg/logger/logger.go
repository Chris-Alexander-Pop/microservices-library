package logger

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"
)

var (
	defaultLogger *slog.Logger
	once          sync.Once
)

type Config struct {
	Level  string `env:"LOG_LEVEL" env-default:"INFO"`
	Format string `env:"LOG_FORMAT" env-default:"JSON"` // JSON or TEXT
}

// Init initializes the global logger
func Init(cfg Config) *slog.Logger {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: parseLevel(cfg.Level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// standard time format
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format(time.RFC3339))
			}
			return a
		},
	}

	if cfg.Format == "TEXT" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(NewTraceHandler(handler))
	slog.SetDefault(logger)

	once.Do(func() {
		defaultLogger = logger
	})

	return logger
}

// Global accessor, though we prefer passing logger or using FromContext if we attach it
func L() *slog.Logger {
	if defaultLogger == nil {
		// Fallback if not initialized
		return slog.Default()
	}
	return defaultLogger
}

func parseLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// TraceHandler adds trace_id and span_id to logs
type TraceHandler struct {
	next slog.Handler
}

func NewTraceHandler(next slog.Handler) *TraceHandler {
	return &TraceHandler{next: next}
}

func (h *TraceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return h.next.Handle(ctx, r)
}

func (h *TraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TraceHandler{next: h.next.WithAttrs(attrs)}
}

func (h *TraceHandler) WithGroup(name string) slog.Handler {
	return &TraceHandler{next: h.next.WithGroup(name)}
}
