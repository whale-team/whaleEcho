FROM golang:1.14.6-alpine3.12 AS builder
WORKDIR /whaleEcho
# 啟動 go module
ENV GO111MODULE=on
 
COPY . .
# 交叉編譯 編譯成 linux 環境的 binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o whaleEcho main.go

FROM alpine:latest
ARG BUILD_TIME
ARG SHA1_VER
LABEL maintainer="vic.xu.development@gmail.com"
LABEL project="whale"



RUN apk update && \
    apk upgrade && \
    apk add --no-cache curl tzdata && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/* 

WORKDIR /whaleEcho
COPY --from=builder /whaleEcho/whaleEcho /whaleEcho/whaleEcho

RUN ls
ENV SHA1_VER=${SHA1_VER}
ENV BUILD_TIME=${BUILD_TIME}

# 設定 group && user 資料夾讀與執行的權限
RUN addgroup -g 1000 whaleEcho && \
    adduser -D -u 1000 -G whaleEcho whaleEcho && \
    chown -R whaleEcho:whaleEcho /whaleEcho && \
    chown 550 /whaleEcho

USER whaleEcho
ENTRYPOINT [ "./whaleEcho" ]
CMD [ "ws" ]