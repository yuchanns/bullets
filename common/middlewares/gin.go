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
		headerData, bodyData, queryData := findDataFromContext(ctx)
		// log action should be asyncronous
		go func() {
			if headerData != nil {
				internal.Logger.Fields(map[string]interface{}{"data": headerData}).DebugInfo(ctx, "Header")
			}
			if queryData != nil {
				internal.Logger.Fields(map[string]interface{}{"data": queryData}).DebugInfo(ctx, "Query")
			}
			if bodyData != nil {
				internal.Logger.Fields(map[string]interface{}{"data": bodyData}).DebugInfo(ctx, "Body")
			}
		}()
		ctx.Next()
	}
}

// findDataFromContext find header, body, and query data from *gin.Context
func findDataFromContext(ctx *gin.Context) (
	headerData, bodyData map[string]interface{},
	queryData url.Values,
) {
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
		headerData = v
	}
	if v, err := url.ParseQuery(ctx.Request.URL.RawQuery); err == nil {
		queryData = v
	}
	if ctx.Request.Body != nil {
		if bodyBuf, err := ioutil.ReadAll(ctx.Request.Body); err == nil {
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuf))
			var v map[string]interface{}
			if err := json.Unmarshal(bodyBuf, &v); err == nil {
				bodyData = v
			}
		}
	}
	return
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
				stack, stackErr := buildStackFromRecover(msg)
				// log action should be asyncronous
				go internal.Logger.Fields(map[string]interface{}{"stack": stack}).Error(ctx, stackErr)
				common.JsonFail(ctx, stackErr.Error(), nil)
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

// buildStackFromRecover build error with stack from panic
func buildStackFromRecover(msg interface{}) (
	stack []string,
	stackErr error,
) {
	switch err := msg.(type) {
	case runtime.Error:
		stackErr = errors.Wrapf(err, "panic runtime error: %v", err)
		stack = internal.BuildStack(stackErr, 4)
	default:
		stackErr = errors.New(fmt.Sprintf("panic error: %v", err))
		stack = internal.BuildStack(stackErr, 4)
	}
	return
}
