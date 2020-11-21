package internal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

type BuiltinLogger struct {
	Zl     *zerolog.Logger
	fields map[string]interface{}
}

func NewBuiltinLogger() *BuiltinLogger {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &BuiltinLogger{
		Zl: &logger,
	}
}

func (l *BuiltinLogger) Fields(fields map[string]interface{}) ILogger {
	l.fields = fields
	return l
}

func (l *BuiltinLogger) Info(ctx context.Context, args ...interface{}) {
	var b bytes.Buffer
	_, _ = fmt.Fprint(&b, args...)
	l.Zl.Info().Fields(l.fields).Msg(b.String())
}

func (l *BuiltinLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	var b bytes.Buffer
	_, _ = fmt.Fprintf(&b, format, args...)
	l.Zl.Info().Fields(l.fields).Msg(b.String())
}

func (l *BuiltinLogger) Error(ctx context.Context, err error, args ...interface{}) {
	var b bytes.Buffer
	_, _ = fmt.Fprint(&b, args...)
	l.Zl.Error().Fields(l.fields).Err(err).Msg(b.String())
}

func (l *BuiltinLogger) Errorf(ctx context.Context, err error, format string, args ...interface{}) {
	var b bytes.Buffer
	_, _ = fmt.Fprintf(&b, format, args...)
	l.Zl.Error().Fields(l.fields).Err(err).Msg(b.String())
}

func (l *BuiltinLogger) Warn(ctx context.Context, args ...interface{}) {
	var b bytes.Buffer
	_, _ = fmt.Fprint(&b, args...)
	l.Zl.Warn().Fields(l.fields).Msg(b.String())
}

func (l *BuiltinLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	var b bytes.Buffer
	_, _ = fmt.Fprintf(&b, format, args...)
	l.Zl.Warn().Fields(l.fields).Msg(b.String())
}

func (l *BuiltinLogger) DebugInfo(ctx context.Context, args ...interface{}) {
	var b bytes.Buffer
	_, _ = fmt.Fprint(&b, args...)
	l.Zl.Debug().Fields(l.fields).Msg(b.String())
}

func (l *BuiltinLogger) DebugInfof(ctx context.Context, format string, args ...interface{}) {
	var b bytes.Buffer
	_, _ = fmt.Fprintf(&b, format, args...)
	l.Zl.Debug().Fields(l.fields).Msg(b.String())
}
