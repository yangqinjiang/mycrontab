# 使用这个镜像编译代码
#
# BUILD 阶段
#
FROM golang:1.12.7 AS build
# 设置我们应用程序的工作目录
WORKDIR /go/src/app
# 添加所有编译的应用代码
COPY . .
#下载代码
RUN go get -d -v ./...
RUN go install -v ./...
# 编译一个静态的go应用 (在二进制构建中包含c语言依赖库)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

#设置我们应用程序的启动命令

CMD ["app"]
#
# 生产阶段
#
FROM scratch AS prod

# 从build阶段拷贝二进制文件
COPY --from=build /go/src/app .
CMD ["./app"]
# 运行
#$ docker build -t my-golang-app .
#$ docker images|grep my-golang-app                       
# 输出 my-golang-app / latest / 43c9e3802f10 /  16 seconds ago /   816MB
#$ docker run -it --rm --name my-running-app my-golang-app
# 输出 `hello world`


# 如何为你的go应用创建轻量级Docker镜像
# http://www.dockone.one/article/8196