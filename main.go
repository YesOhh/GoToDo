package main

import (
	"github.com/gin-gonic/gin"
	"goTodo/controller"
	"goTodo/initialization"
)

func main() {
	r := gin.Default()
	controller.LoadRouters(r)
	ip := initialization.Configuration.Setting.Ip
	port := initialization.Configuration.Setting.Port
	if ip == "" {
		ip = "127.0.0.1"
	}
	if port == "" {
		port = "8080"
	}
	r.Run(ip + ":" + port)
}