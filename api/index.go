package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Crud(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
	return
}

func Index(context *gin.Context) {
	context.HTML(http.StatusOK, fmt.Sprintf("index.tmpl"), gin.H{
		"test": "1234",
	})
}
