package logger

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg any)
	Debugw(msg string, args ...any)
	Debugf(msg string, args ...any)
	Info(msg any)
	Infow(msg string, args ...any)
	InfowContext(ctx context.Context, msg string, args ...any)
	Infof(msg string, args ...any)
	InfofContext(ctx context.Context, msg string, args ...any)
	Error(msg any)
	Errorw(msg string, args ...any)
	Errorf(msg string, args ...any)
	ErrorfContext(ctx context.Context, msg string, args ...any)
	Fatal(msg any)
	Fatalw(msg string, args ...any)
	Fatalf(msg string, args ...any)
}

// Logger -.
type logger struct {
	log     *zap.SugaredLogger
	otel    *otelzap.SugaredLogger
	Tracing bool
}

var _ Logger = (*logger)(nil)

// New -.
func New(level string, opts ...Option) Logger {
	l := &logger{}

	for _, opt := range opts {
		opt(l)
	}

	cfg := zap.NewProductionConfig()

	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	l.log = log.Sugar()

	if l.Tracing {
		l.otel = otelzap.New(log, otelzap.WithMinLevel(zap.DebugLevel)).Sugar()
	}

	return l
}

func (l *logger) Debug(msg any) {
	if l.log == nil {
		l.otel.Debug(msg)
	}
}

func (l *logger) Debugw(msg string, args ...any) {
	l.log.Debugw(msg)
}

func (l *logger) Debugf(msg string, args ...any) {
	l.log.Debugf(msg, args...)
}

func (l *logger) Info(msg any) {
	l.log.Info(msg)
}

func (l *logger) Infow(msg string, args ...any) {
	l.log.Infow(msg, args...)
}

func (l *logger) InfowContext(ctx context.Context, msg string, args ...any) {
	if l.Tracing {
		l.otel.InfowContext(ctx, msg, args...)
	}
}

func (l *logger) Infof(msg string, args ...any) {
	l.log.Infof(msg, args...)
}

func (l *logger) InfofContext(ctx context.Context, msg string, args ...any) {
	if l.Tracing {
		l.otel.InfofContext(ctx, msg, args...)
	}
}

func (l *logger) Error(msg any) {
	l.log.Error(msg)
}

func (l *logger) Errorw(msg string, args ...any) {
	l.log.Errorw(msg, args...)
}

func (l *logger) Errorf(msg string, args ...any) {
	l.log.Errorf(msg, args...)
}

func (l *logger) ErrorfContext(ctx context.Context, msg string, args ...any) {
	if l.Tracing {
		l.otel.ErrorfContext(ctx, msg, args...)
	}
}

func (l *logger) Fatal(msg any) {
	l.log.Fatal(msg)
}

func (l *logger) Fatalw(msg string, args ...any) {
	l.log.Fatalw(msg, args...)
}

func (l *logger) Fatalf(msg string, args ...any) {
	l.log.Fatalf(msg, args...)
}
