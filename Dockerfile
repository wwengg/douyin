FROM golang:alpine as builder

WORKDIR /go/src/douyin
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest
LABEL MAINTAINER="wwwwangg@163.com"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk add --no-cache  gettext tzdata curl  && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" >  /etc/timezone && \
    date && \
    apk del tzdata


WORKDIR /go/src/server
COPY --from=0 /go/src/douyin/server ./
COPY --from=0 /go/src/douyin/certificates ./certificates

RUN cp ./certificates/proxy-ca.crt /usr/local/share/ca-certificates/proxy-ca.crt
RUN update-ca-certificates

EXPOSE 8001,8888
ENTRYPOINT ./server
