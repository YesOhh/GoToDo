package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goTodo/model"
	"goTodo/mylog"
	"net/http"
	"regexp"
)

func Login(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "登录",
		})
		return
	}
	var user model.UserModel
	if err := c.ShouldBind(&user); err != nil {
		mylog.GoTodoLogger.Println("登录信息参数错误", err)
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "登录",
			"flash": model.NewFlash("warning", "登录信息参数错误"),
		})
		return
	}

	if user.ValidUser() {
		session := sessions.Default(c)
		session.Set("todo", true)
		session.Set("username", user.Username)
		err := session.Save()
		if err != nil {
			mylog.GoTodoLogger.Panicln("登录时session出错:", err)
		}
		c.Redirect(http.StatusMovedPermanently, "/")
	} else {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "登录",
			"flash": model.NewFlash("warning", "用户名或密码错误"),
		})
	}
}

func Register(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title": "注册",
		})
		return
	}
	var userRegister model.RegisterModel
	if err := c.ShouldBind(&userRegister); err != nil {
		mylog.GoTodoLogger.Println("登录信息参数错误", err)
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title": "注册",
			"flash": model.NewFlash("warning", "注册信息参数错误"),
		})
		return
	}
	res, err := regexp.MatchString(`^[a-zA-Z0-9_]+$`, userRegister.Username)
	if err != nil {
		mylog.GoTodoLogger.Panicln("正则匹配出错:", err)

	}
	if !res {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title": "注册",
			"flash": model.NewFlash("warning", "用户名仅支持字母、数字和下划线"),
		})
		return
	}

	if userRegister.Password != userRegister.PasswordAgain {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title": "注册",
			"flash": model.NewFlash("warning", "两次填写密码不同"),
		})
		return
	}

	if len(userRegister.Password) < 6 {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title": "注册",
			"flash": model.NewFlash("warning", "密码长度太短，请重新设置"),
		})
		return
	}

	if userRegister.UserModel.ExistUser() {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title": "注册",
			"flash": model.NewFlash("warning", "该用户名已存在"),
		})
		return
	}

	id := userRegister.UserModel.SaveUser()
	mylog.GoTodoLogger.Printf("注册成功，id为：%d", id)
	session := sessions.Default(c)
	session.Set("todo", true)
	session.Set("username", userRegister.Username)
	err = session.Save()
	if err != nil {
		mylog.GoTodoLogger.Panicln("登录时session出错:", err)
	}
	c.Redirect(http.StatusMovedPermanently, "/")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("todo") == true {
		username := session.Get("username")
		session.Clear()
		_ = session.Save()
		mylog.GoTodoLogger.Printf("%s log out, current authenticated: %s", username, session.Get("status"))
		c.Redirect(http.StatusFound, "/login")
	} else {
		c.Redirect(http.StatusFound, "/login")
	}
}