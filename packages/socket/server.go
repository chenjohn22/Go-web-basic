package socket

import (
	"net/http"
	"os"
	"server/packages/utils"
	"strconv"
	"sync"
	"time"

	"github.com/CRGao/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	defaultAotuCloseTime int64 = 60 //秒
)

const (
	ConnectClose ConnStatus = iota
	Connected
	//Connecting
	//ConnectError
)

var SafeConn sync.Map //concurrent map 存储socket连接

var SocketServer *Server

type ConnStatus int

type Server struct {
	TimeOut       int
	ReadBuffer    int
	WriteBuffer   int
	AutoCloseTime int64
	upgrader      websocket.Upgrader
	cmdList       map[int]func(msg *interface{}) error
	regFunc       map[string]*RegFunc
}

type RegFunc struct {
	Connect func(*Conns) error              //创建链接时执行的方法
	HasMsg  func(*ReceiveMsg, *Conns) error //当有消息时，执行的方法
	Close   func(*Conns) error              //当关闭链接时，执行的方法
}

func NewServer() *Server {
	SocketServer = new(Server)
	timeout, _ := strconv.Atoi(os.Getenv("socketTimeOut"))
	readBuffer, _ := strconv.Atoi(os.Getenv("socketReadBuffer"))
	writeBuffer, _ := strconv.Atoi(os.Getenv("socketWriteBuffer"))
	autoClose, _ := strconv.Atoi(os.Getenv("socketAutoClose"))
	SocketServer.TimeOut = timeout
	SocketServer.TimeOut = readBuffer
	SocketServer.TimeOut = writeBuffer
	SocketServer.AutoCloseTime = int64(autoClose)
	if SocketServer.AutoCloseTime == 0 {
		SocketServer.AutoCloseTime = defaultAotuCloseTime
	}
	SocketServer.upgrader = websocket.Upgrader{
		ReadBufferSize:   SocketServer.ReadBuffer,
		WriteBufferSize:  SocketServer.WriteBuffer,
		HandshakeTimeout: time.Duration(SocketServer.TimeOut),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	SocketServer.cmdList = make(map[int]func(msg *interface{}) error)
	SocketServer.upgrader.HandshakeTimeout = time.Duration(SocketServer.TimeOut) * time.Second
	SocketServer.regFunc = make(map[string]*RegFunc)
	return SocketServer
}

// 添加需要监控的方法
func (s *Server) RegFunc(path string, doSomeThing *RegFunc) *Server {
	s.regFunc[path] = doSomeThing
	return s
}

func (s *Server) Begin() {
	go s.timedTick()
}

func (s *Server) Handle(ctx *gin.Context) {
	//特別為測試使用(沒登入訊息的ws)

	conn, err := s.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Error("建立Socket连接失败！原因：" + err.Error())
		return
	}
	defer conn.Close()
	ip := utils.GetClientIp(ctx.Request)

	log.Debug("Websocket:Handle:IP地址为：" + ip + " 的客户端成功连接！")
	if conn, ok := s.GetSafeConn(os.Getenv("socketKey")); ok {
		conn.Handle()
	}
}

// 定时任务
func (s *Server) timedTick() {
	interval := 3 * time.Second
	tick := time.NewTicker(interval)
	for range tick.C {
		SafeConn.Range(func(key, _ interface{}) bool {
			k := key.(string)
			v, ok := s.GetSafeConn(k)
			if ok == false {
				return false
			}
			//如果是使用者斷線60s後關閉連線
			if (time.Now().Unix() - v.lastReceiveTime.Unix()) > s.AutoCloseTime {
				log.Error("ClientID为：" + k + "的连接失效。断开！")
				v.Close()
				SafeConn.Delete(k)
			}
			return true
		})
	}
}

/**
封装并发map  设置conn
*/
func (s *Server) SetSafeConn(keys string, vals *Conns) {
	SafeConn.Store(keys, vals)
}

/**
封装并发map  获取conn
*/

func (s *Server) GetSafeConnByInt(keys int) (regSession *Conns, ok bool) {
	return s.GetSafeConn(strconv.Itoa(keys))
}

func (s *Server) Len() int {
	lengh := 0
	f := func(key, value interface{}) bool {
		lengh++
		return true
	}
	one := lengh
	lengh = 0
	SafeConn.Range(f)
	if one != lengh {
		one = lengh
		lengh = 0
		SafeConn.Range(f)
		if one < lengh {
			return lengh
		}

	}
	return one
}

func (s *Server) ShowIDs() []string {
	var arr []string
	SafeConn.Range(func(key, _ interface{}) bool {
		k := key.(string)
		arr = append(arr, k)
		return true
	})
	return arr
}

func (s *Server) Del(key string) int {
	v, ok := s.GetSafeConn(key)
	if ok == false {
		return 0
	}
	v.Close()
	SafeConn.Delete(key)
	return s.Len()
}

/**
封装并发map  获取conn
*/
func (s *Server) GetSafeConn(keys string) (regSession *Conns, ok bool) {
	//keyss := strconv.Itoa(int(keys))
	keyss := keys
	getVals, ok := SafeConn.Load(keyss)
	if ok == false {
		return nil, false
	}
	regSession, ok = getVals.(*Conns)
	if ok == false {
		return nil, false
	}
	return regSession, true
}
