package redis

import (
	"os"
	"strconv"

	"github.com/CRGao/log"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var redisConn *redis.Client

func NewRedisConn() {
	var ctx *gin.Context
	redisConn = NewRedisStore(ctx)
}

func GetRedisConn(context *gin.Context) *redis.Client {
	if redisConn == nil {
		log.Error("無法獲得數據庫連線！")
		return nil
	}

	times := 0
	conn := redisConn
	for {
		pong, _ := conn.Ping(context).Result()
		if pong == "PONG" || times > 5 {
			break
		}
		//連線未正確回覆 重新獲取一條
		conn = redisConn
		times++
	}

	return conn
}

func NewRedisStore(ctx *gin.Context) *redis.Client {
	poolSize, _ := strconv.Atoi(os.Getenv("poolSize"))
	selectDB, _ := strconv.Atoi(os.Getenv("selectDB"))
	port := os.Getenv("redisPort")
	if port == "" {
		port = ":6379"
	}
	log.Info("connect to redis : ", os.Getenv("redisHost")+":"+port, " db:", os.Getenv("selectDB"))

	Conn := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("redisHost") + ":" + port,
		Password:     os.Getenv("redisPassword"),
		DB:           selectDB,
		PoolSize:     poolSize,     //(10 * CPU number)
		MinIdleConns: poolSize / 2, //最小空閒連接數
		//DialTimeout:  10 * time.Second, //Dial timeout for establishing new connections.(5s)
		//ReadTimeout:  30 * time.Second, //Timeout for socket reads.(3s)
		//WriteTimeout: 30 * time.Second, //Timeout for socket writes. (同ReadTimeout)
		//PoolTimeout:  30 * time.Second,  //(ReadTimeout + 1s)
		//PoolTimeout:        30 * time.Second, //連線持獲取超時時間 (ReadTimeout + 1s)
		//IdleTimeout:        30 * time.Second, //空閒連接的超時時間 (5minute)
		//IdleCheckFrequency: 30 * time.Second, //超時空閒連接的間隔時間 (1minute)
	})
	_, err := Conn.Ping(ctx).Result()
	if err != nil {
		log.Error("連接redis失敗：" + err.Error())
		panic(err)
	}
	return Conn
}
