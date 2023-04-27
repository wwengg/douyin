package impl

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/wwengg/douyin/fay"
	"github.com/wwengg/douyin/proto"
	"log"
	"net/http"
	"time"
)

type FayProxyServer struct {
	// Websocket Addr
	WsAddr string
	//用于缓存接收过的消息ID，判断是否重复接收
	Dictionary map[string][]int64
	//当前Server的链接管理器
	ConnMgr fay.ConnManager

	GenID uint64
}

func NewFayProxyServer() *FayProxyServer {
	return &FayProxyServer{
		WsAddr:     "127.0.0.1:8888",
		Dictionary: make(map[string][]int64),
		ConnMgr:    NewConnManager(),
	}
}

// GetConnMgr 得到链接管理
func (s *FayProxyServer) GetConnMgr() fay.ConnManager {
	return s.ConnMgr
}

// Start Websocket网络服务
func (s *FayProxyServer) StartWebsocket() {
	//logger.ZapLog.Info("Start Websocket server", zap.String("addr", s.WsAddr))
	log.Printf("Start Websocket server addr:%s", s.WsAddr)
	httpServer := &http.Server{
		Addr: s.WsAddr,
		Handler: &WsHandler{upgrader: websocket.Upgrader{
			HandshakeTimeout: 0,
			ReadBufferSize:   0,
			WriteBufferSize:  0,
			WriteBufferPool:  nil,
			Subprotocols:     nil,
			Error:            nil,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			EnableCompression: false,
		},
			sv: s},
		ReadTimeout:    time.Second * time.Duration(60),
		WriteTimeout:   time.Second * time.Duration(60),
		MaxHeaderBytes: 4096,
	}
	httpServer.ListenAndServe()
}

type WsHandler struct {
	sv       *FayProxyServer
	upgrader websocket.Upgrader
}

func (h *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("upgrade error,err:%s", err.Error())
		return
	}

	//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
	if h.sv.ConnMgr.Len() >= 10 {
		conn.Close()
	}

	//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
	h.sv.GenID++
	dealConn := NewConnection(h.sv, h.sv.GenID, NewWsProtocol(conn))

	h.sv.GetConnMgr().Add(dealConn)
	//3.4 启动当前链接的处理业务
	go dealConn.Start()

}

func (s *FayProxyServer) DoMessage(message *proto.Message) {
	list, ok := s.Dictionary[message.Method]
	if !ok {
		list = []int64{}
	} else {
		for _, i := range list {
			if i == message.MsgId {
				return
			}
		}
	}
	if len(list) > 300 {
		list = []int64{}
	}
	log.Printf("method:%s,num:%d", message.Method, len(list))
	list = append(list, message.MsgId)
	s.Dictionary[message.Method] = list
	switch message.Method {
	case "WebcastMemberMessage":
		memberMessage := proto.MemberMessage{}
		memberMessage.XXX_Unmarshal(message.Payload)
		s.onMemberMessgae(memberMessage)
		break
	case "WebcastChatMessage":
		chatMessage := proto.ChatMessage{}
		chatMessage.XXX_Unmarshal(message.Payload)
		s.onChatMessage(chatMessage)
		break
	default:
		break
	}

}

func (s *FayProxyServer) send(pack *fay.MsgPack) {

	data, _ := json.Marshal(pack)
	s.GetConnMgr().SendMsgToAllConn(data)

}

func (s *FayProxyServer) onMemberMessgae(message proto.MemberMessage) {
	data, _ := json.Marshal(message)
	s.send(fay.CreateMsgPack(string(data), fay.MsgType_JoinRoom))

}

func (s *FayProxyServer) onChatMessage(message proto.ChatMessage) {
	data, _ := json.Marshal(message)
	s.send(fay.CreateMsgPack(string(data), fay.MsgType_DanMu))
}
