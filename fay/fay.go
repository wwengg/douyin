package fay

import "github.com/wwengg/douyin/proto"

type FayProxyServer interface {
	GetConnMgr() ConnManager
	StartWebsocket()
	DoMessage(message *proto.Message)
}
