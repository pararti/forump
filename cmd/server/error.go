package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(code int, err string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(code, gin.H{"error": err})
		ctx.HTML(http.StatusOK, "error", gin.H{
			"Error": err,
		})
	}
}

func SuccessHandler(suc string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": suc})
		ctx.HTML(http.StatusOK, "success", gin.H{
			"Success": suc,
		})
	}
}
