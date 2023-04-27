module github.com/wwengg/douyin

go 1.16

require (
	github.com/elazarl/goproxy v0.0.0-20221015165544-a0805db90819
	github.com/gogo/protobuf v1.3.2
	github.com/gorilla/websocket v1.5.0
)

replace github.com/elazarl/goproxy v0.0.0-20221015165544-a0805db90819 => github.com/wwengg/goproxy v0.0.2

//replace github.com/elazarl/goproxy v0.0.0-20221015165544-a0805db90819 => ../../../goproxy
