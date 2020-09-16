package model

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"goTodo/initialization"
	"goTodo/mylog"
	"net/http"
	"strings"
	"time"
)

type Flash map[string]string

func NewFlash(tag, message string) *Flash {
	return &Flash{
		"tag": tag,
		"message": message,
	}
}

type Message struct {
	Id int
	Message string
	NotifyTime string
	HasNotified bool
}

type TotalList struct {
	//Finished []string
	MessageList []Message
	PartFinishedInfo TotalFinishedInfo
}

func (t *TotalList) Search(username string) {
	currentTime := time.Now().Unix()
	rows, err := initialization.Db.Query("SELECT id, message, notify_time FROM " + initialization.DbMessageName + " WHERE username = ? AND finished_time IS NULL ORDER BY create_time DESC", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		mylog.GoTodoLogger.Panicln("查询数据库待办事宜出错：", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var message string
		var notifyTime sql.NullInt64
		err = rows.Scan(&id, &message, &notifyTime)
		if err == sql.ErrNoRows {
			return
		}
		// 只有在有通知时间，并且通知时间大于当前时间，即还没有进行通知的时候，才会把将要通知的时间显现出来，否则显示已通知
		if notifyTime.Valid {
			if notifyTime.Int64 > currentTime {
				t.MessageList = append(t.MessageList, Message{Id: id, Message: message, NotifyTime: time.Unix(notifyTime.Int64, 0).Format("2006/01/02 15:04:05"), HasNotified: false})
			} else {
				t.MessageList = append(t.MessageList, Message{Id: id, Message: message, NotifyTime: "-1", HasNotified: true})
			}

		} else {
			t.MessageList = append(t.MessageList, Message{Id: id, Message: message, NotifyTime: "", HasNotified: false})
		}

	}

	t.PartFinishedInfo.Query(username, 0, 5)
}

func AddMessage(username, message string) int64 {
	res, err := initialization.Db.Exec("INSERT INTO " + initialization.DbMessageName + "(username, message, create_time) values(?, ?, ?)", username, message, time.Now().Unix())
	if err != nil {
		mylog.GoTodoLogger.Panicln("新增待办事宜出错:", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		mylog.GoTodoLogger.Panicln("新增待办事宜出错:", err)
	}
	return id
}

func FinishMessage(username string, idList []string) int64 {
	inSql := strings.Join(idList, ",")
	curSql := "UPDATE " + initialization.DbMessageName + " SET finished_time = ? WHERE username = ? AND id IN ("
	curSql += inSql + ")"
	stmt, err := initialization.Db.Prepare(curSql)
	if err != nil {
		mylog.GoTodoLogger.Panicln("准备标记完成事宜出错:", err)
	}
	res, err := stmt.Exec(time.Now().Unix(), username)
	if err != nil {
		mylog.GoTodoLogger.Panicln("标记完成事宜出错:", err)
	}
	affect, err := res.RowsAffected()
	if err != nil {
		mylog.GoTodoLogger.Panicln("标记完成事宜出错:", err)
	}
	return affect
}

func UpdateMessageNotifyTime(username string, idList []string, settingTime time.Time) int64 {
	inSql := strings.Join(idList, ",")
	curSql := "UPDATE " + initialization.DbMessageName + " SET notify_time = ? WHERE username = ? AND id IN ("
	curSql += inSql + ")"
	stmt, err := initialization.Db.Prepare(curSql)
	if err != nil {
		mylog.GoTodoLogger.Printf("准备更新推送时间出错: %s", err)
		return -1
	}
	res, err := stmt.Exec(settingTime.Unix(), username)
	if err != nil {
		mylog.GoTodoLogger.Printf("更新推送时间出错: %s", err)
		return -1
	}
	affect, err := res.RowsAffected()
	if err != nil {
		mylog.GoTodoLogger.Printf("获取更新推送时间影响条数出错: %s", err)
		return -1
	}
	return affect
}

type Contents struct {
	Msgtype string `json:"msgtype"`
	Text ContentText `json:"text"`
	At ContentAt `json:"at"`
}

type ContentText struct {
	Content string `json:"content"`
}

type ContentAt struct {
	IsAtAll bool `json:"isAtAll"`
}

func NoticeMessage(username string, idList []string) error {
	var webhook string
	err := initialization.Db.QueryRow("SELECT webhook FROM " + initialization.DbUserName + " WHERE username = ?", username).Scan(&webhook)
	if err != nil {
		if err == sql.ErrNoRows {
			mylog.GoTodoLogger.Println("查询webhook时，要查询的用户:" + username + "不存在")
			return errors.New("查询webhook时，要查询的用户:" + username + "不存在")
		}
		mylog.GoTodoLogger.Println("查询数据库webhook出错:", err)
		return err
	}
	if webhook == "" {
		mylog.GoTodoLogger.Println(username + "的钉钉机器人还未设置")
		return errors.New("myError:钉钉机器人还未设置，请先设置")
	}

	inSql := strings.Join(idList, ",")
	// 保证id是对应的username的，以免不匹配导致使用了别人的信息
	curSql := "SELECT message FROM " + initialization.DbMessageName + " WHERE username = ? AND id IN (" + inSql + ")"
	rows, err := initialization.Db.Query(curSql, username)
	defer rows.Close()
	if err != nil {
		if err == sql.ErrNoRows {
			mylog.GoTodoLogger.Println(username + "指定推送的信息不存在")
			return errors.New("myError:指定推送的信息不存在")
		}
		mylog.GoTodoLogger.Println("查询数据库webhook出错:", err)
		return err
	}
	var messageList []string
	for rows.Next() {
		var message string
		err = rows.Scan(&message)
		if err == sql.ErrNoRows {
			continue
		}
		messageList = append(messageList, message)
	}
	c := &Contents{Msgtype: "text", Text: ContentText{Content: "[待办提醒]" + "\n" + strings.Join(messageList, "\n")}, At: ContentAt{IsAtAll: true}}
	data, err := json.Marshal(c)
	if err != nil {
		mylog.GoTodoLogger.Println(username + "需要推送的信息转化为json出错", err)
		return err
	}
	req, err := http.NewRequest("POST", webhook, bytes.NewReader(data))
	if err != nil {
		mylog.GoTodoLogger.Println(username + "需要推送的信息构造request请求出错", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		mylog.GoTodoLogger.Println(username + "需要推送的信息推送出错", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

const SplitSym = ":"

func RedisSubscribe() {
	client := initialization.Client
	sub := client.Subscribe("__keyevent@0__:expired")
	for msg:= range sub.Channel() {
		payload := msg.Payload
		nameId := strings.Split(payload, SplitSym)
		name, id := nameId[0], nameId[1]
		err := NoticeMessage(name, []string{id})
		if err != nil {
			mylog.GoTodoLogger.Printf("%s 的 %s 任务推送失败：%s", name, id, err)
		}
	}
}

func RedisPrepareNotify(username string, idList []string, settingTime time.Time) {
	// redis存储
	client := initialization.Client
	for _, v := range idList {
		key := username + SplitSym + v
		client.Set(key, "0", 0)
		client.ExpireAt(key, settingTime)
	}
	// 数据库存储，存在数据库是为了能在页面上显示，方便更改
	_ = UpdateMessageNotifyTime(username, idList, settingTime)
}


func CancelNotify(username, id string) error {
	// 先取消redis里的过期时间，然后将数据库内的 notify_time 修改
	client := initialization.Client
	key := username + SplitSym + id
	client.Del(key)

	curSql := "UPDATE " + initialization.DbMessageName + " SET notify_time = NULL WHERE username = ? AND id = ?"
	stmt, err := initialization.Db.Prepare(curSql)
	if err != nil {
		mylog.GoTodoLogger.Printf("准备取消推送出错: %s", err)
		return errors.New("myError:已正常取消推送，服务器内部错误导致无法取消显示消息上的详情，请勿担心")
	}
	_, err = stmt.Exec(username, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		mylog.GoTodoLogger.Printf("取消推送出错: %s", err)
		return errors.New("myError:已正常取消推送，服务器内部错误导致无法取消显示消息上的详情，请勿担心")
	}
	return nil
}

type FinishedInfo struct {
	Message string
	CreateTime string
	FinishedTime string
}

type TotalFinishedInfo []FinishedInfo

func (t *TotalFinishedInfo) Query(username string, start, limit int) {
	rows, err := initialization.Db.Query("SELECT message, create_time, finished_time FROM " + initialization.DbMessageName + " WHERE username = ? AND finished_time IS NOT NULL ORDER BY finished_time DESC LIMIT ?, ?", username, start, limit)
	defer rows.Close()
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		mylog.GoTodoLogger.Panicln("查询数据库全部完成事宜出错：", err)
	}
	for rows.Next() {
		var message string
		var createTime int64
		var finishedTime int64
		err = rows.Scan(&message, &createTime, &finishedTime)
		if err == sql.ErrNoRows {
			return
		}
		*t = append(*t, FinishedInfo{Message: message, CreateTime: time.Unix(createTime, 0).Format("2006-01-02 15:04:05"), FinishedTime: time.Unix(finishedTime, 0).Format("2006-01-02 15:04:05")})
	}
}

func ClearAll(username string) int64 {
	rows, err := initialization.Db.Exec("DELETE FROM " + initialization.DbMessageName + " WHERE username = ? AND finished_time IS NOT NULL", username)
	if err != nil {
		mylog.GoTodoLogger.Panicln("删除完成事宜出错:", err)
	}
	affect, err := rows.RowsAffected()
	if err != nil {
		mylog.GoTodoLogger.Panicln("删除完成事宜出错:", err)
	}
	mylog.GoTodoLogger.Printf("已删除%s的全部完成事宜", username)
	return affect
}