package middlewares

import (
	"encoding/json"
	ginopentracing "github.com/Bose/go-gin-opentracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/yuchanns/bullets/common"
	"github.com/yuchanns/bullets/internal"
)

func openTracer(operationPrefix []byte) gin.HandlerFunc {
	if operationPrefix == nil {
		operationPrefix = []byte("api-request-")
	}
	return func(c *gin.Context) {
		var span opentracing.Span
		if cspan, ok := c.Get("tracing-context"); ok {
			span = ginopentracing.StartSpanWithParent(cspan.(opentracing.Span).Context(), string(operationPrefix)+c.Request.Method, c.Request.Method, c.Request.URL.Path)

		} else {
			span = ginopentracing.StartSpanWithHeader(&c.Request.Header, string(operationPrefix)+c.Request.Method, c.Request.Method, c.Request.URL.Path)
		}
		defer span.Finish()
		defer func() {
			if msg := recover(); msg != nil {
				stack, stackErr := buildStackFromRecover(msg)
				if stackErr != nil {
					span.LogFields(log.Error(stackErr))
				}
				if stackJson, err := json.Marshal(map[string]interface{}{"stack": stack}); err == nil {
					// stop to report stack in case the limitation of udp max pack size
					_ = stackJson
					//span.LogFields(log.String("stack", string(stackJson)))
				}
				go internal.Logger.Fields(map[string]interface{}{"stack": stack}).Error(c, stackErr)
				common.JsonFail(c, stackErr.Error(), nil)
				c.Abort()
			}
		}()
		headerData, bodyData, queryData := findDataFromContext(c)
		if headerJson, err := json.Marshal(headerData); err == nil {
			span.LogFields(log.String("header", string(headerJson)))
		}
		if bodyDataJson, err := json.Marshal(bodyData); err == nil {
			span.LogFields(log.String("body", string(bodyDataJson)))
		}
		if queryDataJson, err := json.Marshal(queryData); err == nil {
			span.LogFields(log.String("params", string(queryDataJson)))
		}
		c.Set("tracing-context", span)
		c.Next()

		span.SetTag(string(ext.HTTPStatusCode), c.Writer.Status())
	}
}

func BuildOpenTracerInterceptor(
	serviceName, agentHostPort string,
	operationPrefix []byte,
) (
	closeFunc func(),
	middleware gin.HandlerFunc,
	err error,
) {
	tracer, reporter, closer, err := ginopentracing.InitTracing(serviceName, agentHostPort, ginopentracing.WithEnableInfoLog(true))
	if err != nil {
		panic("unable to init tracing")
	}
	closeFunc = func() {
		reporter.Close()
		closer.Close()
	}
	opentracing.SetGlobalTracer(tracer)
	middleware = openTracer(operationPrefix)
	return
}
