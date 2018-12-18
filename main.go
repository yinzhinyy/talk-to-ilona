package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/v1/log", SimpleLog)

	router.Run()
}

func SimpleLog(c *gin.Context) {
	signature := c.Query("signature")
	log.Println(signature)
}
