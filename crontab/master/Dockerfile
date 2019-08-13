# 使用这个镜像编译代码
#
# BUILD 阶段
#
FROM golang:1.12.7 AS build
LABEL maintainer="yangqinjiang@qq.com"
# 设置go mod的必要条件
ENV GOPROXY https://goproxy.io
ENV GO111MODULE on

WORKDIR /go/cache
# 下载mod
COPY go.mod .
COPY go.sum .
RUN go mod download
# 编译目录
WORKDIR /go/release

COPY . .
# 编译一个静态的go应用 (在二进制构建中包含c语言依赖库)
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o master main.go

#
# 生产阶段
#
FROM scratch AS prod
LABEL maintainer="yangqinjiang@qq.com"
# 从build阶段拷贝二进制文件
COPY --from=build /go/release/master /
# master的web文件
COPY --from=build /go/release/webroot /webroot/
# master的配置文件
COPY --from=build /go/release/config/master.json /config/master.json
# 开放端口
EXPOSE 8080
ENTRYPOINT ["./master"]




# 运行
#$ docker build -t my-golang-app .
#$ docker images|grep my-golang-app                       
# 输出 my-golang-app / latest / 43c9e3802f10 /  16 seconds ago /   816MB
#$ docker run -it --rm --name my-running-app my-golang-app
# 输出 `hello world`


# 如何为你的go应用创建轻量级Docker镜像
# http://www.dockone.one/article/8196