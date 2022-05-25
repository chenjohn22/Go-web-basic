package middleware

import (
	"github.com/CRGao/log"
	"github.com/gin-gonic/gin"
)

func SetLang() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var lang string
		baseLang := "zh-CN"
		defer func() {
			if err := recover(); err != nil {
				log.Error("err->", err)
			}
		}()
		lang, _ = ctx.Cookie("lang")
		if lang == "" {
			ctx.SetCookie("lang", baseLang, 0, "/", "", false, false)
			ctx.Set("lang", baseLang)
		}
	}
}
