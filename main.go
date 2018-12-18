package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/v1/log", SimpleLog)

	router.Run()
}

var token = "yinzhinyy"

func SimpleLog(c *gin.Context) {
	message := ""
	if len(c.Request.URL.Query()) == 0 {
		message = "hello, this is handle view"
	}
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	list := []string{token, timestamp, nonce}
	sort.Strings(list)
	sha1 := sha1.New()
	var srcStr string
	for _, str := range list {
		srcStr = srcStr + str
	}
	sha1.Write([]byte(srcStr))
	hashcode := fmt.Sprintf("%x", sha1.Sum(nil))
	log.Println("handle/GET func: hashcode, signature: %s, %s", hashcode, signature)
	if hashcode == signature {
		log.Println("passed verification")
		message = echostr
	}
	c.String(http.StatusOK, message)
}
