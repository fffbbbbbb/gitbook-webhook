package main

import (
	"fmt"
	"gitbookWebhook/controller"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	err := viper.BindEnv("webhook_secret", "webhook_secret")
	if err != nil {
		err := fmt.Errorf("err:cannot get webhook_secret : %s", err.Error())
		panic(err)
	}
	secret := viper.GetString("webhook_secret")
	if secret == "" {
		err = fmt.Errorf("secret is nil")
		panic(err)
	}
}
func main() {
	r := gin.Default()
	r.Use(controller.Cors(), gin.Recovery())
	r.POST("/webhook", controller.Webhook)
	//初始化http配置
	s := &http.Server{
		Addr:           "0.0.0.0:9001",
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println(`Welcome to labortool-go 
	 默认运行地址: http://0.0.0.0:9001`)

	err := s.ListenAndServe()
	if err != nil {
		panic("bootstrap error")
	}
}
