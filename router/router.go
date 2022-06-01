package router

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"server/api"
	"server/router/middleware"
	"strings"
	"time"

	"github.com/CRGao/log"
	"github.com/gin-gonic/gin"
)

func Start() {
	// 範例
	router := gin.New()
	router.Use(loggerHandler)
	//綁定頁面
	router.LoadHTMLGlob("./web/view/*")
	router.Static("/assetPath", "./web/asset")

	//Router註冊
	InitPageRoute(router)
	InitApiRouter(router)
	err := router.Run(os.Getenv("socketIP") + ":" + os.Getenv("port"))
	if err != nil {
		panic(err)
	}
}

func InitApiRouter(router *gin.Engine) {
	//	全域使用
	groupApi := router.Group("/api")
	groupApi.Use(middleware.Cors())
	// groupApi.Use(middleware.SetLang())
	//router.GET("/account/setlang", v1.SetLang)

	groupApi.GET("/", api.Crud)

}

func InitPageRoute(router *gin.Engine) {
	//	全域使用
	//router.Use(middleware.Cors())
	router.GET("/", api.Index)

}

func loggerHandler(ctx *gin.Context) {
	// Start timer
	start := time.Now()
	path := ctx.Request.URL.Path
	raw := ctx.Request.URL.RawQuery
	method := ctx.Request.Method
	reqBody, _ := ctx.GetRawData()
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) //把請求資料再塞回去

	//把換行符號過濾
	reqBodyStr := strings.Replace(string(reqBody), "\n", "", -1)

	// Process request
	ctx.Next()

	// Stop timer
	end := time.Now()
	latency := end.Sub(start)
	statusCode := ctx.Writer.Status()
	clientIP := ctx.ClientIP()
	if raw != "" {
		path = path + "?" + raw
	}
	log.Info(fmt.Sprintf("METHOD:%s | PATH:%s | BODY:%s | CODE:%d | IP:%s | TIME:%s", method, path, reqBodyStr, statusCode, clientIP, latency))
}
