package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goTodo/mylog"
	"net/http"
	"strings"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/login") || strings.HasPrefix(c.Request.URL.Path, "/register") {
			c.Next()
		} else if strings.HasPrefix(c.Request.URL.Path, "/static") {
			c.Next()
		} else {
			session := sessions.Default(c)
			authenticated := session.Get("todo")
			if authenticated != true {
				mylog.GoTodoLogger.Println("存在过期的session，需要重新登录")
				c.Redirect(http.StatusMovedPermanently, "/login")
				// 不加c.Abort 发现会导致继续进入路由执行，用导致路由中用的session没有需要的信息，然后panic，然后再回到此处的redirect
				c.Abort()
				return
			}
			c.Next()
		}
	}
}