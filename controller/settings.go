package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goTodo/model"
	"net/http"
)

func ShowSettings(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	userSetting := model.ShowSettings(username.(string))
	currentWebHook := "暂未设置"
	if string(userSetting.CurrentWebHook) != "" {
		currentWebHook = string(userSetting.CurrentWebHook)
	}
	c.HTML(http.StatusOK, "setting.tmpl", gin.H{
		"title": "设置",
		"currentWebHook": currentWebHook,
		"username": username.(string),
	})
}

func UpdateWebhook(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	webhook := c.PostForm("webhook")
	_ = model.UpdateWebHook(username.(string), webhook)
	// 此处需要返回当前页面
	c.Redirect(http.StatusMovedPermanently, "/setting")
}