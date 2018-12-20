package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/yinzhinyy/talk-to-ilona/internal"
)

var db = internal.LoadDB("mongo")

func main() {
	router := gin.Default()

	router.GET("/v1/log", validateHandler)
	router.POST("/v1/log", simpleLog)

	router.Run()
}

type ReqBody struct {
	ToUserName   string `xml:"ToUserName"  binding:"required"`
	FromUserName string `xml:"FromUserName"  binding:"required"`
	CreateTime   int    `xml:"CreateTime"  binding:"required"`
	MsgType      string `xml:"MsgType"  binding:"required"`
	Content      string `xml:"Content"  binding:"required"`
	MsgId        int64  `xml:"MsgId"  binding:"required"`
}

func simpleLog(c *gin.Context) {
	var message = ""
	var reqBody ReqBody
	if err := c.ShouldBindXML(&reqBody); err != nil {
		log.Fatalln("illegal xml format")
		c.String(http.StatusBadRequest, message)
		return
	}
	switch reqBody.MsgType {
	case "text":
		handle(reqBody)
	default:
		message = "success"
	}
	c.String(http.StatusOK, message)
}

func handle(reqBody ReqBody) {
	doc := bson.M{
		"userName":   reqBody.FromUserName,
		"createTime": reqBody.CreateTime,
		"activity":   reqBody.Content,
	}
	db.Save("talk", "activity_log", doc)
}

const token = "yinzhinyy"

func validateHandler(c *gin.Context) {
	var message string
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
	log.Printf("handle/GET func: hashcode, signature: %s, %s", hashcode, signature)
	if hashcode == signature {
		log.Println("passed verification")
		message = echostr
	}
	c.String(http.StatusOK, message)
}
