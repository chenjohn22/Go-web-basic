package socket

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/CRGao/log"

	"github.com/gorilla/websocket"
)

type MsgMethod struct {
	Method string `json:"method"`
}
type ReceiveMsg struct {
	MsgMethod
	Data Data `json:"data"`
}
type Data struct {
	Msg string `json:"msg"`
}

type GamePing struct {
	MsgMethod
	Data Data `json:"data"`
}

type Conns struct {
	Key              string
	conn             *websocket.Conn
	lastReceiveTime  time.Time //最後一次接收到數據的時間
	lastSendTime     time.Time //最後一次發送數據的時間
	lastActivityTime time.Time //最後一次操作的時間
	*sync.Mutex
	main          *Server
	lastHeartBeat string
	status        ConnStatus
	errorTime     int  //錯誤次數
	noticeConned  bool //已預告斷線
	disConned     bool //已通知斷線
}

func NewConn(conn *websocket.Conn, main *Server, key string) *Conns {
	server := new(Conns)
	server.Key = key
	server.conn = conn
	server.lastReceiveTime = time.Now()
	server.lastSendTime = time.Now()
	server.lastActivityTime = time.Now()
	server.Mutex = &sync.Mutex{}
	server.main = main
	server.status = Connected
	server.doConnectFunc()
	return server
}

//当创建链接时需要执行的方法
func (c *Conns) doConnectFunc() {
	for _, v := range c.main.regFunc {
		if v.Connect != nil {
			err := v.Connect(c)
			if err != nil {
				log.Error(err)
			}

		}
	}
}

func (c *Conns) doCloseFunc() {
	for _, v := range c.main.regFunc {
		if v.Close != nil {
			err := v.Close(c)
			if err != nil {
				log.Error(err)
			}
		}
	}
}
func (c *Conns) Handle() {
	c.heartBeat()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
				log.Info("socket:server:客户端:" + c.Key + " 斷開連接")
				c.Close()
				break
			}
			log.Error("socket:server:讀取Socket數據失敗！原因：" + err.Error())
			break
		}
		msg, err := c.parseMsg(message)
		if err != nil {
			log.Error("socket:server:解析客服端訊息失敗，錯誤訊息：", err, "。消息內容：", string(message))
			continue
		}
		c.lastReceiveTime = time.Now()
		if c.checkHeartBeatMsg(&msg) == true {
			continue
		}
		if c.main.regFunc[msg.Method] != nil && c.main.regFunc[msg.Method].HasMsg != nil {
			if err := c.main.regFunc[msg.Method].HasMsg(&msg, c); err != nil {
				log.Error("socket:server:執行註冊方法失敗！錯誤訊息：", err)
			}
		}
	}
}

// 解析客户端发送过来的消息
func (c *Conns) parseMsg(msg []byte) (ReceiveMsg, error) {
	var parseMsg ReceiveMsg
	err := json.Unmarshal(msg, &parseMsg)
	return parseMsg, err
}

// 单个连接发送消息
func (c *Conns) Send(msg interface{}) error {
	if c.status == ConnectClose {
		c.lastReceiveTime = c.lastReceiveTime.Add(time.Second * time.Duration(30*-1))
		log.Error("socket:conn:客户端:" + c.Key + " 已斷開連接！")
		return nil
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	err = c.conn.WriteMessage(websocket.TextMessage, data)
	if err == nil {
		c.lastSendTime = time.Now()
		c.errorTime = 0
	} else {
		c.errorTime++
	}
	return err
}

func (c *Conns) Close() {
	//解决 https://github.com/gorilla/websocket/issues/119
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Error(err)
	}
	c.status = ConnectClose
	c.doCloseFunc()
	err = c.conn.Close()
	if err != nil {
		log.Error(err)
	}

}

// 发送心跳和检查心跳
func (c *Conns) heartBeat() {

	tick := time.NewTicker(5 * time.Second) //每5秒进行一次心跳
	go func(tick *time.Ticker) {
		for range tick.C {
			if c.status == ConnectClose {
				tick.Stop()
				break
			}
			if time.Now().Unix()-c.lastReceiveTime.Unix() > 30 {
				// 如果已经有超过30秒没有收到消息，则自动断开
				log.Error("socket:conn:客户端：", c.Key, "，已經連續30秒沒有收到訊息，主動斷開連接！")
				c.Close()
				break
			}

			// if time.Now().Unix()-c.lastActivityTime.Unix() > int64(config.ServerConfig.GameTime*config.ServerConfig.CloseConnRound) {
			// 	log.Error("socket:conn:客戶端：", c.Key, "，已經連續", config.ServerConfig.GameTime*config.ServerConfig.CloseConnRound, "秒沒有操作，強制斷開連結！")
			// 	c.Close()
			// 	break
			// }

			//調整斷線方式 以下先註解
			//if !c.noticeConned && time.Now().Unix()-c.lastActivityTime.Unix() > int64(config.ServerConfig.GameTime * config.ServerConfig.AlertCloseConnRound) {
			//	log.Error("socket:conn:客户端：", c.Key, "，已經連續一段時間没有操作，預先告知前端準備要斷開連接！")
			//	msg := new(GamePing)
			//	msg.Method = "notice_connection"
			//	msg.Data.Msg = c.lastHeartBeat
			//	err := c.Send(msg)
			//	if err != nil {
			//		log.Error("socket:server:發送關閉連結失敗！錯誤訊息：", err)
			//	}
			//	c.noticeConned = true
			//}
			//
			//if !c.disConned && time.Now().Unix()-c.lastActivityTime.Unix() > int64(config.ServerConfig.GameTime * config.ServerConfig.CloseConnRound) {
			//	log.Error("socket:conn:客户端：", c.Key, "，已經連續一段時間没有操作，提醒前端斷開連接！")
			//	msg := new(GamePing)
			//	msg.Method = "close_connection"
			//	msg.Data.Msg = c.lastHeartBeat
			//	err := c.Send(msg)
			//	if err != nil {
			//		log.Error("socket:server:發送關閉連結失敗！錯誤訊息：", err)
			//	}
			//	c.disConned = true
			//}

			if c.errorTime > 10 {
				// 如果已经连继超过10次错误了
				log.Error("socket:conn:客户端：", c.Key, "，已經連續10次錯誤，主動斷開連接！")
				c.Close()
				break
			}
			c.lastHeartBeat = strconv.FormatInt(time.Now().Unix(), 10)
			msg := new(GamePing)
			msg.Method = "ping"
			msg.Data.Msg = c.lastHeartBeat
			err := c.Send(msg)
			if err != nil {
				log.Error("socket:server:發送心跳服務失敗！錯誤訊息：", err)
			}
		}
	}(tick)
}

func (c *Conns) checkHeartBeatMsg(msg *ReceiveMsg) bool {
	if msg.Method == "pong" && msg.Data.Msg == c.lastHeartBeat {
		c.lastReceiveTime = time.Now()
		return true
	}
	if msg.Method == "close" {
		c.Close()
		return true
	}
	return false
}

func (c *Conns) UpdateActivityTime() {
	c.lastActivityTime = time.Now()
	c.noticeConned = false
	c.disConned = false
}
