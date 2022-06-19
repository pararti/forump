package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(code int, err, url, urlhandler string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "error", gin.H{
			"Error":      err,
			"URL":        url,
			"URLHandler": urlhandler,
		})
	}
}

func SuccessHandler(suc string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "success", gin.H{
			"Success": suc,
		})
	}
}
