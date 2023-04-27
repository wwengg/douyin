package main

import (
	"bytes"
	"encoding/json"
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

	proxy.OnRequest().
		HandleConnect(goproxy.AlwaysMitm)

	//proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	//	if req.URL.Scheme == "https" {
	//		req.URL.Scheme = "http"
	//	}
	//	return req, nil
	//})
	proxy.AddWebsocketHandler(func(data []byte, direction goproxy.WebsocketDirection, ctx *goproxy.ProxyCtx) []byte {
		log.Print(direction)
		log.Print(string(data))
		return nil
	})

	proxy.OnResponse(goproxy.UrlHasPrefix("https://webcast.amemv.com/webcast/room/create/")).DoFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			//httpResponse, _ := httputil.DumpResponse(resp, true)
			//res := unmarshalHTTPResponse(httpResponse)
			buf, _ := io.ReadAll(resp.Body)
			responseStream := io.NopCloser(bytes.NewBuffer(buf))
			rtmpLive := RtmpLive{}
			json.Unmarshal(buf, &rtmpLive)

			//log.Println(string(httpResponse))
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

