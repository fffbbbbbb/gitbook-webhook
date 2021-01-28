package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Webhook(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("err:", err.Error())
		FailResp(c, "read Body error"+err.Error())
		return
	}
	key := []byte(viper.GetString("webhook_secret"))
	mac := hmac.New(sha256.New, key)
	_, err = mac.Write(body)
	if err != nil {
		log.Println("err:", err.Error())
		FailResp(c, err.Error())
		return
	}
	if c.GetHeader("X-Hub-Signature-256") != "sha256="+fmt.Sprintf("%x", mac.Sum(nil)) {
		FailResp(c, "wrong Signature")
		return
	}
	log.Println("right Signature")
	SuccessResp(c, nil)
}

type Response struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Router  string      `json:"router"`
	Data    interface{} `json:"data"`
}

//FailResp 错误返回
func FailResp(c *gin.Context, Msg string) {
	resp := Response{
		Success: false,
		Router:  c.Request.URL.RequestURI(),
		Msg:     Msg,
	}
	c.JSON(500, &resp)
	c.Abort()
}

//SuccessResp 错误返回
func SuccessResp(c *gin.Context, data interface{}) {
	resp := Response{
		Success: true,
		Router:  c.Request.URL.RequestURI(),
		Data:    data,
	}
	c.JSON(200, &resp)
	c.Abort()
}

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
