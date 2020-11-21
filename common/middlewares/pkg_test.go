package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/yuchanns/bullet/common"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/test?c=d&e=f", bytes.NewBuffer([]byte("{\"a\":\"b\"}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "111111")
	engine.ServeHTTP(w, req)
}
