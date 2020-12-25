package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	ginopentracing "github.com/Bose/go-gin-opentracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	transport2 "github.com/uber/jaeger-client-go/transport"
	"github.com/yuchanns/bullets/common"
	"github.com/yuchanns/bullets/internal"
	"io"
	"strings"
)

func openTracer(operationPrefix []byte) gin.HandlerFunc {
	if operationPrefix == nil {
		operationPrefix = []byte("api-request-")
	}
	return func(c *gin.Context) {
		var span opentracing.Span
		if cspan, ok := c.Get("tracing-context"); ok {
			span = ginopentracing.StartSpanWithParent(cspan.(opentracing.Span).Context(), c.Request.Method+":"+c.Request.URL.Path, c.Request.Method, c.Request.URL.Path)

		} else {
			span = ginopentracing.StartSpanWithHeader(&c.Request.Header, c.Request.Method+":"+c.Request.URL.Path, c.Request.Method, c.Request.URL.Path)
		}
		defer span.Finish()
		defer func() {
			if msg := recover(); msg != nil {
				span.SetTag("error", true)
				stack, stackErr := buildStackFromRecover(msg)
				if stackErr != nil {
					span.LogFields(log.Error(stackErr))
				}
				if stackJson, err := json.Marshal(map[string]interface{}{"stack": stack}); err == nil {
					span.LogFields(log.String("stack", string(stackJson)))
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

// InitTracing - init opentracing with use http
func InitTracing(serviceName string, url string) (
	tracer opentracing.Tracer,
	reporter jaeger.Reporter,
	closer io.Closer,
	err error) {
	transport := transport2.NewHTTPTransport(url)
	reporter = jaeger.NewRemoteReporter(transport)

	var sampler jaeger.Sampler
	sampler = jaeger.NewConstSampler(true)

	tracer, closer = jaeger.NewTracer(serviceName,
		sampler,
		reporter,
	)
	return tracer, reporter, closer, nil
}

// BuildOpenTracerInterceptor is an alias for BuildOpenTracerAgentInterceptor
func BuildOpenTracerInterceptor(
	serviceName, agentHostPort string,
	operationPrefix []byte,
) (
	closeFunc func(),
	middleware gin.HandlerFunc,
	err error,
) {
	return BuildOpenTracerAgentInterceptor(serviceName, agentHostPort, operationPrefix)
}

// buildWithInitTracing is the common part of BuildOpenTracerCollectorInterceptor and BuildOpenTracerAgentInterceptor
func buildWithInitTracing(tracer opentracing.Tracer, reporter jaeger.Reporter, closer io.Closer, operationPrefix []byte) (
	closeFunc func(),
	middleware gin.HandlerFunc,
	err error,
) {
	internal.Logger = common.Logger
	closeFunc = func() {
		reporter.Close()
		closer.Close()
	}
	opentracing.SetGlobalTracer(tracer)
	middleware = openTracer(operationPrefix)
	return
}

// BuildOpenTracerCollectorInterceptor create an interceptor using the jaeger collector directly
func BuildOpenTracerCollectorInterceptor(
	serviceName, collectorHost string,
	operationPrefix []byte,
) (
	closeFunc func(),
	middleware gin.HandlerFunc,
	err error,
) {
	url := strings.Join([]string{collectorHost, "api/traces?format=jaeger.thrift"}, "/")
	tracer, reporter, closer, err := InitTracing(serviceName, url)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("unable to init collector tracing:%s", err))
	}
	return buildWithInitTracing(tracer, reporter, closer, operationPrefix)
}

// BuildOpenTracerAgentInterceptor create an interceptor using the jaeger agent
func BuildOpenTracerAgentInterceptor(
	serviceName, agentHostPort string,
	operationPrefix []byte,
) (
	closeFunc func(),
	middleware gin.HandlerFunc,
	err error,
) {
	tracer, reporter, closer, err := ginopentracing.InitTracing(serviceName, agentHostPort, ginopentracing.WithEnableInfoLog(true))
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("unable to init agent tracing:%s", err))
	}
	return buildWithInitTracing(tracer, reporter, closer, operationPrefix)
}
