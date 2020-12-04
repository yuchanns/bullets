package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuchanns/bullets/common"
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

func TestBuildOpenTracerInterceptor(t *testing.T) {
	closeOpenTracerFunc, openTracerMiddleware, err := BuildOpenTracerInterceptor("testOpenTrace", os.Getenv("COLLECTOR_HOST"), []byte("api-request-"))
	if err != nil {
		panic(err)
	}
	defer closeOpenTracerFunc()
	engine := gin.New()
	engine.Use(openTracerMiddleware)
	engine.POST("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "")
	})
	engine.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "")
	})
	engine.GET("/panic", func(ctx *gin.Context) {
		var err error
		common.JsonFail(ctx, err.Error(), nil)
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
	req3, _ := http.NewRequest(http.MethodGet, "/panic", nil)
	req4, _ := http.NewRequest(http.MethodGet, "/regular_err", nil)
	engine.ServeHTTP(w, req)
	engine.ServeHTTP(w, req2)
	engine.ServeHTTP(w, req3)
	engine.ServeHTTP(w, req4)
}
