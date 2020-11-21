package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuchanns/bullets/common"
	"github.com/yuchanns/bullets/internal"
	"io/ioutil"
	"net/url"
	"runtime"
	"strings"
)

func NewDefaultRequestInterceptor() gin.HandlerFunc {
	return NewRequestInterceptor(common.Logger)
}

func NewRequestInterceptor(logger internal.ILogger) gin.HandlerFunc {
	internal.Logger = logger
	return func(ctx *gin.Context) {
		// get a new copy of Header
		headerCopy := ctx.Request.Header.Clone()
		headerBuffer := new(bytes.Buffer)
		if err := headerCopy.Write(headerBuffer); err == nil {
			s := strings.Split(headerBuffer.String(), "\r\n")
			v := make(map[string]interface{})
			for i := range s {
				ss := strings.Split(strings.TrimSpace(s[i]), ": ")
				if len(ss) == 2 {
					v[ss[0]] = ss[1]
				}
			}
			internal.Logger.Fields(map[string]interface{}{"data": v}).DebugInfo(ctx, "Header")
		}
		if v, err := url.ParseQuery(ctx.Request.URL.RawQuery); err == nil {
			internal.Logger.Fields(map[string]interface{}{"data": v}).DebugInfo(ctx, "Query")
		}
		// get a new copy of Body
		bodyCopy, _ := ctx.Request.GetBody()
		defer bodyCopy.Close()
		if bodyBuf, err := ioutil.ReadAll(bodyCopy); err == nil {
			var v map[string]interface{}
			if err := json.Unmarshal(bodyBuf, &v); err == nil {
				internal.Logger.Fields(map[string]interface{}{"data": v}).DebugInfo(ctx, "Body")
			}
		}
		ctx.Next()
	}
}

// NewDefaultPanicInterceptor returns a gin middleware with a internal.BuiltinLogger
func NewDefaultPanicInterceptor() gin.HandlerFunc {
	return NewPanicInterceptor(common.Logger)
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
				internal.Logger.Fields(map[string]interface{}{"stack": stack}).Error(ctx, stackErr)
				common.JsonFail(ctx, stackErr.Error(), nil)
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
