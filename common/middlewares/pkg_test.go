package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuchanns/bullets/common"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNewDefaultPanicInterceptor(t *testing.T) {
	engine := gin.New()
	engine.Use(NewDefaultPanicInterceptor())
	engine.GET("/panic", func(ctx *gin.Context) {
		var err error
		common.JsonFail(ctx, err.Error(), nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/panic", nil)
	engine.ServeHTTP(w, req)
	if w.Body.String() != "{\"code\":500,\"data\":null,\"message\":\"panic runtime error: runtime error: invalid memory address or nil pointer dereference: runtime error: invalid memory address or nil pointer dereference\"}" {
		log.Fatal("return response is not equal as expect")
	}
}

func TestNewDefaultRequestInterceptor(t *testing.T) {
	engine := gin.New()
	engine.Use(NewDefaultRequestInterceptor())
	engine.POST("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "")
	})
	engine.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "")
	})
	engine.GET("/regular_err", func(ctx *gin.Context) {
		common.JsonFailWithStack(ctx, errors.Errorf("a regular error"), nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/test?c=d&e=f", bytes.NewBuffer([]byte("{\"a\":\"b\"}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "111111")
	req2, _ := http.NewRequest(http.MethodGet, "/test?a=b&c=d", nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-User-Id", "111111")
	req3, _ := http.NewRequest(http.MethodGet, "/regular_err", nil)
	engine.ServeHTTP(w, req)
	engine.ServeHTTP(w, req2)
	engine.ServeHTTP(w, req3)
	time.Sleep(time.Second)
}

func testBuildOpenTracerInterceptor(t *testing.T, closeOpenTracerFunc func(), openTracerMiddleware gin.HandlerFunc) {
	defer closeOpenTracerFunc()
	engine := gin.New()
	engine.Use(openTracerMiddleware)
	engine.POST("/test", func(ctx *gin.Context) {
		common.JsonSuccess(ctx, "success", gin.H{"hello": "world"})
	})
	engine.GET("/test", func(ctx *gin.Context) {
		common.JsonSuccess(ctx, "success", gin.H{"hello": "world"})
	})
	engine.GET("/panic", func(ctx *gin.Context) {
		var err error
		common.JsonFail(ctx, err.Error(), nil)
	})
	engine.GET("/regular_err", func(ctx *gin.Context) {
		common.JsonFailWithStack(ctx, errors.Errorf("a regular error"), nil)
	})

	w := httptest.NewRecorder()
	reqTable := []func() *http.Request{
		func() *http.Request {
			req, _ := http.NewRequest(http.MethodPost, "/test?c=d&e=f", bytes.NewBuffer([]byte("{\"a\":\"b\"}")))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-Id", "111111")
			return req
		},
		func() *http.Request {
			req, _ := http.NewRequest(http.MethodGet, "/test?a=b&c=d", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-Id", "111111")
			return req
		},
		func() *http.Request {
			req, _ := http.NewRequest(http.MethodGet, "/panic", nil)
			return req
		},
		func() *http.Request {
			req, _ := http.NewRequest(http.MethodGet, "/regular_err", nil)
			return req
		},
	}
	for i := range reqTable {
		engine.ServeHTTP(w, reqTable[i]())
	}
}

func TestBuildOpenTracerCollectorInterceptor(t *testing.T) {
	closeOpenTracerFunc, openTracerMiddleware, err := BuildOpenTracerCollectorInterceptor("testOpenTraceCollector", os.Getenv("COLLECTOR_HOST"), []byte("api-request-"))
	if err != nil {
		panic(err)
	}
	testBuildOpenTracerInterceptor(t, closeOpenTracerFunc, openTracerMiddleware)
}

func TestBuildOpenTracerAgentInterceptor(t *testing.T) {
	closeOpenTracerFunc, openTracerMiddleware, err := BuildOpenTracerInterceptor("testOpenTraceAgent", os.Getenv("AGENT_HOST"), []byte("api-request-"))
	if err != nil {
		panic(err)
	}
	testBuildOpenTracerInterceptor(t, closeOpenTracerFunc, openTracerMiddleware)
}
