package impl

import (
	"errors"
	"fmt"
	"github.com/wwengg/douyin/fay"
	"sync"
)

type ConnManager struct {
	connections map[uint64]fay.Connection
	connLock    sync.RWMutex
}

// NewConnManager 创建一个链接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint64]fay.Connection),
	}
}

// Add 添加链接
func (connMgr *ConnManager) Add(conn fay.Connection) {

	connMgr.connLock.Lock()
	//将conn连接添加到ConnMananger中
	connMgr.connections[conn.GetConnID()] = conn
	connMgr.connLock.Unlock()

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

// Remove 删除连接
func (connMgr *ConnManager) Remove(conn fay.Connection) {

	connMgr.connLock.Lock()
	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	connMgr.connLock.Unlock()
	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

// Get 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint64) (fay.Connection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")

}

// Len 获取当前连接
func (connMgr *ConnManager) Len() int {
	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()
	return length
}

// ClearConn 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	connMgr.connLock.Unlock()
	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}

// ClearOneConn  利用ConnID获取一个链接 并且删除
func (connMgr *ConnManager) ClearOneConn(connID uint64) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connections := connMgr.connections
	if conn, ok := connections[connID]; ok {
		//停止
		conn.Stop()
		//删除
		delete(connections, connID)
		fmt.Println("Clear Connections ID:  ", connID, "succeed")
		return
	}

	fmt.Println("Clear Connections ID:  ", connID, "err")
	return
}

// ClearConn 清除并停止所有连接
func (connMgr *ConnManager) SendMsgToAllConn(data []byte) {
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for _, conn := range connMgr.connections {
		//发送
		_ = conn.SendBuffMsg(data)

	}
	connMgr.connLock.Unlock()
	fmt.Println("Send Data to All Connections successfully: conn num = ", connMgr.Len())
}
