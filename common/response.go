package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// JsonFail serializes the given struct as JSON into the response body.
func JsonFail(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, ToFailMsg(msg, data))
}

// ToFailMsg returns a preset fail content with gin.H
func ToFailMsg(msg string, data interface{}) gin.H {
	return gin.H{
		"code":    500,
		"message": msg,
		"data":    data,
	}
}
