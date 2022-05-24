package router

import (
	"regexp"
	logdata "server/packages/log"

	"github.com/gin-gonic/gin"
)

func ApiRouter() {
	// 範例

	//ip檢查
	//ginEngine.GinEngine.Use(authCheck)
	//建立商戶
	// ginEngine.GinEngine.POST("/create/merchant", CreateMerchant)
	// ginEngine.GinEngine.POST("/update/merchant", UpdateMerchantUri)

	// transaction := ginEngine.GinEngine.Group("/transaction")
	// transaction.Use(authCheck)

}

func authCheck(context *gin.Context) {

	var ok bool
	var err error

	ok, err = regexp.MatchString(`^10\.\d+\.\d+\.\d+$`, context.ClientIP())
	if err != nil {
		name := "authCheck regexp MatchString err"
		logdata.SysErrorLog(map[string]interface{}{
			"name": name,
			"ip":   context.ClientIP(),
		}, err)
		context.Abort()
		return
	}

	if ok {
		return
	}
}
