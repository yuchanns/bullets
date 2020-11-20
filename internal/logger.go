package internal

import "context"

type ILogger interface {
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, err error, args ...interface{})
	Errorf(ctx context.Context, err error, format string, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	DebugInfo(ctx context.Context, args ...interface{})
	DebugInfof(ctx context.Context, format string, args ...interface{})
}

var Logger ILogger
