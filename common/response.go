package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/yuchanns/bullets/internal"
	"net/http"
)

// JsonSuccess serializes the given struct as JSON into the response body.
func JsonSuccess(ctx *gin.Context, msg string, data interface{}, codes ...int) {
	// log action should be asyncronous
	go Logger.Fields(map[string]interface{}{"data": data}).Info(ctx)
	// span
	if cspan, ok := ctx.Get("tracing-context"); ok {
		if span, ok := cspan.(opentracing.Span); ok {
			if dataJson, err := json.Marshal(data); err == nil {
				span.LogFields(log.String("data", string(dataJson)))
			}
		}
	}
	ctx.JSON(http.StatusOK, ToSuccessMsg(msg, data, codes...))
}

// ToSuccessMsg returns a preset success content with gin.H
func ToSuccessMsg(msg string, data interface{}, codes ...int) gin.H {
	code := http.StatusOK
	if len(codes) > 0 {
		code = codes[0]
	}
	return gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	}
}

// JsonFail serializes the given struct as JSON into the response body.
func JsonFail(ctx *gin.Context, msg string, data interface{}, codes ...int) {
	ctx.JSON(http.StatusOK, ToFailMsg(msg, data, codes...))
}

// JsonFailWithStack build a stack of given error and will be logged by Logger
// then call JsonFail
func JsonFailWithStack(ctx *gin.Context, err error, data interface{}, codes ...int) {
	stack := internal.BuildStack(err, 0)
	// log action should be asyncronous
	go Logger.Fields(map[string]interface{}{"stack": stack}).Error(ctx, err)
	// span
	if cspan, ok := ctx.Get("tracing-context"); ok {
		if span, ok := cspan.(opentracing.Span); ok {
			span.SetTag("error", true)
			span.LogFields(log.Error(err))
			if stackJson, err := json.Marshal(map[string]interface{}{"stack": stack}); err == nil {
				span.LogFields(log.String("stack", string(stackJson)))
			}
		}
	}
	JsonFail(ctx, err.Error(), data, codes...)
}

// ToFailMsg returns a preset fail content with gin.H
func ToFailMsg(msg string, data interface{}, codes ...int) gin.H {
	code := http.StatusInternalServerError
	if len(codes) > 0 {
		code = codes[0]
	}
	return gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	}
}
