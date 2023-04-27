package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/wwengg/douyin/fay/impl"
	"github.com/wwengg/douyin/proto"
	"github.com/wwengg/douyin/utils"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

func main() {
	configureCA()
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false

	fayProxyServer := impl.NewFayProxyServer()
	go fayProxyServer.StartWebsocket()

	//ws数据处理
	proxy.AddWebsocketHandler(func(data []byte, direction goproxy.WebsocketDirection, ctx *goproxy.ProxyCtx) (reply []byte) {
		reply = data
		if len(data) == 0 {
			return
		}
		if data[0] != 0x08 {
			return
		}
		wssResponse := proto.WssResponse{}
		if err := wssResponse.XXX_Unmarshal(data); err == nil {
			//检测包格式
			if v, ok := wssResponse.Headers["compress_type"]; !ok && v != "gzip" {
				return
			}
			//解压gzip
			deData, err := utils.GzipDecode(wssResponse.Payload)
			if err != nil {
				ctx.Logf("gzip解压失败")
				return
			}
			res := proto.Response{}
			if err = res.XXX_Unmarshal(deData); err != nil {
				return
			}
			for _, message := range res.Messages {
				fayProxyServer.DoMessage(message)
			}
		}
		return
	})

	proxy.WebSocketHandler = func(dst io.Writer, src io.Reader, direction goproxy.WebsocketDirection, ctx *goproxy.ProxyCtx) error {
		fullPacket := make([]byte, 0)
		buf := make([]byte, 32*1024)
		var err error = nil
		for {
			nr, er := src.Read(buf)

			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}

			if nr > 0 {
				fullPacket = append(fullPacket, buf[:nr]...)
				websocketPacket := newWebsocketPacket(fullPacket)

				if !websocketPacket.valid {
					continue
				}

				websocketPacket.payload = proxy.FilterWebsocketPacket(websocketPacket.payload, direction, ctx)
				encodedPacket := websocketPacket.encode()
				nw, ew := dst.Write(encodedPacket)
				fullPacket = fullPacket[websocketPacket.packetSize:]

				if nw < 0 || len(encodedPacket) < nw {
					nw = 0
					if ew == nil {
						ew = errors.New("invalid write result")
					}
				}
				if ew != nil {
					err = ew
					break
				}
				if len(encodedPacket) != nw {
					err = io.ErrShortWrite
					break
				}
			}
		}
		return err
	}

	proxy.OnRequest(goproxy.ReqHostIs("webcast.amemv.com:443", "frontier-im.douyin.com:443", "webcast100-ws-web-lq.amemv.com:443", "webcast3-ws-web-lf.douyin.com:443", "webcast3-ws-web-hl.douyin.com:443")).
		HandleConnect(goproxy.AlwaysMitm)

	proxy.OnResponse(goproxy.UrlHasPrefix("httpswebcast.amemv.com:443/webcast/room/create/")).DoFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			buf, _ := io.ReadAll(resp.Body)
			responseStream := io.NopCloser(bytes.NewBuffer(buf))
			rtmpLive := RtmpLive{}
			json.Unmarshal(buf, &rtmpLive)

			log.Println(rtmpLive)
			url := rtmpLive.Data.StreamUrl.RtmpPushUrl
			array := strings.Split(url, "/")
			secret := array[len(array)-1]
			serverName := strings.Split(url, secret)[0]
			log.Printf(`服务器：%s`, serverName)
			log.Printf(`推流码：%s`, secret)
			resp.Body = responseStream
			return resp
		},
	)
	log.Println("软件准备就绪，请启动【直播伴侣】并且点击【开始直播】")
	log.Fatal(http.ListenAndServe(":8001", proxy))
}
