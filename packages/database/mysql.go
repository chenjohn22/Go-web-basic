package db

import (
	"os"
	"strconv"
	"time"

	"github.com/CRGao/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func NewMysqlStore() *gorm.DB {
	var db *gorm.DB
	max, _ := strconv.Atoi(os.Getenv("dbMaxConnection"))
	min, _ := strconv.Atoi(os.Getenv("dbMinConnection"))
	link := os.Getenv("dbAccount") + ":" + os.Getenv("dbPassword") + "@tcp(" + os.Getenv("dbHost") + ":" + os.Getenv("dbPort") + ")/" + os.Getenv("dbName") + "?charset=" + os.Getenv("charset") + "&parseTime=true"
	log.Debug("connect to mysql : ", link)
	db, _ = gorm.Open("mysql", link)

	//连接池设置
	db.DB().SetMaxOpenConns(max) //用于设置最大打开的连接数，默认值为0表示不限制
	db.DB().SetMaxIdleConns(min) //用于设置闲置的连接数。
	db.DB().SetConnMaxLifetime(time.Minute * 5)
	err := db.DB().Ping()
	if err != nil {
		log.Error("MySQL Connection Fail：" + err.Error())
		panic(err)
	}
	return db
}
