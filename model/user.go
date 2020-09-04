package model

import (
	"database/sql"
	"goTodo/initialization"
	"goTodo/mylog"
	"goTodo/util"
)

type UserModel struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type RegisterModel struct {
	UserModel
	PasswordAgain string `form:"passwordAgain"`
}

func (user *UserModel) ValidUser() bool {
	row := initialization.Db.QueryRow("SELECT password_hash FROM " + initialization.DbUserName + " WHERE username = ?", user.Username)
	var passwordHash string
	if err := row.Scan(&passwordHash); err != nil {
		if err == sql.ErrNoRows {
			// 说明没有该用户
			return false
		}
		mylog.GoTodoLogger.Panicln("验证用户信息发生错误:", err)
	}
	return passwordHash == util.GenStringHash(user.Password)
}

func (user *UserModel) SaveUser() int64 {
	passwordHash := util.GenStringHash(user.Password)
	res, err := initialization.Db.Exec("INSERT INTO " + initialization.DbUserName + "(username, password_hash) values(?, ?)", user.Username, passwordHash)
	if err != nil {
		mylog.GoTodoLogger.Panicln("注册新用户发生错误:", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		mylog.GoTodoLogger.Panicln("注册新用户发生错误:", err)
	}
	return id
}

func (user *UserModel) ExistUser() bool {
	row := initialization.Db.QueryRow("SELECT id FROM " + initialization.DbUserName + " WHERE username = ?", user.Username)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		mylog.GoTodoLogger.Panicln("查询用户信息发生错误:", err)
	}
	return true
}