package initialization

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var Db *sql.DB
const DbUserName = "user"
const DbMessageName = "message"

func init() {
	var err error
	Db, err = sql.Open("sqlite3", "sqlite3.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	// 用户表
	stmt, err := Db.Prepare("CREATE TABLE IF NOT EXISTS `"+ DbUserName +"` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `username` VARCHAR(20) NOT NULL, `password_hash` VARCHAR(128) NOT NULL)")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err.Error())
	}
	// 代办事情表
	stmt, err = Db.Prepare("CREATE TABLE IF NOT EXISTS `"+ DbMessageName +"` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `username` VARCHAR(20) NOT NULL, `message` VARCHAR(128) NOT NULL, `create_time` INTEGER NOT NULL, `finished_time` INTEGER)")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err.Error())
	}
	// 最大连接数
	Db.SetMaxOpenConns(20)
	// 最大空闲连接数
	Db.SetMaxIdleConns(20)
}
