package migrate

import (
	"fmt"
	"os"

	"github.com/CRGao/log"
	command "github.com/heartz2o2o/db-migrate/command"
)

type ServerMigrate struct {
}

var Server ServerMigrate

func (m *ServerMigrate) Run() {
	dir, err := os.Getwd()
	if err != nil {
		log.Error("無法獲得程式目錄 err->", err)
		return
	}
	sqlDir := dir + "/sql"
	_, err = os.Stat(sqlDir)
	if err != nil {
		log.Error("找不到SQL資料夾 err->", err)
		return
	}
	link := os.Getenv("dbAccount") + ":" + os.Getenv("dbPassword") + "@tcp(" + os.Getenv("dbHost") + ":" + os.Getenv("dbPort") + ")/" + os.Getenv("dbName") + "?charset=" + os.Getenv("charset") + "&parseTime=true"
	env := &command.Environment{
		Dialect:    "mysql",
		DataSource: link,
		Dir:        sqlDir}
	command.SetEnvironment(env)
	migrate := command.UpCommand{}
	migrate.Run([]string{})
	fmt.Print("\n") //幫他換行
}
