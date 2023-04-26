package main

type RtmpLive struct {
	Data struct {
		StreamUrl struct {
			RtmpPushUrl string `json:"rtmp_push_url"`
		} `json:"stream_url"`
	} `json:"data"`
}
