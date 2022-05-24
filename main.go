package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	db "server/packages/database"
	"server/packages/ginEngine"
	"server/packages/redis"
	"server/router"
	"time"
)

func main() {

	timeUnix := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("main timeUnix : [%d] ", timeUnix)

	loadEnv()
	ginEngine.GinInit()
	db.NewConnection()
	redis.NewRedisConn()

	router.ApiRouter()
	ginEngine.GinEngine.Run(os.Getenv("port"))

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
	log.Printf("goServerApiPort : %+v", os.Getenv("port"))
	log.Printf("version : %+v", os.Getenv("version"))
}
