package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goTodo/model"
	"goTodo/mylog"
	"net/http"
	"strconv"
)

func index(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	var totalList = new(model.TotalList)
	totalList.Search(username.(string))
	c.HTML(http.StatusOK, "core.tmpl", gin.H{
		"title":        "Go To Do",
		"messageList":      totalList.MessageList,
		"partFinishedInfo": totalList.PartFinishedInfo,
		"username": username,
	})
}

func Add(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	message := c.PostForm("message")
	id := model.AddMessage(username.(string), message)
	mylog.GoTodoLogger.Printf("%s 新增 id 为 %d 的待办事宜", username, id)
	c.Redirect(http.StatusMovedPermanently, "/")
}

func Finish(c *gin.Context) {
	session := sessions.Default(c)
	idList :=  c.PostFormArray("todo")
	if len(idList) != 0 {
		username := session.Get("username")
		affect := model.FinishMessage(username.(string), idList)
		mylog.GoTodoLogger.Printf("%d 条 %s 的待办事宜已完成", affect, username)
	}
	c.Redirect(http.StatusMovedPermanently, "/")
}

func Clear(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	_ = model.ClearAll(username.(string))
	// 因为此处是301 导致不会执行了！！！！！！！！！！！
	c.Redirect(http.StatusFound, "/")
}

func Finished(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	strStart := c.PostForm("start")
	strLastId := c.PostForm("lastId")
	direction := c.PostForm("submit")
	var start int
	if strStart != "" {
		preStart, err := strconv.Atoi(strStart)
		if err != nil {
			mylog.GoTodoLogger.Panicln("完成事项查询起始start参数无法转成数字", err)
		}
		if direction != "上一页" {
			start = preStart
		} else {
			start = preStart - 10
		}
	}
	if strLastId != "" {
		preLastId, err := strconv.Atoi(strLastId)
		if err != nil {
			mylog.GoTodoLogger.Panicln("完成事项查询展示lastId参数无法转成数字", err)
		}
		if direction != "上一页" {
			start += preLastId + 1
		}
	}

	var hasPre bool
	if start != 0 {
		hasPre = true
	}
	totalFinishedInfo := new(model.TotalFinishedInfo)
	totalFinishedInfo.Query(username.(string), start, 10)
	var hasNext bool
	if len(*totalFinishedInfo) == 10 {
		hasNext = true
	}
	c.HTML(http.StatusOK, "finished.tmpl", gin.H{
		"username": username,
		"totalFinishedInfo": totalFinishedInfo,
		"start": start,
		"hasNext": hasNext,
		"hasPre": hasPre,
	})
}