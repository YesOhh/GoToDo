package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goTodo/initialization"
	"goTodo/model"
	"goTodo/mylog"
	"net/http"
	"strconv"
	"time"
)

func index(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	var totalList = new(model.TotalList)
	totalList.Search(username.(string))
	userSetting := model.ShowSettings(username.(string))
	canNotify := initialization.Configuration.RedisSetting.Exists && (string(userSetting.CurrentWebHook) != "")
	c.HTML(http.StatusOK, "core.tmpl", gin.H{
		"title":        "Go To Do",
		"messageList":      totalList.MessageList,
		"partFinishedInfo": totalList.PartFinishedInfo,
		"username": username,
		"canNotify": canNotify,
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

// 此处判断按钮是submit还是notify来决定做什么
func Operate(c *gin.Context) {
	session := sessions.Default(c)
	idList :=  c.PostFormArray("todo")
	// 如果没有选中任何消息，而且不是取消通知，则直接跳转
	if len(idList) == 0 && c.PostForm("delNotify") == "" {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	username := session.Get("username")

	if c.PostForm("submit") != "" {
		// 点击提交完成
		affect := model.FinishMessage(username.(string), idList)
		mylog.GoTodoLogger.Printf("%d 条 %s 的待办事宜已完成", affect, username)
	} else if c.PostForm("notify") != "" {
		// 点击通知
		userSetting := model.ShowSettings(username.(string))
		canNotify := initialization.Configuration.RedisSetting.Exists && (string(userSetting.CurrentWebHook) != "")
		// 如果没开启redis或者没有配置webhook
		if !canNotify {
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		// 如果没有设置日期，直接推送
		// 如果没有时间，默认为 00:00:00
		stringTime := ""
		var settingTime time.Time
		if c.PostForm("notifyDate") == "" {
			settingTime = time.Now()
		} else {
			stringTime = c.PostForm("notifyDate") + " "
			if c.PostForm("notifyTime") == "" {
				stringTime += "00:00:00"
			} else {
				stringTime += c.PostForm("notifyTime")
			}
			loc, _ := time.LoadLocation("Local")
			var err error
			settingTime, err = time.ParseInLocation("2006/01/02 15:04:05", stringTime, loc)
			if err != nil {
				mylog.GoTodoLogger.Panicln("时间戳转换出错", err)
			}
		}

		if settingTime.Unix() <= time.Now().Unix() {
			// 已经过了过期时间，直接推送
			_ = model.UpdateMessageNotifyTime(username.(string), idList, settingTime)
			err := model.NoticeMessage(username.(string), idList)
			if err != nil {
				mylog.GoTodoLogger.Panicln(err)
			}
			mylog.GoTodoLogger.Printf("%s 的消息推送成功", username)
		} else {
			// 存到redis中，通过订阅监听来准时推送
			model.RedisPrepareNotify(username.(string), idList, settingTime)
		}
	} else if c.PostForm("delNotify") != "" {
		// 点击取消通知
		err := model.CancelNotify(username.(string), c.PostForm("delNotify"))
		if err != nil {
			mylog.GoTodoLogger.Panicln(err)
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/")

}

func Clear(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	_ = model.ClearAll(username.(string))
	// get 方法如果是 301 永久跳转，就会不访问，所以换成302
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