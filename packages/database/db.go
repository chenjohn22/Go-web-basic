package db

import (
	"os"

	"github.com/CRGao/log"
	"github.com/jinzhu/gorm"
)

var storage Storage

type Storage struct {
	dbMysql *gorm.DB
}

func NewConnection() {
	storage.dbMysql = NewMysqlStore()
}

func (sto *Storage) GetDbConn() *gorm.DB {
	if sto.dbMysql == nil {
		log.Error("無法獲得數據庫連線！")
		return nil
	}

	if os.Getenv("version") == "dev" {
		return sto.dbMysql.Debug()
	}

	return sto.dbMysql
}
