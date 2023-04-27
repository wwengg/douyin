// @Title
// @Description
// @Author  Wangwengang  2022/4/30 下午1:53
// @Update  Wangwengang  2022/4/30 下午1:53
package fay

import (
	"io"
	"net"
)

type Protocol interface {
	Write(data []byte) error
	GetReader() (r io.Reader, err error)
	ConnClose()
	RemoteAddr() net.Addr
}
