package common

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func openNotExistFile() error {
	if _, err := os.Open("NOT_EXIST_FILE"); err != nil {
		return errors.Wrap(err, "open error")
	}
	return nil
}

func openNotExistFileHandler(ctx *gin.Context) {
	if err := openNotExistFile(); err != nil {
		JsonFailWithStack(ctx, err, nil)
		return
	}
	JsonSuccess(ctx, "success", nil)
}

func TestJsonFailWithStack(t *testing.T) {
	engine := gin.New()
	engine.GET("/regular_err", openNotExistFileHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/regular_err", nil)
	engine.ServeHTTP(w, req)
}
