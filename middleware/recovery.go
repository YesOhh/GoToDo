package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goTodo/mylog"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				defaultMessage := "服务器内部发生错误"
				message := fmt.Sprintf("%s", err)
				mylog.GoTodoLogger.Printf("发生错误: %s\n", trace(message))
				// 如果错误是以 myError 开头，则显示后面部分的信息
				if strings.HasPrefix(message, "myError") {
					defaultMessage = strings.SplitN(message, ":", 1)[1]
				}
				c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
					"title": "发生错误",
					"error": defaultMessage,
				})
			}
		}()
		c.Next()
	}
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}