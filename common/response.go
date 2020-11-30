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
func JsonSuccess(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, ToSuccessMsg(msg, data))
}

// ToSuccessMsg returns a preset success content with gin.H
func ToSuccessMsg(msg string, data interface{}) gin.H {
	return gin.H{
		"code": http.StatusOK,
		"msg":  msg,
		"data": data,
	}
}

// JsonFail serializes the given struct as JSON into the response body.
func JsonFail(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, ToFailMsg(msg, data))
}

// JsonFailWithStack build a stack of given error and will be logged by Logger
// then call JsonFail
func JsonFailWithStack(ctx *gin.Context, err error, data interface{}) {
	stack := internal.BuildStack(err, 0)
	// log action should be asyncronous
	go Logger.Fields(map[string]interface{}{"stack": stack}).Error(ctx, err)
	// span
	if cspan, ok := ctx.Get("tracing-context"); ok {
		if span, ok := cspan.(opentracing.Span); ok {
			span.LogFields(log.Error(err))
			if stackJson, err := json.Marshal(map[string]interface{}{"stack": stack}); err == nil {
				span.LogFields(log.String("stack", string(stackJson)))
			}
		}
	}
	JsonFail(ctx, err.Error(), data)
}

// ToFailMsg returns a preset fail content with gin.H
func ToFailMsg(msg string, data interface{}) gin.H {
	return gin.H{
		"code": http.StatusInternalServerError,
		"msg":  msg,
		"data": data,
	}
}
