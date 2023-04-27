// @Title
// @Description
// @Author  Wangwengang  2022/3/30 下午11:50
// @Update  Wangwengang  2022/3/30 下午11:50
package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/wwengg/douyin/fay"
	"io"
	"sync"
)

type Connection struct {
	FayProxyServer *FayProxyServer
	//当前连接的socket TCP套接字
	//Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint64
	//消息管理MsgID和对应处理方法的消息管理模块
	//MsgHandler anet2.MsgHandle
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	sync.RWMutex
	//链接属性
	property map[string]interface{}
	////保护当前property的锁
	propertyLock sync.Mutex
	//当前连接的关闭状态
	isClosed bool

	fay.Protocol
}

// NewConnection 创建连接的方法
func NewConnection(fayProxyServer *FayProxyServer, connID uint64, protocol fay.Protocol) fay.Connection {
	//初始化Conn属性
	c := &Connection{
		FayProxyServer: fayProxyServer,
		ConnID:         connID,
		isClosed:       false,
		//MsgHandler:  msgHandler,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, 1024),
		property:    nil,
		Protocol:    protocol,
	}

	return c
}

// StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if err := c.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
			//fmt.Printf("Send data succ! data = %+v\n", data)
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if err := c.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Reader exit!]")
	defer c.Stop()

	// 创建拆包解包的对象
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			r, err := c.GetReader()
			if err != nil {
				//logger.ZapLog.Error("GetReader err:", zap.Error(err))
				return
			}
			if _, err := io.ReadAll(r); err != nil {
				fmt.Println("read msg head error ", err)
				return
			}
		}
	}
}

// Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	//c.TCPServer.CallOnConnStart(c)
}

// Stop 停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	//c.TCPServer.CallOnConnStop(c)

	c.Lock()
	defer c.Unlock()

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)

	// 关闭socket链接
	c.ConnClose()
	//关闭Writer
	c.cancel()

	//将链接从连接管理器中删除
	c.FayProxyServer.GetConnMgr().Remove(c)

	//关闭该链接全部管道
	close(c.msgBuffChan)
	//设置标志位
	c.isClosed = true

}

// GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint64 {
	return c.ConnID
}

// SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	//写回客户端
	c.msgChan <- data

	return nil
}

// SendBuffMsg  发生BuffMsg
func (c *Connection) SendBuffMsg(data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}

	//写回客户端
	c.msgBuffChan <- data

	return nil
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// 返回ctx，用于用户自定义的go程获取连接退出状态
func (c *Connection) Context() context.Context {
	return c.ctx
}
