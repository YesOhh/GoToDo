package model

import (
	"database/sql"
	"goTodo/initialization"
	"goTodo/mylog"
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
}

type TotalList struct {
	//Finished []string
	MessageList []Message
	PartFinishedInfo TotalFinishedInfo
}

func (t *TotalList) Search(username string) {
	rows, err := initialization.Db.Query("SELECT id, message FROM " + initialization.DbMessageName + " WHERE username = ? AND finished_time IS NULL ORDER BY create_time DESC", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		mylog.GoTodoLogger.Panicln("查询数据库待办事宜出错：", err)
	}
	for rows.Next() {
		var id int
		var message string
		err = rows.Scan(&id, &message)
		if err == sql.ErrNoRows {
			return
		}
		t.MessageList = append(t.MessageList, Message{Id: id, Message: message})
	}

	//rows, err = initialization.Db.Query("SELECT message FROM " + initialization.DbMessageName + " WHERE username = ? AND finished_time IS NOT NULL ORDER BY finished_time DESC LIMIT 7", username)
	//if err != nil {
	//	if err == sql.ErrNoRows {
	//		return
	//	}
	//	mylog.GoTodoLogger.Panicln("查询数据库完成事宜出错：", err)
	//}
	//for rows.Next() {
	//	var message string
	//	err = rows.Scan(&message)
	//	if err == sql.ErrNoRows {
	//		return
	//	}
	//	t.Finished = append(t.Finished, message)
	//}

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
		mylog.GoTodoLogger.Panicln("准备标价完成事宜出错:", err)
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

type FinishedInfo struct {
	Message string
	CreateTime string
	FinishedTime string
}

type TotalFinishedInfo []FinishedInfo

func (t *TotalFinishedInfo) Query(username string, start, limit int) {
	rows, err := initialization.Db.Query("SELECT message, create_time, finished_time FROM " + initialization.DbMessageName + " WHERE username = ? AND finished_time IS NOT NULL ORDER BY finished_time DESC LIMIT ?, ?", username, start, limit)
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