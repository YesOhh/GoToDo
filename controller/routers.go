package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"goTodo/middleware"
	"net/http"
)

const MaxAge = 60*60*24

func LoadRouters(r *gin.Engine)  {
	// 第二个是存放静态资源的地址，绑到第一个参数上，用于模板中调用
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	store := cookie.NewStore([]byte("ToDoSecret"))
	store.Options(sessions.Options{MaxAge: MaxAge})
	// 因为 Authorization 用了session的key，所以要写在后面
	r.Use(sessions.Sessions("mysession", store))
	r.Use(middleware.Recovery())
	r.Use(middleware.Authorization())
	r.GET("/login", Login)
	r.POST("/login", Login)
	r.GET("/register", Register)
	r.POST("/register", Register)
	r.GET("/logout", Logout)
	r.GET("/", index)
	r.GET("/finished", Finished)
	r.POST("/finished", Finished)
	v := r.Group("/operation")
	v.POST("/add", Add)
	v.POST("/finish", Finish)
	v.GET("/clear", Clear)
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"title": "发生错误",
			"error": "未找到该页面",
		})
	})
}
