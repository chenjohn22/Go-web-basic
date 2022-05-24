package ginEngine

import (
	"net/http"
	"os"

	logdata "server/packages/log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var sessionStore sessions.Store
var GinEngine *gin.Engine

func GinInit() {

	// 初始化引擎
	gin.SetMode(gin.ReleaseMode)
	GinEngine = gin.Default()
	// pprof.Register(GinEngine) // 性能

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	GinEngine.Use(cors.New(config))

	// GinEngine.Use(cors())
	GinEngine.Use(esLog)

	// healthcheck
	GinEngine.GET("/healthcheck", healthCheck)

	return
}

func healthCheck(c *gin.Context) {
	c.String(http.StatusOK, "healthCheck ver : "+os.Getenv("version"))
}

func esLog(context *gin.Context) {
	// log.Printf("esLog : %+v", context.Request.URL)
	if context.HandlerName() != "server/external/ginEngine.healthCheck" {
		context.PostForm("")

		logdata.SysLog(map[string]interface{}{
			"name":   context.HandlerName(),
			"header": context.Request.Header,
			"form":   context.Request.Form,
		})
	}

	// Pass on to the next-in-chain
	context.Next()
}
