package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"sync"
)

var db *sql.DB
var once sync.Once

func init() {
	once.Do(initMysql)
}

func initMysql() {
	_db, err := sql.Open("mysql", "dev:123456@/friend?charset=utf8")
	db = _db
	if err != nil {
		log.Println("database initialize error : ", err.Error())
		os.Exit(-1)
	}
}
