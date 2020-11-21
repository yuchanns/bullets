package common

import (
	"github.com/gin-gonic/gin"
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
	Logger.Fields(map[string]interface{}{"stack": stack}).Error(ctx, err)
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
