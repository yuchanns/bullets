package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuchanns/bullet/common"
	"github.com/yuchanns/bullet/internal"
	"runtime"
)

func NewDefaultRequestInterceptor() gin.HandlerFunc {
	return NewRequestInterceptor(common.DefaultLogger)
}

func NewRequestInterceptor(logger internal.ILogger) gin.HandlerFunc {
	internal.Logger = logger
	return func(ctx *gin.Context) {
		// TODO: log request data
		ctx.Next()
	}
}

// NewDefaultPanicInterceptor returns a gin middleware with a internal.BuiltinLogger
func NewDefaultPanicInterceptor() gin.HandlerFunc {
	return NewPanicInterceptor(common.DefaultLogger)
}

// NewPanicInterceptor returns a gin middleware
func NewPanicInterceptor(logger internal.ILogger) gin.HandlerFunc {
	internal.Logger = logger
	return func(ctx *gin.Context) {
		// recover from panic and record
		defer func() {
			if msg := recover(); msg != nil {
				var (
					stack    []string
					stackErr error
				)
				switch err := msg.(type) {
				case runtime.Error:
					stackErr = errors.Wrapf(err, "panic runtime error: %v", err)
					stack = internal.BuildStack(stackErr, 4)
				default:
					stackErr = errors.New(fmt.Sprintf("panic error: %v", err))
					stack = internal.BuildStack(stackErr, 4)
				}
				internal.Logger.Error(ctx, stackErr, stack)
				common.JsonFail(ctx, stackErr.Error(), nil)
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
