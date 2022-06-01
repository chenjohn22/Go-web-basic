package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	logs "log"
	"os"
	"os/signal"
	"server/packages/ginEngine"
	"server/router"
	"syscall"
	"time"

	"github.com/CRGao/log"
)

func main() {

	//timeUnix := time.Now().UnixNano() / int64(time.Millisecond)

	loadEnv()
	//log 初始化
	conf := log.LogConfig{
		Level:      os.Getenv("version"),
		FileName:   "./logs/log.log",
		HasConsole: true,
		Color:      true,
		Json:       true,
		MaxSize:    20,
		MaxAge:     20,
		DateSlice:  "m",
		Format:     "%{time:2006/01/02 15:04:05.000} [%{level:.4s}] %{shortfile} %{shortfunc} %{message}",
	}
	log.InitByConfigStruct(&conf)
	ginEngine.GinInit()
	//db.NewConnection()
	//redis.NewRedisConn()
	//migrate
	//go migrate.Server.Run()

	// router.ApiRouter()
	// ginEngine.GinEngine.Run(os.Getenv("port"))
	go router.Start()

	//crontab
	// go crontab.Start()

	//優雅的結束
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	select {
	case s := <-signalChan:
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			time.Sleep(10 * time.Second)
			fmt.Println("退出")
		case syscall.SIGUSR1:
			fmt.Println("usr1")
		case syscall.SIGUSR2:
			fmt.Println("usr2")
		default:
			fmt.Println("other")
		}
		fmt.Println("程式獲得訊號：", s)
	}
}

func loadEnv() {

	config, err := ioutil.ReadFile("config/env.json")
	if err != nil {
		log.Fatal("找不到env.json")
	}

	configHost := make(map[string]string)
	// log.Printf("configHost : %v\n", configHost)

	err = json.Unmarshal(config, &configHost)
	if err != nil {
		// log.Printf("configHost err: %v\n", err)
		return
	}

	// log.Printf("configHost : %v\n", configHost)
	for k, v := range configHost {
		// log.Printf("%s : %s\n", k, v)
		_ = os.Setenv(k, v)
	}

	//os.Setenv("url", "http://localhost")
	logs.Printf("goServerApiPort : %+v", os.Getenv("port"))
	logs.Printf("version : %+v", os.Getenv("version"))
}
