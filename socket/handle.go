package handle

import (
	"server/packages/socket"
)

func HandleFunc() *socket.RegFunc {
	handle := new(socket.RegFunc)
	handle.Connect = onConnect
	handle.HasMsg = onMessage
	handle.Close = onClose
	return handle
}

//建立連線時傳一次
func onConnect(conn *socket.Conns) error {

	//任何type 都可以
	err := conn.Send("test")
	return err
}

// 接收到socket消息
func onMessage(receiveMsg *socket.ReceiveMsg, conn *socket.Conns) error {
	err := conn.Send("1234")
	if err != nil {
		return err
	}

	return nil
}

func onClose(conn *socket.Conns) error {
	return nil
}
