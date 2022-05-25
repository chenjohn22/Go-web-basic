package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code int
		code = http.StatusBadRequest
		token := ctx.GetHeader("X-COMM-SECRET")

		if token == "7Lim-8nTkv4n9Mr4Q" {
			code = http.StatusOK
		} else {
			code = http.StatusUnauthorized
		}

		if code != http.StatusOK {
			ctx.JSON(code, gin.H{
				"msg":    "Token error",
				"result": code,
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
