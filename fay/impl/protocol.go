package impl

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wwengg/douyin/fay"
	"io"
	"net"
)

type WsProtocol struct {
	//当前连接的socket套接字
	Conn *websocket.Conn
}

// NewWsProtocol 创建连接的方法
func NewWsProtocol(conn *websocket.Conn) fay.Protocol {
	return &WsProtocol{Conn: conn}
}

// Write 写消息Goroutine， 用户将数据发送给客户端
func (c *WsProtocol) Write(data []byte) error {
	if err := c.Conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		fmt.Println("Send Websocket Data error:, ", err)
		return err
	}
	return nil
}

// GetReader 读消息Goroutine，用于从客户端中读取数据
func (c *WsProtocol) GetReader() (r io.Reader, err error) {
	messageType, r, err := c.Conn.NextReader()
	if err != nil {
		fmt.Println("websocket read msg error ", err)
		return nil, err
	}
	if websocket.BinaryMessage != messageType {
		fmt.Println("messageType != websocket.BinaryMessage")
		//return nil, errors.New("messageType != websocket.BinaryMessage")
		return r, nil
	}
	return r, nil
}

func (c *WsProtocol) ConnClose() {
	c.Conn.Close()
}

// RemoteAddr 获取远程客户端地址信息
func (c *WsProtocol) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
