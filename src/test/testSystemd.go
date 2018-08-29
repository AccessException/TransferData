package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	router.POST("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.GET("/test/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"err_code": 1,
			"err_msg":  "You have reached maximum request limit.",
			"data": gin.H{},
		})
	})

	router.Run(":9009")
}