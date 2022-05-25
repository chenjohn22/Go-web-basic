package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Crud(context *gin.Context) {

}

func Index(context *gin.Context) {
	context.HTML(http.StatusOK, fmt.Sprintf("index.tmpl"), gin.H{
		"test": "1234",
	})
}
